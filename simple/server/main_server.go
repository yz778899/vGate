package main

import (
	"github.com/14132465/vGate/net"
	"github.com/14132465/vGate/net/app"
	"github.com/14132465/vGate/net/handler"
	msg_handler "github.com/14132465/vGate/simple/server/handler"
	"github.com/gorilla/websocket"
)

func main() {

	//注册消息管理者
	iniRegistry()

	//创建服务端
	server := net.NewWsClient().Config("ws://localhost:8080/")

	//业务处理器
	server.Handler(&handler.ServerHandler{})
	//通信密钥
	var secretKey string = app.VGate.Config.Gate.SecretKey

	//请求连接
	server.Connect(func(conn *websocket.Conn) {
		//连接成功后，就订阅以下的主题消息

		//绑定连接，向网关订阅 登录注册等 topic
		sender := app.Sender
		sender.BindConn(conn) // 帐号服务器 绑定连接
		sender.Config(true, "Server of account", secretKey)
		//"ga-23xk=v"
		sender.Subscription(User_Login) //订阅用户登录命令
		sender.Subscription(User_Register)
		sender.Subscription(Hall_Game_List) //订阅游戏列表

	})

}

func iniRegistry() {
	registry := handler.NewRegistry()

	topic := "/user/login"
	registry.Register(handler.MsgHandlerCreate{
		Topic:      topic,
		CreateFunc: msg_handler.NewLoginHandler,
	})

}
