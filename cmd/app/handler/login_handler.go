package msg_handler

import (
	"fmt"
	"math/rand"

	appmsg "github.com/yz778899/vGate/cmd/app/app_msg"
	"github.com/yz778899/vGate/net/handler"
	"github.com/yz778899/vGate/net/logic"
	"github.com/yz778899/vGate/net/msg"
)

type LoginHandler struct {
	handler.BaseMsgHandler
	Request *appmsg.LoginRequest
}

func NewLoginHandler(topic string, session *msg.Session, msg *msg.WebsocketMsg) handler.MsgHandlerInterface {
	hdl := LoginHandler{}
	hdl.Topic = topic
	hdl.Msg = msg
	return &hdl
}

// Init() 使用父类 handler.BaseMsgHandler 默认实现

// PreProcess处理前
func (this *LoginHandler) BeforeProcess() error {
	//解码得到请求消息体
	this.Request = &appmsg.LoginRequest{}
	err := appmsg.Decoder(this.Msg, this.Request)
	return err
}

// Process 需要子类实现
func (this *LoginHandler) Process() error {

	sid := this.Msg.SessionId
	newId := sid
	info := "登录成功!"
	if sid > 10000*10000 {
		newId = int64(1000 + rand.Intn(9000))
		//log.Info("摸拟登录中……")
		info = fmt.Sprintf("登录成功! 你的用户名 %v  密码 %v, 原sessionId = %v  新id = %v  !  向网关 发送请求变更  ！", this.Request.User, this.Request.Pass, sid, newId)
		//用户session ID变更 消息结构体
		changeMsg := logic.SessionIdChange{SessionId: sid,
			NewId: newId}
		logic.Sender.Notice(logic.Session_Id_Change, changeMsg)

	}

	resp := &appmsg.LoginResponse{Info: info}
	logic.Sender.Resp(newId, this.Msg.GetTopic(), resp)

	//log.Info(info)
	return nil
}

// AfterProcess 处理后
func (this *LoginHandler) AfterProcess() {

}
