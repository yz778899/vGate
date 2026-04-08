package main

import (
	"runtime"

	"github.com/gorilla/websocket"
	msg_handler "github.com/yz778899/vGate/cmd/app/handler"
	"github.com/yz778899/vGate/net"
	_ "github.com/yz778899/vGate/net/env"
	"github.com/yz778899/vGate/net/handler"
	"github.com/yz778899/vGate/net/logic"
	pool "github.com/yz778899/vGate/net/poll"
)

var AppService *net.AppService

func main() {
	//注册消息管理者
	iniRegistry()

	//创建服务端  host.docker.internal = docker的宿主机
	app := net.NewAppService().Config("ws://host.docker.internal:5566/", runtime.NumCPU()*2)
	//app := net.NewAppService().Config("ws://localhost:5566/", runtime.NumCPU()*2)
	AppService = app
	//业务处理器
	app.Handler(&handler.ServerHandler{Pool: app.Pool})
	//服务端收到消息,初步过滤后.即会将消息组装成任务,进入任务队列,然后再并行处理,
	app.Pool.Handler(func(task *pool.MessageTask) {
		//此处是设置任务子线程的具体处理方法
		handler.RegistryInstance.RunHandler(task.Msg, AppService.Session)
	})
	//请求连接
	app.Connect(func(conn *websocket.Conn) {
		//绑定连接，向网关订阅 登录注册等 topic

		logic.Sender.BindConn(app.Session) // 帐号服务器 绑定连接
		logic.Sender.Config(true, "AccountApp")
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
