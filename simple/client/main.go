package main

import (
	"time"

	"github.com/yz778899/vGate/net/app"
	_ "github.com/yz778899/vGate/net/app"
)

func main() {

	app.Log.Info(" first log  ")

	for {
		time.Sleep(time.Microsecond)
	}

}
