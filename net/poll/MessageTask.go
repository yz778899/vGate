package pool

import "github.com/yz778899/vGate/net/msg"

//消息任务
type MessageTask struct {
	SessionId int64
	Msg       *msg.WebsocketMsg
}

func (this *MessageTask) GetZoneId() int {
	return int(this.SessionId)
}

type Task interface {
	//使用该值以判断使用哪个 slave来处理
	GetZoneId() int
}
