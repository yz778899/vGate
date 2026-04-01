package net

import (
	"fmt"
	"log"
	"math/rand"
	"sync/atomic"
	"time"

	"github.com/14132465/vGate/net/app"
	"github.com/14132465/vGate/net/coroutine"
	"github.com/14132465/vGate/net/data"
	"github.com/14132465/vGate/net/handler"
	"github.com/gorilla/websocket"
)

var uuid atomic.Int64

type WsClient struct {
	Path    string
	pool    *coroutine.CoroutineGroup
	handler handler.WsHandlerInterface
	//Conn    *websocket.Conn
	*data.Session
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

func NewWsClient() *WsClient {
	WsClient := WsClient{}
	WsClient.pool = coroutine.NewCoroutineGroup(1, "WsClient_msg_group", 4)
	return &WsClient
}

// 连接成功
func (this *WsClient) Connect(onConnectedCallBack func(conn *websocket.Conn)) {
	// 连接服务器
	conn, _, err := websocket.DefaultDialer.Dial(this.Path, nil)
	if err != nil {
		log.Fatal("连接失败:", err)
	}
	defer conn.Close()

	//连接成功
	session := this.handler.OnConnect(conn)
	this.Session = session

	this.setupHeartbeat()

	if onConnectedCallBack != nil {
		go onConnectedCallBack(conn)
	}
	// 接收消息
	for {
		_, jsonMsg, err := conn.ReadMessage()
		if err != nil {
			log.Println("接收失败:", err)
			this.handler.OnDisconnect(session)
			return
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
			log.Printf("发送 heartbeatMsg 失败: %v", err)
			return
		}
		// 设置读取超时
		this.Conn.SetReadDeadline(time.Now().Add(time.Duration(app.VGate.Config.Gate.ReadOverTime) * time.Second))
		log.Println("发送 heartbeatMsg")
	}
}
