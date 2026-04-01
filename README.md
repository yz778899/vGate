# vGate 

一个用go编写的网关组件，通信目前仅支持 websocket 。socket 和 udp 版后续将会加入

该项目为游戏分布式、集群而设计，也可用于传统软件的web项目，目前是个人开源项目，可修改，商用，没有限制。

本项目按照三个职能来设计

- **网关：
   维护服务端的订阅数据，可Subscription（订阅） UnSubscription(取消订阅)

   路由消息，功能有如下几点

      1 将客户端的请求转发到服务端

      2 将服务端的响应转发给客户端

      3 针对 Notice （通知）类型的消息，对所有服务端广播，可由网关发起，也可以任一个服务端发起。

- **服务端：

   1 可向网关 Subscription（订阅） UnSubscription(取消订阅)，订阅的topic ，网关即会将topic相符的消息发布，从而实现服务分布式与集群

   2 同一个 topic 有多个服务端 订阅，多个服务端都会收到该 topic所属的消息

   3 将业务处理的结果，通知到指定的客户端

   4 断线重连 、离线订阅 待开发

- **客户端 
   1 仅向网关发送请求，即可得到服务端的响应。服务端的端口不对外暴露。
   2 多个服务器共用一个网关端口，无感切换服务，一次单个连接处理所有业务。

   
## 技术栈

- **Websocket 框架**：gorilla/websocket

## 项目结构

## 测试代码

- ** 登录
{"cmd":"request","topic":"/user/login","content":{"user":"jack" , "pass":"123456"}}

- ** 游戏列表
{"cmd":"request","topic":"/hall/game_list","content":{}

