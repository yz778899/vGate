## vGate 

一个用go编写的网关组件，通信目前支持 websocket 。socket 和 udp 版后续将会加入

该项目为游戏分布式、集群而设计，也可用于传统软件的web项目，目前是个人开源项目，没有使用限制。

##  使用示例



### 网关示例
```go

func main() {
	defer env.Log.Sync()

	handler := handler.GateHandler{}
	err := net.NewWsServer().Config(8080, "/").Handler(&handler).Run()
	if err != nil {
		env.Log.Fatal("Server failed to start: ", zap.Error(err))
	}
}
```
### 服务端示例
```go


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

```

本项目按照三个职能来设计
-  **网关：**

   维护服务端的订阅数据，可Subscription（订阅） UnSubscription(取消订阅)

   路由消息，功能有如下几点

      1 将客户端的请求转发到服务端

      2 将服务端的响应转发给客户端

      3 针对 Notice （通知）类型的消息，对所有服务端广播，可由网关发起，也可以任一个服务端发起。

-  **服务端：**


   1 可向网关 Subscription（订阅） UnSubscription(取消订阅) 。 网关会将topic相符的消息发布给服务

   2 同一个 topic 有多个服务端 订阅，多个服务端都会收到该 topic所属的消息

   3 将业务处理的结果，通知到指定的客户端

   4 断线重连 、离线订阅 待开发






- **客户端**

   1 仅向网关发送请求，即可得到服务端的响应。服务端的端口不对外暴露。

   2 多个服务器共用一个网关端口，无感切换服务，一次单个连接处理所有业务。


测试使用第三方 websocket http://www.websocket-test.com/ 输入  ws://localhost:8080  点连接

网关与服务端启动后，用下面的测试消息测试

- ** 登录 

```bash
{"cmd":"request","topic":"/user/login","content":{"user":"jack" , "pass":"123456"}}
```


- ** 游戏列表

```bash

{"cmd":"request","topic":"/hall/game_list","content":{}

```

同时也可以用第三方 websocket 来当服务器用，订阅消息。

作为服务端订阅消息 ( secretKey 必须与网关一致，否则网关将会拒绝订阅 )

```bash
{
    "cmd": "subscription",
    "topic": "/user/login",
    "serverName": "Server of account",
    "secretKey": "sdklPY#$xks-23ksd%^dfskljkl[@#345]"
}
```

以上请求将成功订阅 /user/login topic ，当有其它客户端使用 下面的请求时将会收到它的消息

```bash
{"cmd":"request","topic":"/user/login","content":{"user":"jack" , "pass":"123456"}}
```
