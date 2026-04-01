package handler

import (
	"fmt"

	"github.com/14132465/vGate/net/app"
	"github.com/14132465/vGate/net/data"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gorilla/websocket"
)

// ServerHandler 服务端 处理器，负责处理WebSocket连接和消息
type ServerHandler struct {
	Session *data.Session
}

// 收到消息
func (this *ServerHandler) OnMessage(conn *websocket.Conn, msg *data.WsMsg) error {

	defer func() {
		if err := recover(); err != nil {
			log.Error(fmt.Printf("处理消息时发生错误: %v\n", err))
		}
	}()

	switch msg.Cmd {
	// case data.Publish:
	// 	//发布消息
	case data.Heartbeat:
		//心跳忽略
		return nil
	case data.Notice:
		log.Info(fmt.Printf("### ServerHandler  cmd = Notice, Topic = %v 通知消息，没有订阅，也会收到的类型 \n", msg.Topic))
	case data.Request:
		//客户端请求消息
		by, err := msg.Content.MarshalJSON()
		if err != nil {
			return err
		} else {

			log.Info(fmt.Printf(" recv  topic = %v msg = %v \n", msg.Topic, string(by)))

			RegistryInstance.RunHandler(msg, this.Session)

			// creator, ok := RegistryInstance.GetMsgHandlerCreate(msg.Topic)
			// if ok {
			// 	creator.CreateFunc(msg.Topic, &data.Session{}, &data.WsMsg{}).Init()
			// } else {
			// 	log.Error("### 缺少 MsgHandler 对应 topic 是 " + msg.Topic + " ，该消息将丢弃处理 !\n")
			// }
		}
	default:
		log.Error(fmt.Sprintf("未知的消息指令 %#v \n ", msg))
	}
	log.Info(fmt.Sprintf("	ServerHandler  OnMessage  msg = %#v  \n", msg))
	return nil

}

func (this *ServerHandler) OnError(conn *websocket.Conn, err error) {
	log.Error(fmt.Sprintf("  serverHandler :  OnError  %v \n", err))
}

// 连接建立
func (this *ServerHandler) OnConnect(conn *websocket.Conn) *data.Session {
	// 将新连接添加到会话管理器
	session := app.VGate.SessionManager.AddSession(&data.Session{
		UUID:   -1,
		Status: 1,
		Conn:   conn,
	})
	log.Info(fmt.Sprintf("  serverHandler :  OnConnect session = %#v ", session))
	return session
}

// 连接断开
func (this *ServerHandler) OnDisconnect(session *data.Session) {
	log.Error(fmt.Sprintf("  serverHandler :  OnDisconnect session = %#v ", session))
}
