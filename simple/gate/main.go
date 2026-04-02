package main

import (
	"github.com/14132465/vGate/net"
	"github.com/14132465/vGate/net/app"
	"github.com/14132465/vGate/net/handler"
	"go.uber.org/zap"
)

func main() {

	defer app.Log.Sync() //

	handler := handler.GateHandler{}
	err := net.NewWsServer().Config(8080, "/").Handler(&handler).Run()
	if err != nil {
		app.Log.Fatal("Server failed to start: ", zap.Error(err))
	}

}
