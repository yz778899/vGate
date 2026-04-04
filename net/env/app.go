package env

import (
	"github.com/yz778899/vGate/net/data"
	"github.com/yz778899/vGate/net/env/config"
)

// 入口
type gate struct {
	//全局会话管理器实例
	SessionMgr *data.SessionManager
	//服务器管理器实例
	AppSessionMgr *data.AppServiceManager
	//全局订阅信息管理器实例
	//SubHelper *logic.SubscriptionHelper
	//配置
	Config *config.RootConfig
}

var (
	VGate *gate
)

func init() {
	VGate = &gate{
		Config:        config.GetConfig("config.yaml"),
		SessionMgr:    data.SessionManagerInstance,
		AppSessionMgr: data.ServerManagerInstance,
		//SubHelper:     logic.SubHelper,
	}

	if err := InitLogger(VGate.Config); err != nil {
		panic(err)
	}
}

// CheckSecretKey检查提供的密钥是否与全局密钥匹配
func (this *gate) CheckSecretKey(key string) bool {
	return this.Config.Gate.SecretKey == "" || this.Config.Gate.SecretKey == key
}
