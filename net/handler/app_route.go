package handler

import (
	"fmt"
	"sync"

	"github.com/gofiber/fiber/v2/log"
	"github.com/yz778899/vGate/net/msg"
)

// AppRoute 消息路由
type AppRoute struct {
	Creaters map[string]Creater // topic -> MsgHandlerCreate
	mu       sync.RWMutex
}

var (
	Default *AppRoute
)

// NewAppRoute 创建注册中心
func NewAppRoute() *AppRoute {
	if Default == nil {
		Default = &AppRoute{
			Creaters: make(map[string]Creater),
		}
	}
	return Default
}

// Add 注册处理器
func (r *AppRoute) Add(handlerCreate Creater) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	topic := handlerCreate.Topic
	if _, exists := r.Creaters[topic]; exists {
		return fmt.Errorf("MsgHandlerCreate for topic %s already registered", topic)
	}

	r.Creaters[topic] = handlerCreate
	return nil
}

// GetMsgHandlerCreate 根据主题获取处理器
func (r *AppRoute) GetMsgHandlerCreate(topic string) (Creater, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	MsgHandlerCreate, ok := r.Creaters[topic]
	return MsgHandlerCreate, ok
}

// Remove 注销处理器
func (r *AppRoute) Remove(topic string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.Creaters[topic]; !exists {
		return fmt.Errorf("MsgHandlerCreate for topic %s not found", topic)
	}

	delete(r.Creaters, topic)
	return nil
}

// ListTopics 列出所有已注册的主题
func (r *AppRoute) ListTopics() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	topics := make([]string, 0, len(r.Creaters))
	for topic := range r.Creaters {
		topics = append(topics, topic)
	}
	return topics
}

// 创建一个处理器，并控制它的生周期运行
func (r *AppRoute) Exec(msg *msg.WebsocketMsg, session *msg.Session) error {

	defer func() {
		if err := recover(); err != nil {
			log.Error(fmt.Printf("处理消息时发生错误: %v\n", err))
		}
	}()

	creator, ok := Default.GetMsgHandlerCreate(msg.Topic)
	if ok {

		hdl := creator.CreateFunc(msg.Topic, session, msg)

		//初始化
		err := hdl.Init()
		if err != nil {
			return err
		}
		//处理前
		err = hdl.BeforeProcess()
		if err != nil {
			return err
		}
		//处理中
		err = hdl.Process()
		if err != nil {
			return err
		}
		//处理后
		hdl.AfterProcess()

	} else {
		log.Error("### 缺少 MsgHandler 对应 topic 是 " + msg.Topic + " ，该消息将丢弃处理 !\n")
	}
	return nil
}
