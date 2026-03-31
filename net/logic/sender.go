package logic

import (
	"encoding/json"
	"errors"
	"fmt"
	"runtime/debug"

	"github.com/14132465/vGate/net/data"
	"github.com/gorilla/websocket"
)

// 消息发送者 - 服务端|客户端
type sender struct {
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
func (this *sender) Config(isServer bool, serverName string, secretKey string) *sender {
	this.isServer = isServer
	this.serverName = serverName
	this.SecretKey = secretKey
	return this
}

//secretKey

// 通知所有服务端
func (this *sender) Notice(topic string, msg *data.WsMsg) error {
	if !this.isServer {
		return errors.New("客户端不能使用该方法!")
	}
	return this.sendMsg(data.Notice, topic, msg)
}

// 响应
func (this *sender) Response(topic string, msg *data.WsMsg) error {
	if !this.isServer {
		return errors.New("客户端不能使用该方法!")
	}
	return this.sendMsg(data.Response, topic, msg)
}

// 请求,
func (this *sender) Request(topic string, msg *data.WsMsg) error {

	if this.isServer {
		return errors.New("服务器不能使用该方法!")
	}
	return this.sendMsg(data.Request, topic, msg)
}

// 发送消息到网关
func (this *sender) sendMsg(cmd string, topic string, msg *data.WsMsg) error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("sendMsg Exception : %v\n", err)
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
		case data.Request:
			sendMsg = data.BuildResponseMsg(topic, content)
		case data.Response:
			sendMsg = data.BuildRequestMsg(topic, content)
		default:
			{
				//TODO 未定义的cmd
			}
		}

		//sendMsg = data.BuildResponseMsg(topic, content)

		// by, err := json.Marshal(sendMsg)
		// if err != nil {
		// 	return err
		// }
		// this.Conn.WriteJSON(by)

		defer func() {
			if err := recover(); err != nil {
				fmt.Printf(" panic: %v\n", err)
				fmt.Printf(" Stack Info:\n %s \n", debug.Stack())
			}
		}()
		err = this.Conn.WriteJSON(sendMsg)
		if err != nil {
			fmt.Printf("SendMessage  error %v \n", err)
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
