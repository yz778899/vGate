package main

import (
	"github.com/14132465/vGate/net"
	"github.com/14132465/vGate/net/app"
	"github.com/14132465/vGate/net/handler"
)

func main() {

	app.VGate.SetSecretKey("ga-23xk=v") // 设置全局密钥 ，如果没设则不检查
	handler := handler.GateHandler{}
	net.NewWsServer().Config(8080, "/").Handler(&handler).Run()

}
