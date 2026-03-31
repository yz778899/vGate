package data

import (
	"encoding/json"
)

// 未解析之前的消息
type NoDecoderMsg struct {
	SessionId int64  `json:"sessionId"` //会话ID
	SnId      int    `json:"snId"`      //消息序列号
	Msg       string `json:"msg"`       //消息内容
}

func (this NoDecoderMsg) MsgSnId() int {
	return this.SnId
}

// 解码消息
func Decoder(ndMsg NoDecoderMsg) (error, *WsMsg) {

	msg := WsMsg{}
	err := json.Unmarshal([]byte(ndMsg.Msg), &msg)
	msg.SessionId = ndMsg.SessionId

	if err == nil {
		switch msg.Cmd {
		case Subscription:
			//订阅消息
			//logic.SubHelper.AddSubscriptionInfo(msg.Topic, ServerManagerInstance.GetSessionById(msg.SessionId).Server)
		//case Publish:
		//发布消息
		case UnSubscription:
			//取消订阅消息
			//logic.SubHelper.UnSubscriptionInfo(msg.Topic, ServerManagerInstance.GetSessionById(msg.SessionId).Server)
		case Notice:
			//通知消息
		case Request:
			//请求消息
		case Response:
			//回复消息
		default:
			//fmt.Printf("未知的消息指令 %v ", msg.Cmd)
			msg.Cmd = Request
			msg.Content = json.RawMessage(ndMsg.Msg)
			//return nil, &msg
		}
	} else {

	}

	return err, &msg
}
