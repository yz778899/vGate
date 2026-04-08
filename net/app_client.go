package net

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/gorilla/websocket"
	"github.com/yz778899/vGate/net/env"
	"github.com/yz778899/vGate/net/msg"
	data "github.com/yz778899/vGate/net/msg"

	"github.com/yz778899/vGate/net/handler"
)

type AppClient struct {
	ConnSession   *msg.ConnSession
	Url           string
	handler       handler.ServiceAcceptInterface
	isConnected   bool
	maxRetries    int           //最大重连次数
	retryInterval time.Duration //重连时间间隔
}

// 配置消息处理器
func (this *AppClient) Handler(handler handler.ServiceAcceptInterface) *AppClient {
	this.handler = handler
	return this
}

// 配置 WsServer 的端口和路径
func (this *AppClient) Config(url string) *AppClient {
	this.Url = url
	return this
}

// 常规客户端 重连间隔1秒 连续偿试30分钟
func NewAppClient() *AppClient {
	WsClient := AppClient{}
	WsClient.maxRetries = 1000 * 60 * 1                                                    //1分钟  毫秒数
	WsClient.retryInterval = time.Millisecond * 1000                                       //1秒
	WsClient.maxRetries = WsClient.maxRetries / int(WsClient.retryInterval.Milliseconds()) //保证能重连24小时
	return &WsClient
}

// 连接成功
func (this *AppClient) Connect(onConnectedCallBack func(conn *websocket.Conn)) (*websocket.Conn, error) {

	var conn *websocket.Conn
	var err error
	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()

	logOutNum := 20

	for i := 0; i < this.maxRetries; i++ {

		if i > 0 && i%logOutNum == 0 {
			env.Log.Info(fmt.Sprintf("连接 %v (第 %d/%d 次)...", this.Url, i+1, this.maxRetries))
		}

		conn, _, err = websocket.DefaultDialer.Dial(this.Url, nil)
		if err == nil {
			//env.Log.Info("连接成功！")

			//连接成功
			this.handler.OnConnect(conn)
			connSession := msg.ConnSession{Conn: conn,
				UUID: int64(0),
			}
			this.ConnSession = &connSession
			this.isConnected = true
			this.setupHeartbeat()
			if onConnectedCallBack != nil {
				onConnectedCallBack(conn)
			}
			// 循环读取消息
			err = this.readMsg()
		}

		// 检查是否是连接拒绝错误
		if this.isConnectionRefused(err) {
			if i%logOutNum == 0 {
				env.Log.Info(fmt.Sprintf("连接被拒绝，服务可能未启动，%v 后重试...", this.retryInterval))
			}
		} else {
			if i%logOutNum == 0 {
				env.Log.Info(fmt.Sprintf("连接失败: %v，  %v 后重试...", err, this.retryInterval))
			}
		}

		time.Sleep(this.retryInterval)
	}

	return nil, fmt.Errorf("连接失败，已重试 %d 次: %w", this.maxRetries, err)

}

// isConnectionRefused 检查是否是连接拒绝错误
func (this *AppClient) isConnectionRefused(err error) bool {
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
func (this *AppClient) readMsg() error {

	conn := this.ConnSession.Conn
	session := this.handler.OnConnect(conn)

	// 接收消息
	for {
		_, originalMsgByteArray, err := conn.ReadMessage()
		if err != nil {
			log.Println("接收失败:", err)
			conn.Close()
			this.isConnected = false
			//this.handler.OnDisconnect(this.SesConnSession)
			this.handler.OnDisconnect(nil)
			return err
		}

		theMsg := data.NoDecoderMsg{
			SessionId: uuid.Add(1),
			Msg:       string(originalMsgByteArray),
			SnId:      rand.Intn(128), //够slave取模就可以了
		}

		countMsg.Add(1)
		num := countMsg.Load()
		if num%int64(2) == 0 {
			fmt.Printf(" app_service revicer , msg count = %v  : %v", countMsg.Load(), theMsg)
		}

		var v data.NoDecoderMsg = theMsg
		WsMsg, _ := data.ServerDecoder(v)

		this.handler.OnMessage(&handler.WebSocketContext{
			Session:  session,
			Original: &originalMsgByteArray,
			WsMsg:    WsMsg,
		})
	}
}

// 设置心跳
func (this *AppClient) setupHeartbeat() {
	// 启动心跳发送 goroutine
	go this.sendHeartbeat()
}

// sendHeartbeat 定期发送 Ping
func (this *AppClient) sendHeartbeat() {
	ticker := time.NewTicker(time.Duration(env.VGate.Config.Gate.HeartbeatTime) * time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C // 阻塞等待 ticker 信号
		if err := this.ConnSession.Conn.WriteJSON(data.HeartbeatMsg()); err != nil {
			env.Log.Info(fmt.Sprintf("发送 heartbeatMsg 失败: %v", err))
			return
		}
		// 设置读取超时
		this.ConnSession.Conn.SetReadDeadline(time.Now().Add(time.Duration(env.VGate.Config.Gate.ReadOverTime) * time.Second))
		env.Log.Info("发送 heartbeatMsg")
	}
}
