package net

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/14132465/vGate/net/app"
	"github.com/14132465/vGate/net/data"
	"github.com/14132465/vGate/net/handler"

	ws "github.com/gorilla/websocket"
)

// WsServer结构体表示一个WebSocket服务器，包含端口、路径、协程池和消息处理器等信息
type WsServer struct {
	Port string
	Path string
	//pool    *coroutine.CoroutineGroup
	handler handler.WsHandlerInterface
	//fun     func(msg data.WsMsg)
}

// 配置 WsServer 的端口和路径
func (this *WsServer) Config(Port int, Path string) *WsServer {
	this.Port = strconv.Itoa(Port)
	this.Path = Path
	return this
}

// 创建 WsServer 实例
func NewWsServer() *WsServer {
	ws := WsServer{}
	return &ws
}

// 配置消息处理器
func (this *WsServer) Handler(handler handler.WsHandlerInterface) *WsServer {
	this.handler = handler
	return this
}

// 运行 WsServer，监听指定端口并处理 WebSocket 连接和消息
func (this *WsServer) Run() *WsServer {
	http.HandleFunc(this.Path, this.wsServerHandler)
	log.Println("WsServer run , port = " + this.Port)
	log.Fatal(http.ListenAndServe(":"+this.Port, nil))
	return this
}

// 配置 Upgrader，用于将 HTTP 连接升级为 WsServer
var upgrader = ws.Upgrader{
	ReadBufferSize:  1024 * 8,
	WriteBufferSize: 1024 * 8,
	// 开发时允许所有跨域请求，生产环境需要严格校验
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 处理 WebSocket 连接和消息
func (this *WsServer) wsServerHandler(w http.ResponseWriter, r *http.Request) {
	// 升级 HTTP 连接为 WsServer
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("升级失败:", err)
		return
	}
	defer conn.Close()

	//读超时 10秒
	conn.SetReadDeadline(time.Now().Add(time.Duration(app.VGate.Config.Gate.ReadOverTime) * time.Second))

	session := this.handler.OnConnect(conn)

	for {
		// 读取客户端消息
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("关闭通道:", err)
			this.handler.OnDisconnect(session)
			break
		}

		fmt.Printf("网关收到消息 msg = %v \n", string(msg))

		var theMsg data.NoDecoderMsg

		theMsg = data.NoDecoderMsg{
			SessionId: session.UUID,
			Msg:       string(msg),
			SnId:      rand.Intn(1024),
		}

		WsMsg, _ := data.GateDecoder(theMsg)

		fmt.Printf("处理后的消息 msg = %#v \n", WsMsg)

		this.handler.OnMessage(conn, WsMsg)

	}
	this.handler.OnDisconnect(session)
}

// // sendHeartbeat 定期发送 Ping
// func (c *WsServer) sendHeartbeat(session *data.Session) {
// 	ticker := time.NewTicker(5 * time.Second)
// 	defer ticker.Stop()

// 	for {
// 		select {
// 		case <-ticker.C:
// 			// 发送 Ping 消息
// 			if err := session.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
// 				log.Printf("发送 Ping 失败: %v", err)
// 				return
// 			}
// 			log.Println("发送 Ping")

// 		case <-c.send:
// 			// 有其他消息发送，继续
// 			continue
// 		}
// 	}
// }
