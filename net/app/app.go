package app

import (
	"github.com/yz778899/vGate/net/app/config"
	"github.com/yz778899/vGate/net/data"
	"github.com/yz778899/vGate/net/logic"
)

// 入口
type gate struct {
	//全局会话管理器实例
	SessionManager *data.SessionManager
	//全局服务器管理器实例
	ServerManager *data.ServerManager
	//全局订阅信息管理器实例
	SubHelper *logic.SubscriptionHelper
	//密钥 用于网关与服务器通讯，判断是否一致
	//SecretKey string
	Config *config.RootConfig
	//ZapLog *zap.Logger
}

//const Gate  = &gate{}

var (
	VGate *gate
)

func init() {
	VGate = &gate{
		Config:         config.GetConfig("config.yaml"),
		SessionManager: data.SessionManagerInstance,
		ServerManager:  data.ServerManagerInstance,
		SubHelper:      logic.SubHelper,
		//ZapLog:         Log,
	}

	//defer Log.Sync()

	if err := InitLogger(VGate.Config); err != nil {
		panic(err)
	}
	Log.Info("服务启动成功")

}

// CheckSecretKey检查提供的密钥是否与全局密钥匹配
func (this *gate) CheckSecretKey(key string) bool {
	return this.Config.Gate.SecretKey == "" || this.Config.Gate.SecretKey == key
}
