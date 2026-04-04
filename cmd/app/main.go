package main

import (
	"github.com/gorilla/websocket"
	msg_handler "github.com/yz778899/vGate/cmd/app/handler"
	"github.com/yz778899/vGate/net"
	"github.com/yz778899/vGate/net/handler"
	"github.com/yz778899/vGate/net/logic"
)

func main() {
	//注册消息管理者
	iniRegistry()
	//创建服务端
	app := net.NewAppService().Config("ws://localhost:8080/")
	//业务处理器
	app.Handler(&handler.ServerHandler{})
	//请求连接
	app.Connect(func(conn *websocket.Conn) {
		//绑定连接，向网关订阅 登录注册等 topic
		logic.Sender.BindConn(conn) // 帐号服务器 绑定连接
		logic.Sender.Config(true, "Server of account")
		//阅以下的主题消息
		logic.Sender.Subscription(User_Login) //订阅用户登录命令
		logic.Sender.Subscription(User_Register)
		logic.Sender.Subscription(Hall_Game_List) //订阅游戏列表
	})
}

// 注册处理器
func iniRegistry() {
	registry := handler.NewRegistry()
	registry.Register(handler.MsgHandlerCreate{
		Topic:      "/user/login",
		CreateFunc: msg_handler.NewLoginHandler,
	})
}
