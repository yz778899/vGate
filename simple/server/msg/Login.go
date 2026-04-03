package msg

import (
	"encoding/json"

	"github.com/yz778899/vGate/net/data"
)

func Decoder(wsMsg *data.WsMsg, reqMsg any) error {
	err := json.Unmarshal(wsMsg.Content, reqMsg)
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
