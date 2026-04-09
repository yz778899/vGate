package handler

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/yz778899/vGate/net/env"
	"github.com/yz778899/vGate/net/msg"
	pool "github.com/yz778899/vGate/net/poll"
	"go.uber.org/zap"
)

// ServerHandler 服务端 处理器，负责处理WebSocket连接和消息
type ServerHandler struct {
	Session *msg.Session
	Pool    *pool.QueueMaster[*pool.MessageTask]
}

var Log *zap.Logger

func init() {
	Log = env.Log
}

// 收到消息
func (this *ServerHandler) OnMessage(ctx *WebSocketContext) error {

	defer func() {
		if err := recover(); err != nil {
			Log.Error("处理消息时发生错误:", zap.Any("err", err))
		}
	}()

	//conn := ctx.Session.Conn
	wsMsg := ctx.WsMsg

	switch wsMsg.Cmd {
	// case data.Publish:
	// 	//发布消息
	case msg.Heartbeat:
		//心跳忽略
		return nil
	case msg.Notice:
		//Log.Info("### ServerHandler  cmd = Notice, 通知消息，没有订阅，也会收到的类型 ", zap.Any("topic", wsMsg.Topic))
	case msg.Request:

		//先进队列,再并行处理业务
		task := pool.MessageTask{
			//SessionId: wsMsg.SessionId,
			Msg: wsMsg,
		}
		this.Pool.Accept(&task)

	default:

		Log.Error(fmt.Sprintf("未知的消息指令 %#v \n ", wsMsg))
	}

	custom := msg.CustomMessage{
		WebsocketMsg: *wsMsg,
		HideFields:   []string{"data", "secretKey"}, // 隐藏敏感字段
	}
	jsonData, _ := json.Marshal(custom)
	Log.Info("  serverHandler :  OnMessage  %v \n", zap.String("msg", string(jsonData)))
	return nil

}

func (this *ServerHandler) OnError(conn *websocket.Conn, err error) {
	Log.Error(fmt.Sprintf("  serverHandler :  OnError  %v \n", err))
}

// 连接建立
func (this *ServerHandler) OnConnect(conn *websocket.Conn) *msg.Session {
	// 将新连接添加到会话管理器
	session := env.VGate.SessionMgr.AddSession(&msg.Session{
		ConnSession: &msg.ConnSession{
			Conn: conn,
			UUID: -1,
		},
		Status: 1,
	})
	if session == nil {
		fmt.Print("debug break!")
	} else {

		Log.Info(fmt.Sprintf("  serverHandler :  OnConnect session = %#v ", session))
	}
	return session
}

// 连接断开
func (this *ServerHandler) OnDisconnect(session *msg.Session) {
	Log.Error(fmt.Sprintf("  serverHandler :  OnDisconnect session = %#v ", session))
}
