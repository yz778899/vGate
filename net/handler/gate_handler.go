package handler

import (
	"fmt"

	"github.com/14132465/vGate/net/app"
	"github.com/14132465/vGate/net/data"
	"github.com/gorilla/websocket"
)

// GateHandler网关处理器，负责处理WebSocket连接和消息
type GateHandler struct {
}

func (this *GateHandler) checkSecretKey(key string) bool {

	if !app.CheckSecretKey(key) {
		fmt.Printf("密钥不匹配，拒绝处理消息\n")
		return false
	}
	return true
}

// 收到消息
func (this *GateHandler) OnMessage(conn *websocket.Conn, msg *data.WsMsg) {

	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("处理消息时发生错误: %v\n", err)
		}
	}()

	switch msg.Cmd {
	case data.Subscription:
		//订阅消息
		if !this.checkSecretKey(msg.SecretKey) {
			return
		} else {
			server := app.ServerManager.GetServer(msg.SessionId)
			if (server) == nil {
				server = &data.Server{}
			}

			if server != nil {
				app.SubHelper.AddSubscriptionInfo(msg.Topic, server)
			} else {
				fmt.Printf("未找到会话ID %d 对应的服务器\n", msg.SessionId)
			}
		}
	case data.Publish:
		//发布消息
	case data.UnSubscription:
		//取消订阅消息
		if !this.checkSecretKey(msg.SecretKey) {
			return
		} else {
			server := app.ServerManager.GetServer(msg.SessionId)
			if server != nil {
				app.SubHelper.UnSubscriptionInfo(msg.Topic, server)
			} else {
				fmt.Printf("未找到会话ID %d 对应的服务器\n", msg.SessionId)
			}
		}
	case data.Notice:
		//通知消息
	case data.Request:
		//请求消息
		app.SubHelper.Broadcast(msg.Topic, msg)

	case data.Response:
		//转发回复消息
		session := app.SessionManager.GetSession(msg.SessionId)
		if session != nil {
			session.SendMessage(msg)
		} else {
			fmt.Printf("未找到会话ID %d 对应的客户端\n", msg.SessionId)
		}
	default:
		//fmt.Printf("未知的消息指令 %v ", msg.Cmd)

	}

}

func (this *GateHandler) OnError(conn *websocket.Conn, err error) {
	fmt.Printf("  main  ---- handler :  OnError  %v \n", err)
}

// 连接建立
func (this *GateHandler) OnConnect(conn *websocket.Conn) *data.Session {
	fmt.Printf("  main  ---- handler :  OnConnect  \n")
	// 将新连接添加到会话管理器
	session := app.SessionManager.AddSession(&data.Session{
		UUID:   -1,
		Status: 1,
		Conn:   conn,
	})
	//同时添加到了服务器列表中， 这里需要判断
	// app.ServerManager.AddServer(&data.Server{
	// 	UUID:   session.UUID,
	// 	Status: 1,
	// 	Conn:   conn,
	// })
	return session
}

// 连接断开
func (this *GateHandler) OnDisconnect(session *data.Session) {
	fmt.Printf("  main  ---- handler :  OnDisconnect  \n")
}
