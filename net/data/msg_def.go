package data

import "encoding/json"

const (
	//订阅消息   一般是服务器向网关，订阅来自客户端消息
	Subscription = "subscription"
	//取消订阅
	UnSubscription = "unsubscription"
	//发布消息   网关收到客户端的消息，根据订阅规则来发布 通知到所有匹配的服务器
	Publish = "publish"
	//通知 网关广播给所有服务器
	Notice = "notice"
	//请求 客户端发起请求
	Request = "request"
	//回复 服务器回复客户端的请求
	Response = "response"
)

// 基础消息
type BaseMsg struct {
	Cmd   string `json:"cmd"`   //消息指令 如：订阅Subscription、发布Publish、通知Notice、请求Request、回复Response等
	Topic string `json:"topic"` //订阅主题

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
type PublishMsg struct {
	BaseMsg
	ClientId string          `json:"clientId"` //发布消息的客户端ID
	Content  json.RawMessage `json:"content"`  // 保留原始 JSON
}

// 通知消息 服务器向客户端发起通知
type NoticeMsg struct {
	BaseMsg
	SecretKey string          `json:"secretKey"` //密钥
	Content   json.RawMessage `json:"content"`   // 保留原始 JSON
}

// 请求消息 服务器向客户端发起请求
type RequestMsg = struct {
	BaseMsg
	SessionId int64           `json:"sessionId"` //发布消息的客户端ID
	Content   json.RawMessage `json:"content"`   // 保留原始 JSON
}

// 回复消息 客户端回复服务器的请求
type ResponseMsg = RequestMsg

// 兼容所有种类的消息
type WsMsg struct {
	BaseMsg
	SessionId  int64           `json:"sessionId"`  //发布消息的客户端ID
	ServerName string          `json:"serverName"` //服务器名称
	SecretKey  string          `json:"secretKey"`  //密钥
	Content    json.RawMessage `json:"content"`    // 保留原始 JSON
}
