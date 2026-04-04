package handler

import (
	"fmt"

	"github.com/yz778899/vGate/net/data"
)

// BaseMsgHandler 基础处理器实现
type BaseMsgHandler struct {
	Topic   string
	Session *data.Session
	Msg     *data.WebsocketMsg
}

// NewBaseMsgHandler 创建基础处理器
func NewBaseMsgHandler(topic string, session *data.Session, msg *data.WebsocketMsg) MsgHandlerInterface {
	return &BaseMsgHandler{
		Session: session,
		Msg:     msg,
		Topic:   topic,
	}
}

func (this *BaseMsgHandler) GetTopic() string {
	return this.Topic
}

// Init 初始化
func (this *BaseMsgHandler) Init() error {

	return nil
}

// BeforeProcess 业务处理前
func (this *BaseMsgHandler) BeforeProcess() error {
	// 子类可以重写
	return nil
}

// Process 处理业务  需要子类实现
func (this *BaseMsgHandler) Process() error {
	return fmt.Errorf("Process method not implemented")
}

// AfterProcess 处理后
func (this *BaseMsgHandler) AfterProcess() {

}

// Release 释放
func (this *BaseMsgHandler) Release() error {
	// 子类可以重写
	return nil
}

// OnError 错误处理
func (this *BaseMsgHandler) OnError(stage string, err error) {
	// 可以记录日志、发送告警等
	fmt.Printf("[%s] Error in stage %s: %v, MsgId=%s\n", this.GetTopic(), stage, err, this.Msg.Cmd)
}
