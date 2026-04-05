package main

import (
	"fmt"

	"github.com/yz778899/vGate/net"
	"github.com/yz778899/vGate/net/env/config"
)

func main() {

	//defer env.Log.Sync()

	err := net.NewWsServer().WithConfig(&config.GateConfig{
		WsPath:        "/",  //websocket路径
		WsPort:        6789, //网关启动端口
		SecretKey:     "",   //密钥 如设置 app 与 gate 双方需要一致，为空则不较验
		HeartbeatTime: 3,    //心跳频率 -- 仅 app服务需要 gate不需要
		ReadOverTime:  7,    //读写超时秒数
	}).Run()

	//err := net.NewWsServer().Run()  //启用上面默认参数启动
	if err != nil {
		fmt.Printf("gate failed to start: %v ", err)
	}

}
