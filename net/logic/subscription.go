package logic

import (
	"sync"

	"github.com/yz778899/vGate/net/data"
)

// SubscriptionInfo结构体表示一个订阅信息，包含订阅的主题和对应的服务器会话信息
type SubscriptionInfo struct {
	Topic  string
	Server *data.AppServer
}

// SubscriptionInfoList []SubscriptionInfo
type SubscriptionHelper struct {
	sync.RWMutex
	// key topic value 订阅服务器列表
	SubscriptionMap map[string][]SubscriptionInfo
}

// 全局订阅信息管理器实例
var SubHelper = &SubscriptionHelper{}

func init() {
	SubHelper.SubscriptionMap = make(map[string][]SubscriptionInfo)
}

// 添加订阅信息
func (this *SubscriptionHelper) AddSubscriptionInfo(topic string, server *data.AppServer) {
	this.Lock()
	defer this.Unlock()
	sub := SubscriptionInfo{
		Topic:  topic,
		Server: server,
	}
	this.SubscriptionMap[topic] = append(this.SubscriptionMap[topic], sub)
}

// 获取订阅信息列表
func (this *SubscriptionHelper) GetSubscriptionInfo(topic string) []SubscriptionInfo {
	this.RLock()
	defer this.RUnlock()
	return this.SubscriptionMap[topic]
}

// 移除订阅信息
func (this *SubscriptionHelper) UnSubscriptionInfo(topic string, server *data.AppServer) {
	this.Lock()
	defer this.Unlock()
	list := this.SubscriptionMap[topic]

	for i, sub := range list {
		if sub.Server == server {
			// 从切片中移除该订阅信息
			this.SubscriptionMap[topic] = append(list[:i], list[i+1:]...)
			break
		}
	}
}

// 服务器断开，需要清除所有相关的订阅数据
func (this *SubscriptionHelper) ServerClose(server *data.AppServer) {
	this.Lock()
	defer this.Unlock()

	for topic, list := range this.SubscriptionMap {
		for i, sub := range list {
			if sub.Server == server {
				// 从切片中移除该订阅信息
				this.SubscriptionMap[topic] = append(list[:i], list[i+1:]...)
				break
			}
		}

	}

}

// 广播消息给订阅了指定主题的所有服务器
func (this *SubscriptionHelper) Broadcast(topic string, msg *data.WebsocketMsg) {
	subs := this.GetSubscriptionInfo(topic)
	for _, sub := range subs {
		sub.Server.SendMessage(msg)
	}
}
