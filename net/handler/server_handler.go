package handler

import (
	"fmt"

	"github.com/14132465/vGate/net/app"
	"github.com/14132465/vGate/net/data"
	"github.com/gorilla/websocket"
)

// ServerHandler 服务端 处理器，负责处理WebSocket连接和消息
type ServerHandler struct {
}

// 收到消息
func (this *ServerHandler) OnMessage(conn *websocket.Conn, msg *data.WsMsg) {

	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("处理消息时发生错误: %v\n", err)
		}
	}()

	switch msg.Cmd {
	// case data.Publish:
	// 	//发布消息
	case data.Notice:
		fmt.Printf("### ServerHandler  cmd = Notice, Topic = %v 通知消息，没有订阅，也会收到的类型 \n", msg.Topic)
	case data.Request:
		//客户端请求消息
		//
	default:
		fmt.Printf("未知的消息指令 %v ", msg.Cmd)

	}

	fmt.Printf("	ServerHandler  OnMessage  msg = %#v  \n", msg)

}

func (this *ServerHandler) OnError(conn *websocket.Conn, err error) {
	fmt.Printf("  serverHandler :  OnError  %v \n", err)
}

// 连接建立
func (this *ServerHandler) OnConnect(conn *websocket.Conn) *data.Session {
	// 将新连接添加到会话管理器
	session := app.VGate.SessionManager.AddSession(&data.Session{
		UUID:   -1,
		Status: 1,
		Conn:   conn,
	})
	fmt.Printf("  serverHandler :  OnConnect session = %#v \n", session)
	return session
}

// 连接断开
func (this *ServerHandler) OnDisconnect(session *data.Session) {
	fmt.Printf("  serverHandler :  OnDisconnect session = %#v \n", session)
}
