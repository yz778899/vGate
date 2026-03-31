package handler

import (
	"encoding/json"
	"fmt"

	"github.com/14132465/vGate/net/app"
	"github.com/14132465/vGate/net/data"
	"github.com/14132465/vGate/net/logic"
	"github.com/gorilla/websocket"
)

// GateHandler网关处理器，负责处理WebSocket连接和消息
type GateHandler struct {
}

func (this *GateHandler) checkSecretKey(key string) bool {

	if !app.VGate.CheckSecretKey(key) {
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
			server := app.VGate.ServerManager.GetAndCreateServer(msg.SessionId)
			if server != nil {
				app.VGate.SubHelper.AddSubscriptionInfo(msg.Topic, server)
			} else {
				fmt.Printf("未找到会话ID %d 对应的服务器\n", msg.SessionId)
			}
		}
	// case data.Publish:
	// 	//发布消息
	case data.UnSubscription:
		//取消订阅消息
		if !this.checkSecretKey(msg.SecretKey) {
			return
		} else {
			server := app.VGate.ServerManager.GetServerOnly(msg.SessionId)
			if server != nil {
				app.VGate.SubHelper.UnSubscriptionInfo(msg.Topic, server)
			} else {
				fmt.Printf("未找到会话ID %d 对应的服务器\n", msg.SessionId)
			}
		}
	case data.Notice:
		//通知消息
	case data.Request:
		//客户端请求消息，将通过订阅信息管理器转发给订阅了指定主题的服务器
		app.VGate.SubHelper.Broadcast(msg.Topic, msg)

	case data.Response:
		//转发回复消息
		session := app.VGate.SessionManager.GetSession(msg.SessionId)
		if session != nil {
			session.SendMessage(msg)
		} else {
			fmt.Printf("未找到会话ID %d 对应的客户端\n", msg.SessionId)
		}
	default:
		//fmt.Printf("未知的消息指令 %v ", msg.Cmd)

	}

	fmt.Printf("  GateHandler :  OnMessage  %v \n", msg)

}

func (this *GateHandler) OnError(conn *websocket.Conn, err error) {
	fmt.Printf("  GateHandler :  OnError  %v \n", err)
}

// 连接建立
func (this *GateHandler) OnConnect(conn *websocket.Conn) *data.Session {
	// 将新连接添加到会话管理器
	session := app.VGate.SessionManager.AddSession(&data.Session{
		UUID:   -1,
		Status: 1,
		Conn:   conn,
	})
	fmt.Printf("  GateHandler :  OnConnect sessionId = %#v \n", session.UUID)

	//通知客户端上线
	lst := app.VGate.ServerManager.GetAlls()
	by, _ := json.Marshal(session)
	noticeMsg := data.BuildNoticeMsg(app.VGate.SecretKey, logic.Notice_On_Line, string(by))
	for _, server := range lst {
		if server != nil {
			server.SendMessage(noticeMsg)
		}
	}
	return session
}

// 连接断开
func (this *GateHandler) OnDisconnect(session *data.Session) {
	fmt.Printf("  GateHandler :  OnDisconnect session = %#v \n", session)
	app.VGate.SessionManager.RemoveSession(session.UUID)
	server := app.VGate.ServerManager.GetServerOnly(session.UUID)
	if server != nil {

		logic.SubHelper.ServerClose(server)
	} else {
		//通知客户端下线
		lst := app.VGate.ServerManager.GetAlls()
		for _, server := range lst {
			by, _ := json.Marshal(session)
			noticeMsg := data.BuildNoticeMsg(app.VGate.SecretKey, logic.Notice_Off_Line, string(by))
			server.SendMessage(noticeMsg)
		}

		//logic.Sender.Response(topic string, msg *data.WsMsg)
	}

}
