package msg_handler

import (
	"github.com/gofiber/fiber/v2/log"
	appmsg "github.com/yz778899/vGate/cmd/app/app_msg"
	"github.com/yz778899/vGate/net/handler"
	"github.com/yz778899/vGate/net/logic"
	"github.com/yz778899/vGate/net/msg"
)

// 游戏列表处理器
type GameListHandler struct {
	handler.BaseMsgHandler
	Request *appmsg.GameListRequest
}

func NewGameListHandler(topic string, session *msg.Session, msg *msg.WebsocketMsg) handler.MsgHandlerInterface {
	hdl := GameListHandler{}
	hdl.Topic = topic
	hdl.Msg = msg
	return &hdl
}

// 处理前
func (this *GameListHandler) BeforeProcess() error {
	//解码得到请求消息体
	this.Request = &appmsg.GameListRequest{}
	err := appmsg.Decoder(this.Msg, this.Request)
	return err
}

// 处理
func (this *GameListHandler) Process() error {
	resp := &appmsg.GameListResponse{Games: []appmsg.Game{}}
	logic.Sender.Resp(this.Msg.SessionId, this.Msg.GetTopic(), resp)
	return nil
}

// 处理后
func (this *GameListHandler) AfterProcess() {
}

// 释放
func (this *GameListHandler) Release() error {
	return nil
}

// 错误处理
func (this *GameListHandler) OnError(stage string, err error) {
	log.Error(err)
}
