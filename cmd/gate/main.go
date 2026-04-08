package main

import (
	"fmt"

	//_ "net/http/pprof" // 自动注册 pprof handlers

	"github.com/yz778899/vGate/net"
	"github.com/yz778899/vGate/net/env"
)

func main() {

	// 启动 pprof HTTP 服务（独立端口，不影响业务）
	// go func() {
	// 	log.Println("pprof server started on :6060")
	// 	if err := http.ListenAndServe("localhost:6060", nil); err != nil {
	// 		log.Printf("pprof server error: %v", err)
	// 	}
	// }()
	defer env.Log.Sync()

	err := net.NewWsServer().Run()

	//err := net.NewWsServer().Run()  //启用上面默认参数启动
	if err != nil {
		fmt.Printf("gate failed to 1 start: %v ", err)
	}

}
