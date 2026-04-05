package msg

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

// 服务端  解码消息，将 NoDecoderMsg 转换为 WsMsg
func ServerDecoder(ndMsg NoDecoderMsg) (*WebsocketMsg, error) {
	var msg WebsocketMsg

	// 解析 JSON
	if err := json.Unmarshal([]byte(ndMsg.Msg), &msg); err != nil {
		// 解析失败时，返回未知命令的消息
		msg = WebsocketMsg{
			BaseMsg: BaseMsg{
				Cmd:   Unknown,
				Topic: "",
			},
			SessionId: ndMsg.SessionId,
			Data:      json.RawMessage(ndMsg.Msg),
		}
		return &msg, err
	}
	return &msg, nil
}

// 网关 Decoder 解码消息，将 NoDecoderMsg 转换为 WsMsg
func GateDecoder(ndMsg NoDecoderMsg) (*WebsocketMsg, error) {
	var msg WebsocketMsg

	// 解析 JSON
	if err := json.Unmarshal([]byte(ndMsg.Msg), &msg); err != nil {
		// 解析失败时，返回未知命令的消息
		msg = WebsocketMsg{
			BaseMsg: BaseMsg{
				Cmd:   Unknown,
				Topic: "",
			},
			SessionId: ndMsg.SessionId,
			Data:      json.RawMessage(ndMsg.Msg),
		}
		return &msg, err
	}

	// 根据命令类型设置 SessionId
	switch msg.Cmd {
	case Response:
		//Response类型是app发给client的，其SessionId指向的是client的sessionId,
		// Notice
		// SessionId 保持原样
		return &msg, nil
	case Request, Heartbeat, Subscription, UnSubscription, Notice:
		//Request 需要设置 SessionId
		msg.SessionId = ndMsg.SessionId
		return &msg, nil

	default:
		// 未知命令：设置 SessionId 并标记cmd为 Unknown
		msg.SessionId = ndMsg.SessionId
		//msg.Cmd = Unknown
		msg.Cmd = Request //简化，默认为客户端请求，可以省掉一个字段
		msg.Data = json.RawMessage(ndMsg.Msg)
		return &msg, nil
	}
}
