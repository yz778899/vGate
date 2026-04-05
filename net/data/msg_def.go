package data

import (
	"encoding/json"
	"reflect"
	"strings"
)

const (
	//订阅消息   一般是服务器向网关，订阅来自客户端消息
	Subscription = "subscription"
	//取消订阅
	UnSubscription = "unsubscription"
	//发布消息   网关收到客户端的消息，根据订阅规则来发布 通知到所有匹配的服务器
	//Publish = "publish"

	//通知 网关广播给所有服务器 或者 先由服务器通知网关，网关再转发到所有服务器
	Notice = "notice"
	//请求 客户端发起请求
	Request = "request"
	//回复 服务器回复客户端的请求
	Response = "response"
	//心跳 收到心跳后 网关会立即回复
	Heartbeat = "heartbeat"
	//未知指令
	Unknown = "unknown"
)

type BaseMsgInterFace interface {
	GetCmd() string
	GetTopic() string
	GetData() json.RawMessage
}

// 基础消息
type BaseMsg struct {
	Cmd   string          `json:"cmd"`   //消息指令 如：订阅Subscription、发布Publish、通知Notice、请求Request、回复Response等
	Topic string          `json:"topic"` //订阅主题
	Data  json.RawMessage `json:"data"`  // 保留原始 JSON

}

// 最终发给客户端的消息
type ToClientMsg struct {
	Topic string          `json:"topic"` //订阅主题
	Data  json.RawMessage `json:"data"`  // 保留原始 JSON
}

func (this *ToClientMsg) TransitionOf(msg *WebsocketMsg) *ToClientMsg {
	this.Topic = msg.Topic
	this.Data = msg.Data
	return this
}

func (this *BaseMsg) GetCmd() string {
	return this.Cmd
}

func (this *BaseMsg) GetTopic() string {
	return this.Topic
}
func (this *BaseMsg) GetData() json.RawMessage {
	return this.Data
}

// 心跳
type heartbeatMsg struct {
	Cmd string `json:"cmd"` //消息指令 如：订阅Subscription、发布Publish、通知Notice、请求Request、回复Response等
}

var heartbeatMsgInstance *heartbeatMsg

// 心跳消息
func HeartbeatMsg() *heartbeatMsg {
	if heartbeatMsgInstance == nil {
		heartbeatMsgInstance = &heartbeatMsg{Cmd: Heartbeat}
	}
	return heartbeatMsgInstance
}

// 订阅消息 服务器向网关发起订阅请求
type SubscriptionMsg = struct {
	BaseMsg
	ServerName string `json:"serverName"` //服务器名称
	SecretKey  string `json:"secretKey"`  //密钥 核对密钥是否与网关一致，否则无法订阅

}

// 取消订阅消息 服务器向网关发起取消订阅请求
type UnSubscriptionMsg = SubscriptionMsg

// 发布消息 客户端向服务器发布消息
// type PublishMsg struct {
// 	BaseMsg
// 	ClientId string          `json:"clientId"` //发布消息的客户端ID
// 	Data  json.RawMessage `json:"data"`  // 保留原始 JSON
// }

// 通知消息 向所有服务器 发起通知 , 服务器无需订阅
type NoticeMsg struct {
	BaseMsg
	SecretKey string `json:"secretKey"` //密钥
}

// 请求消息 服务器向客户端发起请求
type RequestMsg = struct {
	BaseMsg
	SessionId int64           `json:"sessionId"` //发布消息的客户端ID
	Data      json.RawMessage `json:"data"`      // 保留原始 JSON
}

// 回复消息 客户端回复服务器的请求
type ResponseMsg = RequestMsg

// 兼容所有种类的消息
type WebsocketMsg struct {
	BaseMsg
	ServerName string          `json:"serverName"` //服务器名称
	SessionId  int64           `json:"sessionId"`  //发布消息的客户端ID
	SecretKey  string          `json:"secretKey"`  //密钥
	Data       json.RawMessage `json:"data"`       // 保留原始 JSON
}

// CustomMessage 自定义序列化
type CustomMessage struct {
	WebsocketMsg
	HideFields []string `json:"-"` // 要隐藏的字段
}

// 自定义 MarshalJSON 方法
func (c CustomMessage) MarshalJSON() ([]byte, error) {
	// 创建一个map来存储要序列化的字段
	result := make(map[string]interface{})

	// 使用反射获取Message结构体的值
	msgValue := reflect.ValueOf(c.WebsocketMsg)
	msgType := reflect.TypeOf(c.WebsocketMsg)

	// 遍历Message的所有字段
	for i := 0; i < msgValue.NumField(); i++ {
		field := msgType.Field(i)
		fieldValue := msgValue.Field(i)

		// 获取json标签
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			jsonTag = strings.ToLower(field.Name)
		}

		// 检查是否需要隐藏此字段
		shouldHide := false
		for _, hideField := range c.HideFields {
			if strings.EqualFold(hideField, jsonTag) || strings.EqualFold(hideField, field.Name) {
				shouldHide = true
				break
			}
		}

		// 如果不需要隐藏，则添加到结果map中
		if !shouldHide {
			result[jsonTag] = fieldValue.Interface()
		}
	}

	return json.Marshal(result)
}
