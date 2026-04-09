package net

import (
	"fmt"
	"math/rand"
	"net"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"github.com/yz778899/vGate/net/env"
	"github.com/yz778899/vGate/net/msg"
	pool "github.com/yz778899/vGate/net/poll"
	"go.uber.org/zap"

	"github.com/yz778899/vGate/net/handler"
)

var uuid atomic.Int64

var countMsg atomic.Int64
var Log *zap.Logger

func init() {
	Log = env.Log
}

type AppService struct {
	*msg.Session
	Url           string
	Pool          *pool.QueueMaster[*pool.MessageTask]
	handler       handler.ServiceAcceptInterface
	isConnected   bool
	maxRetries    int           //最大重连次数
	retryInterval time.Duration //重连时间间隔
	poolNum       int           //消息池的线程数量
}

// 配置消息处理器
func (this *AppService) Handler(handler handler.ServiceAcceptInterface) *AppService {
	this.handler = handler
	return this
}

// 配置 WsServer 的端口和路径
func (this *AppService) Config(url string, poolNum int) *AppService {
	this.Url = url
	this.poolNum = poolNum
	this.Pool = pool.NewCoroutineGroup[*pool.MessageTask](1, "WsLine_msg_group", this.poolNum)
	return this
}

// socket客户端， 重连间隔100毫秒 连续偿试24小时
func NewAppService() *AppService {
	WsLine := AppService{}
	WsLine.maxRetries = 1000 * 60 * 60 * 24 //24小时 毫秒数
	WsLine.retryInterval = time.Millisecond * 1
	WsLine.maxRetries = WsLine.maxRetries / int(WsLine.retryInterval.Milliseconds()) //保证能重连24小时
	return &WsLine
}

// 连接成功
func (this *AppService) Connect(onConnectedCallBack func(conn *websocket.Conn)) (*websocket.Conn, error) {

	var conn *websocket.Conn
	var err error
	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()

	logOutNum := 10

	for i := 0; i < this.maxRetries; i++ {

		if i%logOutNum == 0 {
			env.Log.Error(fmt.Sprintf("尝试连接 %v (第 %d/%d 次)...", this.Url, i+1, this.maxRetries))
		}

		conn, _, err = websocket.DefaultDialer.Dial(this.Url, nil)
		if err == nil {
			env.Log.Error("连接成功！")

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
func (this *AppService) isConnectionRefused(err error) bool {
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
func (this *AppService) readMsg() error {

	conn := this.Session.Conn
	session := this.handler.OnConnect(conn)

	// 接收消息
	for {
		_, originalMsgByteArray, err := conn.ReadMessage()
		if err != nil {
			Log.Error("接收失败:", zap.Any("error", err))
			conn.Close()
			this.isConnected = false
			this.handler.OnDisconnect(this.Session)
			return err
		}

		theMsg := msg.NoDecoderMsg{
			SessionId: uuid.Add(1),
			Msg:       string(originalMsgByteArray),
			SnId:      rand.Intn(len(this.Pool.Slave)), //够slave取模就可以了
		}

		var v msg.NoDecoderMsg = theMsg
		WsMsg, _ := msg.ServerDecoder(v)
		if WsMsg.Cmd == msg.Notice || WsMsg.Cmd == msg.Unknown {
			//通知暂时不处理
		} else {
			this.handler.OnMessage(&handler.WebSocketContext{
				Session:  session,
				Original: &originalMsgByteArray,
				WsMsg:    WsMsg,
			})
		}

	}
}

// 设置心跳
func (this *AppService) setupHeartbeat() {
	// 启动心跳发送 goroutine
	go this.sendHeartbeat()
}

// sendHeartbeat 定期发送 Ping
func (this *AppService) sendHeartbeat() {
	ticker := time.NewTicker(time.Duration(env.VGate.Config.Gate.HeartbeatTime) * time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C // 阻塞等待 ticker 信号
		if err := this.Session.Conn.WriteJSON(msg.HeartbeatMsg()); err != nil {
			env.Log.Info(fmt.Sprintf("send heartbeatMsg faild: %v", err))
			return
		}
		// 设置读取超时
		this.Session.Mutex.Lock()
		this.Session.Conn.SetReadDeadline(time.Now().Add(time.Duration(env.VGate.Config.Gate.ReadOverTime) * time.Second))
		this.Session.Mutex.Unlock()
		env.Log.Info("send heartbeatMsg")
	}
}
