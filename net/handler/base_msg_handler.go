package handler

import (
	"fmt"

	"github.com/14132465/vGate/net/data"
)

// BaseMsgHandler 基础处理器实现
type BaseMsgHandler struct {
	Topic   string
	Session *data.Session
	Msg     *data.WsMsg
}

// NewBaseMsgHandler 创建基础处理器
func NewBaseMsgHandler(topic string, session *data.Session, msg *data.WsMsg) MsgHandlerInterface {
	return &BaseMsgHandler{
		Session: session,
		Msg:     msg,
		Topic:   topic,
	}
}

func (this *BaseMsgHandler) GetTopic() string {
	return this.Topic
}

// Init 默认实现
func (this *BaseMsgHandler) Init() error {

	return nil
}

// PreProcess 默认实现
func (this *BaseMsgHandler) BeforeProcess() error {
	// 子类可以重写
	return nil
}

// Process 需要子类实现
func (this *BaseMsgHandler) Process() error {
	return fmt.Errorf("Process method not implemented")
}

// PostProcess 默认实现
func (this *BaseMsgHandler) AfterProcess() {
	// 记录处理时间
	// 注意：需要在消息中存储开始时间，这里简化处理

}

// Release 默认实现
func (this *BaseMsgHandler) Release() error {
	// 子类可以重写
	return nil
}

// OnError 错误处理
func (this *BaseMsgHandler) OnError(stage string, err error) {
	// 可以记录日志、发送告警等
	fmt.Printf("[%s] Error in stage %s: %v, MsgId=%s\n", this.GetTopic(), stage, err, this.Msg.Cmd)
}
