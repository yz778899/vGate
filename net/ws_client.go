package net

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"github.com/yz778899/vGate/net/app"
	"github.com/yz778899/vGate/net/coroutine"
	"github.com/yz778899/vGate/net/data"
	"github.com/yz778899/vGate/net/handler"
)

var uuid atomic.Int64

type WsClient struct {
	*data.Session
	Path        string
	pool        *coroutine.CoroutineGroup
	handler     handler.WsHandlerInterface
	isConnected bool
	//Conn    *websocket.Conn
	maxRetries    int           //最大重连次数
	retryInterval time.Duration //重连时间间隔
}

func NewWsWsClient() *WsClient {
	return &WsClient{}
}

// 配置消息处理器
func (this *WsClient) Handler(handler handler.WsHandlerInterface) *WsClient {
	this.handler = handler
	return this
}

// 配置 WsServer 的端口和路径
func (this *WsClient) Config(Path string) *WsClient {
	this.Path = Path
	return this
}

// 常规客户端 重连间隔1秒 连续偿试30分钟
func NewWsClient() *WsClient {
	WsClient := WsClient{}
	WsClient.pool = coroutine.NewCoroutineGroup(1, "WsClient_msg_group", 4)
	WsClient.maxRetries = 1000 * 60 * 30                                                   //30分钟  毫秒数
	WsClient.retryInterval = time.Millisecond * 1000                                       //1秒
	WsClient.maxRetries = WsClient.maxRetries / int(WsClient.retryInterval.Milliseconds()) //保证能重连24小时
	return &WsClient
}

// socket客户端， 重连间隔100毫秒 连续偿试24小时
func NewWsClientAlwaysOnlie() *WsClient {
	WsClient := WsClient{}
	WsClient.pool = coroutine.NewCoroutineGroup(1, "WsClient_msg_group", 4)
	WsClient.maxRetries = 1000 * 60 * 60 * 24 //24小时 毫秒数
	WsClient.retryInterval = time.Millisecond * 1
	WsClient.maxRetries = WsClient.maxRetries / int(WsClient.retryInterval.Milliseconds()) //保证能重连24小时
	return &WsClient
}

// 连接成功
func (this *WsClient) Connect(onConnectedCallBack func(conn *websocket.Conn)) (*websocket.Conn, error) {

	var conn *websocket.Conn
	var err error
	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()

	logOutNum := 20

	for i := 0; i < this.maxRetries; i++ {

		if i%logOutNum == 0 {

			app.Log.Info(fmt.Sprintf("尝试连接 (第 %d/%d 次)...", i+1, this.maxRetries))
		}

		conn, _, err = websocket.DefaultDialer.Dial(this.Path, nil)
		if err == nil {
			app.Log.Info("连接成功！")

			//连接成功
			session := this.handler.OnConnect(conn)
			this.Session = session
			this.isConnected = true

			this.setupHeartbeat()

			if onConnectedCallBack != nil {
				onConnectedCallBack(conn)
			}
			// 循环读取消息
			err = this.readMsg()

			//return conn, err
		}

		// 检查是否是连接拒绝错误
		if isConnectionRefused(err) {
			if i%logOutNum == 0 {
				app.Log.Info(fmt.Sprintf("连接被拒绝，服务可能未启动，%v 后重试...", this.retryInterval))
			}
		} else {
			if i%logOutNum == 0 {
				app.Log.Info(fmt.Sprintf("连接失败: %v，  %v 后重试...", err, this.retryInterval))
			}
		}

		time.Sleep(this.retryInterval)
	}

	return nil, fmt.Errorf("连接失败，已重试 %d 次: %w", this.maxRetries, err)

	// 连接服务器
	// conn, _, err := websocket.DefaultDialer.Dial(this.Path, nil)
	// if err != nil {
	// 	log.Fatal("连接失败:", err)
	// 	return err
	// }
	// defer conn.Close()

}

// isConnectionRefused 检查是否是连接拒绝错误
func isConnectionRefused(err error) bool {
	if err == nil {
		return false
	}

	// 检查是否是 net.OpError 类型
	if opErr, ok := err.(*net.OpError); ok {
		if sysErr, ok := opErr.Err.(*net.OpError); ok {
			return sysErr.Err.Error() == "connectex: No connection could be made because the target machine actively refused it."
		}
	}

	return false
}

// 循环读取消息
func (this *WsClient) readMsg() error {

	conn := this.Session.Conn
	// 接收消息
	for {
		_, jsonMsg, err := conn.ReadMessage()
		if err != nil {
			log.Println("接收失败:", err)
			conn.Close()
			this.isConnected = false
			this.handler.OnDisconnect(this.Session)
			return err
		}

		theMsg := data.NoDecoderMsg{
			SessionId: uuid.Add(1),
			Msg:       string(jsonMsg),
			SnId:      rand.Intn(len(this.pool.Slave)), //够slave取模就可以了
		}

		fmt.Print(theMsg)

		var v data.NoDecoderMsg = theMsg
		WsMsg, _ := data.ServerDecoder(v)
		this.handler.OnMessage(conn, WsMsg)
	}
}

// 设置心跳
func (this *WsClient) setupHeartbeat() {
	// 启动心跳发送 goroutine
	go this.sendHeartbeat()
}

// sendHeartbeat 定期发送 Ping
func (this *WsClient) sendHeartbeat() {
	ticker := time.NewTicker(time.Duration(app.VGate.Config.Gate.HeartbeatTime) * time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C // 阻塞等待 ticker 信号
		if err := this.Session.Conn.WriteJSON(data.HeartbeatMsg()); err != nil {
			app.Log.Info(fmt.Sprintf("发送 heartbeatMsg 失败: %v", err))
			return
		}
		// 设置读取超时
		this.Conn.SetReadDeadline(time.Now().Add(time.Duration(app.VGate.Config.Gate.ReadOverTime) * time.Second))
		app.Log.Info("发送 heartbeatMsg")
	}
}
