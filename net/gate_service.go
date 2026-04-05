package net

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/yz778899/vGate/net/env"
	"github.com/yz778899/vGate/net/env/config"
	"github.com/yz778899/vGate/net/handler"
	data "github.com/yz778899/vGate/net/msg"

	ws "github.com/gorilla/websocket"
)

// GateServer结构体表示一个WebSocket服务器，包含端口、路径、协程池和消息处理器等信息
type GateServer struct {
	conf *config.GateConfig
	// Port    string
	// Path    string
	handler handler.ServiceAcceptInterface
}

// 配置 WsServer 的端口和路径
func (this *GateServer) WithConfig(conf *config.GateConfig) *GateServer {
	this.conf = conf
	return this
}

// 创建 WsServer 实例
func NewWsServer() *GateServer {
	ws := GateServer{
		conf:    &env.VGate.Config.Gate,
		handler: &handler.GateHandler{},
	}
	return &ws
}

// 配置消息处理器
func (this *GateServer) Handler(handler handler.ServiceAcceptInterface) *GateServer {
	this.handler = handler
	return this
}

// 运行 WsServer，监听指定端口并处理 WebSocket 连接和消息
func (this *GateServer) Run() error {
	// if this.handler == nil {
	// 	this.handler = &handler.GateHandler{}
	// }
	// if this.Conf == nil {
	// 	this.Conf = &env.VGate.Config.Gate
	// }
	http.HandleFunc(this.conf.WsPath, this.wsServerHandler)
	env.Log.Info(fmt.Sprintf("WsServer run , port : %v", this.conf.WsPort))
	return http.ListenAndServe(fmt.Sprintf(":%v", this.conf.WsPort), nil)
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
func (this *GateServer) wsServerHandler(w http.ResponseWriter, r *http.Request) {
	// 升级 HTTP 连接为 WsServer
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		env.Log.Info(fmt.Sprintf("升级失败: %v", err))
		return
	}
	defer conn.Close()

	//读超时 10秒
	conn.SetReadDeadline(time.Now().Add(time.Duration(this.conf.ReadOverTime) * time.Second))

	session := this.handler.OnConnect(conn)

	for {
		// 读取客户端消息
		_, originalMsgByteArray, err := conn.ReadMessage()
		if err != nil {
			log.Println("关闭通道:", err)
			this.handler.OnDisconnect(session)
			break
		}

		//fmt.Printf("网关收到消息 msg = %v \n", string(msg))

		var theMsg data.NoDecoderMsg

		theMsg = data.NoDecoderMsg{
			SessionId: session.UUID,
			Msg:       string(originalMsgByteArray),
			SnId:      rand.Intn(1024),
		}

		WsMsg, _ := data.GateDecoder(theMsg)

		//fmt.Printf("data.GateDecoder 处理后的消息 msg = %#v \n", WsMsg)
		this.handler.OnMessage(handler.WebSocketContext{
			Session:  session,
			Original: &originalMsgByteArray,
			WsMsg:    WsMsg,
		})

	}
	this.handler.OnDisconnect(session)
}
