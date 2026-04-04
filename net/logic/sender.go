package logic

import (
	"encoding/json"
	"errors"
	"fmt"
	"runtime/debug"

	"github.com/gofiber/fiber/v2/log"
	"github.com/gorilla/websocket"
	"github.com/yz778899/vGate/net/data"
	"github.com/yz778899/vGate/net/env"
	"github.com/yz778899/vGate/net/env/config"
)

// 消息发送者 - 服务端|客户端
type sender struct {
	Conf       *config.RootConfig
	Conn       *websocket.Conn
	serverName string
	isServer   bool   //是否为服务端
	SecretKey  string //密钥
}

// 消息发送者 网关
type SenderForGate struct {
	data.Session
}

// 用于 服务器|客户端 的消息发送者
var Sender *sender

func init() {
	Sender = &sender{}
}

// 绑定
func (this *sender) BindConn(connInstance *websocket.Conn) *sender {
	this.Conn = connInstance
	return this
}

// 绑定
// secretKey 安全密钥
func (this *sender) Config(isServer bool, serverName string) *sender {
	this.isServer = isServer
	this.serverName = serverName
	this.Conf = env.VGate.Config
	return this
	//config env.VGate.Config
}

//secretKey

// 通知所有服务端
func (this *sender) Notice(topic string, msg any) error {
	if !this.isServer {
		return errors.New("客户端不能使用该方法!")
	}
	return this.sendMsg(0, data.Notice, topic, msg)
}

// 响应 request
// func (this *sender) Response(userId int64, msg data.BaseMsgInterFace) error {
// 	if !this.isServer {
// 		return errors.New("客户端不能使用该方法!")
// 	}
// 	return this.sendMsg(userId, data.Response, msg.GetTopic(), msg.GetContent())
// }

// 响应 request
func (this *sender) Resp(userId int64, topic string, msg any) error {
	if !this.isServer {
		return errors.New("客户端不能使用该方法!")
	}
	return this.sendMsg(userId, data.Response, topic, msg)
}

// 请求,
func (this *sender) Request(userId int64, topic string, msg any) error {

	if this.isServer {
		return errors.New("服务器不能使用该方法!")
	}
	return this.sendMsg(userId, data.Request, topic, msg)
}

// 发送消息到网关
func (this *sender) sendMsg(userId int64, cmd string, topic string, msg any) error {
	defer func() {
		if err := recover(); err != nil {
			log.Error(fmt.Printf(" panic: %v\n", err))
			log.Error(fmt.Printf(" Stack Info:\n %s \n", debug.Stack()))
		}
	}()
	if this.Conn == nil {
		return errors.New("需要绑定连接, 方法 BindConn(conn websocket.Conn) !")
	} else {

		content, err := json.Marshal(msg)
		if err != nil {
			return err
		}
		//连接可用
		var sendMsg any
		switch cmd {
		case data.Notice:
			sendMsg = data.BuildNoticeMsg(this.Conf.Gate.SecretKey, topic, content)
		case data.Request:
			sendMsg = data.BuildRequestMsg(userId, topic, content)
		case data.Response:
			sendMsg = data.BuildResponseMsg(userId, topic, content)
		default:

			log.Error("要发送的消息类型 在意料中外 ， 将会丢弃消息")
			return nil
		}
		err = this.Conn.WriteJSON(sendMsg)
		if err != nil {
			log.Error("SendMessage  error  %v \n", err)
		}

	}
	return nil
}

// 网关订阅
func (this *sender) Subscription(topic string) error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf(" panic: %v\n", err)
			fmt.Printf(" Stack Info:\n %s \n", debug.Stack())
		}
	}()
	if !this.isServer {
		return errors.New("客户端不能使用该方法!")
	}
	if this.Conn == nil {
		return errors.New("需要绑定连接, 方法 BindConn(conn websocket.Conn) !")
	} else {
		//连接可用
		msg := data.BuildSubscriptionMsg(topic, this.serverName, this.SecretKey)
		err := this.Conn.WriteJSON(msg)
		if err != nil {
			fmt.Printf("SendMessage  error %v \n", err)
		}

	}
	return nil
}

// 网关取消订阅
func (this *sender) UnSubscription(topic string) error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf(" panic: %v\n", err)
			fmt.Printf(" Stack Info:\n %s \n", debug.Stack())
		}
	}()

	if !this.isServer {
		return errors.New("客户端不能使用该方法!")
	}
	if this.Conn == nil {
		return errors.New("需要绑定连接, 方法 BindConn(conn websocket.Conn) !")
	} else {
		//连接可用
		msg := data.BuildUnSubscriptionMsg(topic, this.serverName)
		// by, err := json.Marshal(msg)
		// if err != nil {
		// 	return err
		// }
		// this.Conn.WriteJSON(by)

		err := this.Conn.WriteJSON(msg)
		if err != nil {
			fmt.Printf("SendMessage  error %v \n", err)
		}

	}
	return nil
}
