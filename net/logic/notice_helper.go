package logic

import (
	"encoding/json"
	"fmt"

	"github.com/14132465/vGate/net/data"
	"github.com/gofiber/fiber/v2/log"
)

const (
	//客户端上线通知
	Notice_On_Line string = "online"
	//客户端下线通知
	Notice_Off_Line string = "offline"
	//session_id 变更 [一般用户登录后可以通知网关主动变更，但需要保证其唯一性，]
	Session_Id_Change string = "session_id_change"
)

// 用户session ID变更
type SessionIdChange struct {
	SessionId int64
	NewId     int64
}

// 网关处理通知的逻辑
type NoticeHelper struct {
}

var NoticeHelperInstance *NoticeHelper

func init() {
	if NoticeHelperInstance == nil {
		NoticeHelperInstance = &NoticeHelper{}
	}
}

func (this *NoticeHelper) Handler(msg *data.WsMsg) bool {

	if msg.Topic == Session_Id_Change {
		change := SessionIdChange{}
		err := json.Unmarshal(msg.Content, &change)
		if err != nil {
			log.Error(fmt.Printf("SessionIdChange 反序列化出错 %v ", msg))
		} else {
			data.SessionManagerInstance.ChangeId(change.SessionId, change.NewId)
		}
		return true
	} else {
		return false
	}

}
