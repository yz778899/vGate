package net

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/14132465/vGate/net/data"
	"github.com/14132465/vGate/net/handler"

	ws "github.com/gorilla/websocket"
)

// WsServer结构体表示一个WebSocket服务器，包含端口、路径、协程池和消息处理器等信息
type WsServer struct {
	Port string
	Path string
	//pool    *coroutine.CoroutineGroup
	handler handler.HandlerInterface
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
	//ws.pool = coroutine.NewCoroutineGroup(1, "ws_msg_group", 4)

	// fun := func(msg coroutine.V1Msg) {

	// 	if nd, ok := msg.(data.NoDecoderMsg); ok {
	// 		_, wsMsg := data.Decoder(nd)

	// 		fmt.Print("wsMsg ######## ,", wsMsg)
	// 		session := app.SessionManager.GetSession(wsMsg.SessionId)
	// 		if session != nil {
	// 			ws.handler.OnMessage(session.Conn, wsMsg)
	// 		} else {
	// 			//"未找到会话ID  可能已经下线了;
	// 			fmt.Printf("未找到会话ID %d 对应的服务器\n", wsMsg.SessionId)
	// 		}

	// 	} else {
	// 		fmt.Printf("无法解析的消息类型 %v ", msg)
	// 	}
	// }

	// ws.pool.Handler(fun)
	return &ws
}

// 配置消息处理器
func (this *WsServer) Handler(handler handler.HandlerInterface) *WsServer {
	this.handler = handler
	return this
}

// 运行 WsServer，监听指定端口并处理 WebSocket 连接和消息
func (this *WsServer) Run() *WsServer {
	http.HandleFunc(this.Path, this.handleWsServer)
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
func (this *WsServer) handleWsServer(w http.ResponseWriter, r *http.Request) {
	// 升级 HTTP 连接为 WsServer
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("升级失败:", err)
		return
	}
	defer conn.Close()

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

		var v data.NoDecoderMsg = theMsg
		_, WsMsg := data.Decoder(v)
		this.handler.OnMessage(conn, WsMsg)

		//this.pool.Accept(theMsg)

		//log.Printf("收到消息: %s\n", msg)

		// 原样返回消息（Echo）
		// if err := conn.WriteMessage(messageType, msg); err != nil {
		// 	log.Println("发送消息失败:", err)
		// 	break
		// }
	}
	this.handler.OnDisconnect(session)
}
