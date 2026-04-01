package handler

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/14132465/vGate/net/app"
	"github.com/14132465/vGate/net/data"
	"github.com/14132465/vGate/net/logic"
	"github.com/gofiber/fiber/v2/log"
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
func (this *GateHandler) OnMessage(conn *websocket.Conn, msg *data.WsMsg) error {

	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("处理消息时发生错误: %v\n", err)
		}
	}()

	conn.SetReadDeadline(time.Now().Add(time.Duration(app.VGate.Config.Gate.ReadOverTime) * time.Second))

	fmt.Printf("  GateHandler :  OnMessage  %#v \n", msg)

	switch msg.Cmd {
	case data.Heartbeat:
		//心跳
		err := conn.WriteJSON(msg)
		if err != nil {
			log.Error("SendMessage Heartbeat error %v \n", err)
		}
		return nil
	case data.Subscription:
		//订阅消息
		if !this.checkSecretKey(msg.SecretKey) {
			return nil
		} else {
			server := app.VGate.ServerManager.GetAndCreateServer(msg.SessionId)
			if server != nil {
				app.VGate.SubHelper.AddSubscriptionInfo(msg.Topic, server)
			} else {
				log.Error("未找到会话ID %d 对应的服务器\n", msg.SessionId)
			}
		}
		return nil
	// case data.Publish:
	// 	//发布消息
	case data.UnSubscription:
		//取消订阅消息
		if !this.checkSecretKey(msg.SecretKey) {
			return nil
		} else {
			server := app.VGate.ServerManager.GetServerOnly(msg.SessionId)
			if server != nil {
				app.VGate.SubHelper.UnSubscriptionInfo(msg.Topic, server)
			} else {
				log.Error("未找到会话ID %d 对应的服务器\n", msg.SessionId)
			}
		}
		return nil
	case data.Notice:
		//通知消息
		isHandler := logic.NoticeHelperInstance.Handler(msg)
		if !isHandler {
			//其它通知，转发给所有的服务器
			app.VGate.SubHelper.Broadcast(msg.Topic, msg)
		}
		return nil
	case data.Request:
		//客户端请求消息，将通过订阅信息管理器转发给订阅了指定主题的服务器
		app.VGate.SubHelper.Broadcast(msg.Topic, msg)
		return nil
	case data.Response:
		//转发回复消息
		session := app.VGate.SessionManager.GetSession(msg.SessionId)
		if session != nil {
			session.SendMessage(msg)
		} else {
			log.Error("未找到会话ID %d 对应的客户端\n", msg.SessionId)
		}
		return nil
	default:
		fmt.Printf("未知的消息指令 %#v ", msg)
		//fmt.Printf("  GateHandler :  OnMessage  %v \n", msg)
	}

	return nil
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
	noticeMsg := data.BuildNoticeMsg(app.VGate.Config.Gate.SecretKey, logic.Notice_On_Line, by)
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
		app.VGate.ServerManager.RemoveServer(session.UUID)
	} else {
		//通知客户端下线
		lst := app.VGate.ServerManager.GetAlls()
		by, _ := json.Marshal(session)
		noticeMsg := data.BuildNoticeMsg(app.VGate.Config.Gate.SecretKey, logic.Notice_Off_Line, by)
		for _, server := range lst {
			if server != nil {
				server.SendMessage(noticeMsg)
			} else {
				fmt.Print("break point of debug !")
			}
		}

		//logic.Sender.Response(topic string, msg *data.WsMsg)
	}

}
