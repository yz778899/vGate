package handler

import (
	"github.com/14132465/vGate/net/data"
	"github.com/gorilla/websocket"
)

// HandlerInterface定义了处理WebSocket连接和消息的接口
type HandlerInterface interface {
	// 收到消息
	OnMessage(conn *websocket.Conn, msg *data.WsMsg)
	// 连接建立
	OnConnect(conn *websocket.Conn) *data.Session
	// 连接断开
	OnDisconnect(session *data.Session)
}
