package main

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2/log"
	"github.com/gorilla/websocket"
	appmsg "github.com/yz778899/vGate/cmd/app/app_msg"
	"github.com/yz778899/vGate/net/handler"
	"github.com/yz778899/vGate/net/msg"
)

// ClientHandler 服务端 处理器，负责处理WebSocket连接和消息
type ClientHandler struct {
	Session *msg.Session
}

// 收到消息
func (this *ClientHandler) OnMessage(ctx *handler.WebSocketContext) error {

	defer func() {
		if err := recover(); err != nil {
			log.Error(fmt.Printf("处理消息时发生错误: %v\n", err))
		}
	}()

	//收到消息,0.5秒继续发送 登录消息
	//time.Sleep(time.Millisecond * 1)
	this.SendToGate(ctx.Session, "/user/login", LoginMsg())
	return nil
}

// 发送消息[到网关]
func (this *ClientHandler) SendToGate(session *msg.Session, Topic string, _msg any) error {

	msgByteArray, err := json.Marshal(_msg)
	if err != nil {
		log.Error("Json解析出错, ", err)
		return err
	}
	pack := msg.ToClientMsg{
		Topic: Topic,
		Data:  msgByteArray,
	}
	err = session.Conn.WriteJSON(pack)
	if err != nil {
		log.Error("client handler SendMessage  error %v \n", err)
	}
	return err
}

func LoginMsg() any {
	//登录的消息
	loginMsg := appmsg.LoginRequest{
		User: "asdfsf",
		Pass: "sdf",
	}
	return &loginMsg
}

func (this *ClientHandler) OnError(conn *websocket.Conn, err error) {
	log.Error(fmt.Sprintf("  ClientHandler :  OnError  %v \n", err))
}

// 连接建立
func (this *ClientHandler) OnConnect(conn *websocket.Conn) *msg.Session {

	return &msg.Session{
		ConnSession: &msg.ConnSession{
			Conn: conn,
			UUID: -1,
		},
	}
}

// 连接断开
func (this *ClientHandler) OnDisconnect(session *msg.Session) {
	log.Error(fmt.Sprintf("  ClientHandler :  OnDisconnect session = %#v ", session))
}
