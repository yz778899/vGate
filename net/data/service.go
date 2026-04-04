package data

import (
	"fmt"
	"runtime/debug"
	"sync"
	"sync/atomic"
)

// AppServer会话结构体，包含客户端ID、会话状态、HTTP请求和响应对象等信息
type AppServer struct {
	Session Session
	//Name    string
}

// 发送消息
func (this *AppServer) SendMessage(msg any) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf(" panic: %v\n", err)
			fmt.Printf(" Stack Info:\n %s \n", debug.Stack())
		}
	}()
	err := this.Session.Conn.WriteJSON(msg)
	if err != nil {
		fmt.Printf("SendMessage  error %v \n", err)
	}
}

type AppServiceManager struct {
	ServerMap map[int64]*AppServer //会话映射表，存储所有客户端的会话信息
	uuid      atomic.Int64
	mutex     sync.RWMutex
}

// 全局会话管理器实例
var ServerManagerInstance *AppServiceManager

// 初始化会话管理器实例
func init() {
	ServerManagerInstance = &AppServiceManager{
		ServerMap: make(map[int64]*AppServer),
		uuid:      atomic.Int64{},
	}
	ServerManagerInstance.uuid.Store(10000*10000*20 - 1000) //初始值设置为一个较大的数[20亿]，避免与实际客户端ID冲突
}

// 根据客户端ID获取会话信息 , 只读不创建
func (sm *AppServiceManager) GetServerOnly(uuid int64) *AppServer {
	defer sm.mutex.RLocker().Unlock()
	sm.mutex.RLocker().Lock()
	if Server, ok := sm.ServerMap[uuid]; ok {
		return Server
	}
	return nil
}

// 根据客户端ID获取会话信息
func (sm *AppServiceManager) GetAndCreateServer(uuid int64) *AppServer {
	defer sm.mutex.Unlock()
	sm.mutex.Lock()
	if Server, ok := sm.ServerMap[uuid]; ok {
		return Server
	}
	//上面只读，这里有可能要写了

	//如果在会话管理器中找不到对应的会话信息，则创建一个新的Server对象并返回
	session := SessionManagerInstance.GetSession(uuid)
	if session != nil {
		sm.ServerMap[uuid] = &AppServer{
			Session: *session,
		}

		return sm.ServerMap[uuid]
	}
	return nil
}

// 添加会话信息
func (sm *AppServiceManager) AddServer(service *AppServer) *AppServer {

	if service.Session.UUID <= 0 {
		id := sm.uuid.Add(1)
		service.Session.UUID = id
		defer sm.mutex.Unlock()
		sm.mutex.Lock()
		sm.ServerMap[service.Session.UUID] = service
	} else {
		//客户端ID已存在
	}
	return service

}

// 移除会话信息
func (sm *AppServiceManager) RemoveServer(uuid int64) {
	defer sm.mutex.Unlock()
	sm.mutex.Lock()
	delete(sm.ServerMap, uuid)
}

// 更新会话状态
func (sm *AppServiceManager) UpdateServerStatus(uuid int64, status int8) {
	defer sm.mutex.Unlock()
	sm.mutex.Lock()
	if Server, ok := sm.ServerMap[uuid]; ok {
		Server.Session.Status = status
	}
}

// 取得所有服务器
func (sm *AppServiceManager) GetAlls() []*AppServer {
	defer sm.mutex.RLocker().Unlock()
	sm.mutex.RLocker().Lock()
	lst := make([]*AppServer, 0)
	for _, v := range sm.ServerMap {
		if v != nil {
			lst = append(lst, v)
		}
	}
	return lst

}
