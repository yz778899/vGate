package main

import (
	"github.com/14132465/vGate/net"
	"github.com/14132465/vGate/net/handler"
	"github.com/14132465/vGate/net/logic"
	"github.com/gorilla/websocket"
)

func main() {

	//创建服务端
	server := net.NewWsClient().Config("ws://localhost:8080/")

	//业务处理器
	server.Handler(&handler.ServerHandler{})

	var secretKey string = "ga-23xk=v"

	//请求连接
	server.Connect(func(conn *websocket.Conn) {
		//连接成功后，就订阅以下的主题消息
		sender := logic.Sender
		sender.BindConn(conn) // 帐号服务器 绑定连接
		sender.Config(true, "Server of account", secretKey)
		//"ga-23xk=v"
		sender.Subscription(User_Login) //订阅用户登录命令
		sender.Subscription(User_Register)
		sender.Subscription(Hall_Game_List) //订阅游戏列表

	})

}
