package msg

import (
	"encoding/json"
	"sync"
	"sync/atomic"

	"github.com/gofiber/fiber/v2/log"
	"github.com/gorilla/websocket"
)

// 客户端专属
type ConnSession struct {
	//客户端ID
	UUID int64
	//connect
	Conn *websocket.Conn
	//针对conn读写加锁
	Mutex sync.Mutex

	Closed atomic.Bool
}

// // 发送消息[到网关]
// func (this *ConnSession) SendToClient(msg *ToClientMsg) {
// 	defer func() {
// 		if err := recover(); err != nil {
// 			log.Error(" panic: %v\n", err)
// 			log.Error(" Stack Info:\n %s \n", debug.Stack())
// 		}
// 	}()
// 	err := this.Conn.WriteJSON(msg)
// 	if err != nil {
// 		log.Error("SendToClient  error %v \n", err)
// 	}
// }

// 发送消息[到网关]
func (this *ConnSession) SendToGate(Topic string, msg any) error {

	msgByteArray, err := json.Marshal(msg)
	if err != nil {
		log.Error("Json解析出错, ", err)
		return err
	}
	pack := ToClientMsg{
		Topic: Topic,
		Data:  msgByteArray,
	}
	this.Mutex.Lock()
	err = this.Conn.WriteJSON(pack)
	this.Mutex.Unlock()
	if err != nil {
		log.Error("send_session SendMessage  error %v \n", err)
	}
	return err
}

// 发送消息给服务
// func (this *ConnSession) SendToService(msg *WebsocketMsg) {
// 	defer func() {
// 		if err := recover(); err != nil {
// 			log.Error(" panic: %v\n", err)
// 			log.Error(" Stack Info:\n %s \n", debug.Stack())
// 		}
// 	}()
// 	err := this.Conn.WriteJSON(msg)
// 	if err != nil {
// 		log.Error("SendMessage  error %v \n", err)
// 	}
// }
