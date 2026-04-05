package appmsg

import (
	"encoding/json"

	data "github.com/yz778899/vGate/net/msg"
)

func Decoder(wsMsg *data.WebsocketMsg, reqMsg any) error {
	err := json.Unmarshal(wsMsg.Data, reqMsg)
	return err
}

// 请求 登录
type LoginRequest struct {
	User string
	Pass string
}

// 登录返回
type LoginResponse struct {
	Info string
}
