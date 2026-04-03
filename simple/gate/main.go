package main

import (
	"github.com/yz778899/vGate/net"
	"github.com/yz778899/vGate/net/app"
	"github.com/yz778899/vGate/net/handler"
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
