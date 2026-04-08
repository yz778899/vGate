package handler

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"github.com/yz778899/vGate/net/env"
	"github.com/yz778899/vGate/net/logic"
	"github.com/yz778899/vGate/net/msg"
	data "github.com/yz778899/vGate/net/msg"
	"go.uber.org/zap"
)

// GateHandler网关处理器，负责处理WebSocket连接和消息
type GateHandler struct {
}

func (this *GateHandler) checkSecretKey(key string) bool {

	if !env.VGate.CheckSecretKey(key) {
		Log.Error("密钥不匹配，拒绝处理消息")
		return false
	}
	return true
}

// 收到消息
func (this *GateHandler) OnMessage(ctx *WebSocketContext) error {

	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("处理消息时发生错误: %v\n", err)
		}
	}()
	conn := ctx.Session.Conn
	msg := ctx.WsMsg

	conn.SetReadDeadline(time.Now().Add(time.Duration(env.VGate.Config.Gate.ReadOverTime) * time.Second))

	// custom := data.CustomMessage{
	// 	WebsocketMsg: *msg,
	// 	HideFields:   []string{"data", "secretKey"}, // 隐藏敏感字段
	// }
	//jsonData, _ := json.Marshal(custom)
	//jsonData, _ := json.MarshalIndent(custom, "", "  ")

	//fmt.Printf("  GateHandler :  OnMessage  %v \n", string(jsonData))

	switch msg.Cmd {
	case data.Heartbeat:
		//心跳
		err := conn.WriteJSON(msg)
		if err != nil {
			Log.Error("SendMessage Heartbeat error  \n", zap.Any("err", err))
		}
		return nil
	case data.Subscription:
		//订阅消息
		if !this.checkSecretKey(msg.SecretKey) {
			return nil
		} else {
			server := env.VGate.AppSessionMgr.GetAndCreateServer(msg.SessionId)
			if server != nil {
				logic.SubHelper.AddSubscriptionInfo(msg.Topic, server)
			} else {
				Log.Error("A未找到会话ID  对应的服务器\n", zap.Any("msg.SessionId", msg.SessionId), zap.Any("msg", msg))
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
			server := env.VGate.AppSessionMgr.GetServerOnly(msg.SessionId)
			if server != nil {
				logic.SubHelper.UnSubscriptionInfo(msg.Topic, server)
			} else {
				Log.Error("未找到会话ID  对应的服务器\n", zap.Any("sessionId", msg.SessionId), zap.Any("msg", msg))
			}
		}
		return nil
	case data.Notice:
		//通知消息
		isHandler := logic.NoticeHelperInstance.Handler(msg)
		if !isHandler {
			//其它通知，转发给所有的服务器
			logic.SubHelper.Broadcast(msg.Topic, msg)
		}
		return nil
	case data.Request:
		//客户端请求消息，将通过订阅信息管理器转发给订阅了指定主题的服务器
		logic.SubHelper.Broadcast(msg.Topic, msg)
		return nil
	case data.Response:
		//转发回复消息
		session := env.VGate.SessionMgr.GetSession(msg.SessionId)
		if session != nil {
			toClient := data.ToClientMsg{}
			session.SendToClient(toClient.TransitionOf(msg))
		} else {
			Log.Error("未找到会话ID %d 对应的客户端\n", zap.Any("msg.SessionId", msg.SessionId), zap.Any("msg", msg))
		}
		return nil
	default:
		Log.Error("收到未知的消息 %#v ", zap.Any("msg.data", msg.Data))
		session := env.VGate.SessionMgr.GetSession(msg.SessionId)
		if session != nil {
			session.SendToClient(data.GetUnknownMsg("未知消息，请按规范请求"))
		}
	}

	return nil
}

func (this *GateHandler) OnError(conn *websocket.Conn, err error) {
	Log.Error("  GateHandler :  OnError  %v \n", zap.Any("err", err))
}

// 连接建立
func (this *GateHandler) OnConnect(conn *websocket.Conn) *data.Session {
	// 将新连接添加到会话管理器
	session := env.VGate.SessionMgr.AddSession(&data.Session{
		ConnSession: &msg.ConnSession{
			Conn: conn,
			UUID: -1,
		},
		Status: 1,
	})

	Log.Debug("  GateHandler :  new Connect sessionId = ", zap.Any("session.UUID", session.UUID))

	//通知客户端上线
	lst := env.VGate.AppSessionMgr.GetAlls()
	by, _ := json.Marshal(session)
	noticeMsg := data.BuildNoticeMsg(env.VGate.Config.Gate.SecretKey, logic.Notice_On_Line, by)
	for _, server := range lst {
		if server != nil {
			server.SendMessage(noticeMsg)
		}
	}
	return session
}

// 连接断开
func (this *GateHandler) OnDisconnect(session *data.Session) {
	//fmt.Printf("  GateHandler :  OnDisconnect session = %#v \n", session)
	env.VGate.SessionMgr.RemoveSession(session.UUID)
	server := env.VGate.AppSessionMgr.GetServerOnly(session.UUID)
	if server != nil {

		logic.SubHelper.ServerClose(server)
		env.VGate.AppSessionMgr.RemoveServer(session.UUID)
	} else {
		//通知客户端下线
		lst := env.VGate.AppSessionMgr.GetAlls()
		by, _ := json.Marshal(session)
		noticeMsg := data.BuildNoticeMsg(env.VGate.Config.Gate.SecretKey, logic.Notice_Off_Line, by)
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
