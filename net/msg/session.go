package msg

import (
	"runtime/debug"
	"sync"
	"sync/atomic"

	"github.com/gofiber/fiber/v2/log"
)

// Session会话结构体，包含客户端ID、会话状态、HTTP请求和响应对象等信息
type Session struct {
	*ConnSession
	Status int8 //会话状态 0：未连接 1：已连接 2：已断开
}

// 发送消息
func (this *Session) SendToClient(msg *ToClientMsg) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(" panic: %v\n", err)
			log.Error(" Stack Info:\n %s \n", debug.Stack())
		}
	}()
	err := this.Conn.WriteJSON(msg)
	if err != nil {
		log.Error("SendToClient  error %v \n", err)
	}
}

// 客户端发给网关
func (this *Session) SendToGate(msg *ToClientMsg) {

}

// 发送消息给服务
func (this *Session) SendToService(msg *WebsocketMsg) {
	defer func() {
		if err := recover(); err != nil {
			log.Error(" panic: %v\n", err)
			log.Error(" Stack Info:\n %s \n", debug.Stack())
		}
	}()
	err := this.Conn.WriteJSON(msg)
	if err != nil {
		log.Error("session SendMessage  error %v \n", err)
	}
}

type SessionManager struct {
	sessionMap map[int64]*Session //会话映射表，存储所有客户端的会话信息
	uuid       atomic.Int64
	mutex      sync.RWMutex
}

// 全局会话管理器实例
var SessionManagerInstance *SessionManager

// 初始化会话管理器实例
func init() {
	SessionManagerInstance = &SessionManager{
		sessionMap: make(map[int64]*Session, 512),
		uuid:       atomic.Int64{},
	}
	SessionManagerInstance.uuid.Store(10000 * 10000 * 20) //初始值设置为一个较大的数[20亿]，避免与实际客户端ID冲突
}

// 根据客户端ID获取会话信息
func (sm *SessionManager) GetSession(uuid int64) *Session {
	defer sm.mutex.RLocker().Unlock()
	sm.mutex.RLocker().Lock()
	if session, ok := sm.sessionMap[uuid]; ok {
		return session
	}
	return nil
}

// 添加会话信息
func (sm *SessionManager) AddSession(session *Session) *Session {

	if session.UUID <= 0 {
		id := sm.uuid.Add(1)
		session.UUID = id
		defer sm.mutex.Unlock()
		sm.mutex.Lock()
		sm.sessionMap[session.UUID] = session
	} else {
		//客户端ID已存在
	}
	return session

}

// 移除会话信息
func (sm *SessionManager) RemoveSession(uuid int64) {
	defer sm.mutex.Unlock()
	sm.mutex.Lock()
	delete(sm.sessionMap, uuid)
}

// 更新会话状态
func (sm *SessionManager) UpdateSessionStatus(uuid int64, status int8) {
	defer sm.mutex.Unlock()
	sm.mutex.Lock()
	if session, ok := sm.sessionMap[uuid]; ok {
		session.Status = status
	}
}

// 变更sessionID
func (sm *SessionManager) ChangeId(uuid int64, newId int64) {
	defer sm.mutex.Unlock()
	sm.mutex.Lock()
	if session, ok := sm.sessionMap[uuid]; ok {
		session.UUID = newId
		delete(sm.sessionMap, uuid)
		sm.sessionMap[newId] = session
	}
}
