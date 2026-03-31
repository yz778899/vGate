package app

import (
	"github.com/14132465/vGate/net/data"
	"github.com/14132465/vGate/net/logic"
)

// 入口
// type gate struct {
// 	//SessionManager *data.SessionManager
// }

var (
	//VGate          = &gate{}

	//全局会话管理器实例
	SessionManager *data.SessionManager

	//全局服务器管理器实例
	ServerManager *data.ServerManager
	//全局订阅信息管理器实例
	SubHelper *logic.SubscriptionHelper
	//密钥 用于网关与服务器通讯，判断是否一致
	secretKey string
)

func init() {
	SessionManager = data.SessionManagerInstance
	SubHelper = logic.SubHelper
}

// SetSecretKey设置全局密钥
func SetSecretKey(key string) {
	secretKey = key
}

// CheckSecretKey检查提供的密钥是否与全局密钥匹配
func CheckSecretKey(key string) bool {
	return secretKey == "" || secretKey == key
}
