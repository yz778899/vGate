package main

import (
	"time"

	"github.com/yz778899/vGate/tools"
)

func main() {

	time.Sleep(time.Millisecond * 2)
	for i := 0; i < 100; i++ {
		tools.SimpleTimeExpend()
	}
}

// func testTimeExpend() {
// 	t := tools.NewTimeExpend()
// 	mem := make([]byte, 1024*1024*5)
// 	t.EndAndPrint("申请内存,消耗时长")
// 	t.Reset()
// 	for x := 0; x < len(mem); x++ {
// 		mem[x] = byte(rand.Intn(255))
// 	}
// 	t.EndAndPrint("	写入内存,消耗时长")
// }
