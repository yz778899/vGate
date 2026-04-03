package handler

import (
	"github.com/gorilla/websocket"
	"github.com/yz778899/vGate/net/data"
)

// WsHandlerInterface 定义了处理WebSocket连接和消息的接口
type WsHandlerInterface interface {
	// 收到消息
	OnMessage(conn *websocket.Conn, msg *data.WsMsg) error
	// 连接建立
	OnConnect(conn *websocket.Conn) *data.Session
	// 连接断开
	OnDisconnect(session *data.Session)
}

// MsgHandlerInterface 单条消息处理接口
type MsgHandlerInterface interface {
	GetTopic() string // 处理器对应的主题
	// 核心生命周期方法
	Init() error          // 初始化
	BeforeProcess() error // 处理前
	Process() error       // 处理中
	AfterProcess()        // 处理后
	Release() error       // 释放
	// 辅助方法
	OnError(stage string, err error) // 错误处理钩子
}

type MsgHandlerCreate struct {
	Topic string
	//创建实例的方法
	CreateFunc func(topic string, session *data.Session, msg *data.WsMsg) MsgHandlerInterface
}

//type HandlerFunc func(session *data.Session, msg *data.WsMsg)
