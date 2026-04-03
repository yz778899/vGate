package main

import (
	"github.com/gorilla/websocket"
	"github.com/yz778899/vGate/net"
	"github.com/yz778899/vGate/net/app"
	"github.com/yz778899/vGate/net/handler"
	msg_handler "github.com/yz778899/vGate/simple/server/handler"
)

func main() {

	//注册消息管理者
	iniRegistry()

	//创建服务端
	server := net.NewWsClientAlwaysOnlie().Config("ws://localhost:8080/")

	//业务处理器
	server.Handler(&handler.ServerHandler{})
	//通信密钥
	var secretKey string = app.VGate.Config.Gate.SecretKey

	//请求连接
	server.Connect(func(conn *websocket.Conn) {
		//绑定连接，向网关订阅 登录注册等 topic
		app.Sender.BindConn(conn) // 帐号服务器 绑定连接
		app.Sender.Config(true, "Server of account", secretKey)
		//阅以下的主题消息
		app.Sender.Subscription(User_Login) //订阅用户登录命令
		app.Sender.Subscription(User_Register)
		app.Sender.Subscription(Hall_Game_List) //订阅游戏列表
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
