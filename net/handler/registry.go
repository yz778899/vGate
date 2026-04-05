package handler

import (
	"fmt"
	"sync"

	"github.com/gofiber/fiber/v2/log"
	"github.com/yz778899/vGate/net/msg"
)

// Registry 处理器注册中心
type Registry struct {
	MsgHandlerCreates map[string]MsgHandlerCreate // topic -> MsgHandlerCreate
	mu                sync.RWMutex
}

var (
	RegistryInstance *Registry
)

// NewRegistry 创建注册中心
func NewRegistry() *Registry {
	if RegistryInstance == nil {
		RegistryInstance = &Registry{
			MsgHandlerCreates: make(map[string]MsgHandlerCreate),
		}
	}
	return RegistryInstance
}

// Register 注册处理器
func (r *Registry) Register(handlerCreate MsgHandlerCreate) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	topic := handlerCreate.Topic
	if _, exists := r.MsgHandlerCreates[topic]; exists {
		return fmt.Errorf("MsgHandlerCreate for topic %s already registered", topic)
	}

	r.MsgHandlerCreates[topic] = handlerCreate
	return nil
}

// GetMsgHandlerCreate 根据主题获取处理器
func (r *Registry) GetMsgHandlerCreate(topic string) (MsgHandlerCreate, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	MsgHandlerCreate, ok := r.MsgHandlerCreates[topic]
	return MsgHandlerCreate, ok
}

// Unregister 注销处理器
func (r *Registry) Unregister(topic string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.MsgHandlerCreates[topic]; !exists {
		return fmt.Errorf("MsgHandlerCreate for topic %s not found", topic)
	}

	delete(r.MsgHandlerCreates, topic)
	return nil
}

// ListTopics 列出所有已注册的主题
func (r *Registry) ListTopics() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	topics := make([]string, 0, len(r.MsgHandlerCreates))
	for topic := range r.MsgHandlerCreates {
		topics = append(topics, topic)
	}
	return topics
}

// 创建一个处理器，并控制它的生周期运行
func (r *Registry) RunHandler(msg *msg.WebsocketMsg, session *msg.Session) error {

	defer func() {
		if err := recover(); err != nil {
			log.Error(fmt.Printf("处理消息时发生错误: %v\n", err))
		}
	}()

	creator, ok := RegistryInstance.GetMsgHandlerCreate(msg.Topic)
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
