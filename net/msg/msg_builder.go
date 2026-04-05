package msg

import "encoding/json"

//构建 订阅 的消息
func BuildSubscriptionMsg(topic string, serverName string, secretKey string) *SubscriptionMsg {
	return &SubscriptionMsg{
		BaseMsg: BaseMsg{Cmd: Subscription, Topic: topic},

		ServerName: serverName,
		SecretKey:  secretKey,
	}
}

//构建 取消订阅 的消息
func BuildUnSubscriptionMsg(topic string, serverName string) *UnSubscriptionMsg {
	return &UnSubscriptionMsg{
		BaseMsg:    BaseMsg{Cmd: UnSubscription, Topic: topic},
		ServerName: serverName,
	}
}

//构建 发布 的消息
// func BuildPublishMsg(clientId string, topic string, data string) *PublishMsg {
// 	return &PublishMsg{
// 		BaseMsg:  BaseMsg{Cmd: Subscription, Topic: topic},
// 		ClientId: clientId,
// 		Data:  json.RawMessage(data),
// 	}
// }

//构建 通知 的消息
func BuildNoticeMsg(secretKey string, topic string, data []byte) *NoticeMsg {
	return &NoticeMsg{
		BaseMsg: BaseMsg{Cmd: Notice, Topic: topic, Data: json.RawMessage(data)},

		SecretKey: secretKey,
	}
}

//构建 客户端请求 的消息
func BuildRequestMsg(sessionId int64, topic string, data []byte) *RequestMsg {
	return &RequestMsg{
		BaseMsg:   BaseMsg{Cmd: Request, Topic: topic},
		Data:      json.RawMessage(data),
		SessionId: sessionId,
	}
}

//构建 客户端请求 的消息
func BuildResponseMsg(sessionId int64, topic string, data []byte) *ResponseMsg {
	return &ResponseMsg{
		BaseMsg:   BaseMsg{Cmd: Response, Topic: topic},
		Data:      json.RawMessage(data),
		SessionId: sessionId,
	}
}
