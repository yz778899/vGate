package pool

import "github.com/yz778899/vGate/net/msg"

//消息任务
type MessageTask struct {
	SessionId int64
	Msg       *msg.WebsocketMsg
}

func (this *MessageTask) SlaveId() int {
	return int(this.SessionId)
}
