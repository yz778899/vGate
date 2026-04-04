package msg_handler

import (
	"fmt"
	"math/rand"

	"github.com/gofiber/fiber/v2/log"
	"github.com/yz778899/vGate/cmd/app/msg"
	"github.com/yz778899/vGate/net/data"
	"github.com/yz778899/vGate/net/handler"
	"github.com/yz778899/vGate/net/logic"
)

type LoginHandler struct {
	handler.BaseMsgHandler
	Request *msg.LoginRequest
}

func NewLoginHandler(topic string, session *data.Session, msg *data.WebsocketMsg) handler.MsgHandlerInterface {
	hdl := LoginHandler{}
	hdl.Topic = topic
	hdl.Msg = msg
	return &hdl
}

// Init() 使用父类 handler.BaseMsgHandler 默认实现

// PreProcess处理前
func (this *LoginHandler) BeforeProcess() error {
	//解码得到请求消息体
	this.Request = &msg.LoginRequest{}
	err := msg.Decoder(this.Msg, this.Request)
	return err
}

// Process 需要子类实现
func (this *LoginHandler) Process() error {

	sid := this.Msg.SessionId
	newId := int64(1000 + rand.Intn(9000))

	log.Info("摸拟登录中……")
	info := fmt.Sprintf("登录成功! 你的用户名 %v  密码 %v, 原sessionId = %v  新id = %v  !  向网关 发送请求变更  ！", this.Request.User, this.Request.Pass, sid, newId)
	log.Info(info)

	//用户session ID变更 消息结构体
	changeMsg := logic.SessionIdChange{SessionId: sid,
		NewId: newId}
	logic.Sender.Notice(logic.Session_Id_Change, changeMsg)

	// this.Session.SendMessage(logic.Session_Id_Change, changeMsg)

	resp := &msg.LoginResponse{Info: info}

	logic.Sender.Resp(newId, this.Msg.GetTopic(), resp)

	return nil
}

// AfterProcess 处理后
func (this *LoginHandler) AfterProcess() {

}
