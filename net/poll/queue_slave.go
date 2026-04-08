package pool

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"sync"
	"time"
)

// 队列处理线程
type QueueSlave struct {
	//消息队列
	queue []*MessageTask
	//消息处理器
	MsgHandler Handler
	Name       string
	Id         int
	mutex      sync.Mutex
}

func (this *QueueSlave) Init() {
	this.queue = make([]*MessageTask, 0)
}

// 启动协程，持续处理消息 会阻塞 应当以  go Start() 启动
func (this *QueueSlave) Start() {

	count := 0
	for {

		this.mutex.Lock()
		// 出队 (Dequeue)
		len := len(this.queue)
		if len == 0 {
			this.mutex.Unlock()
			count = 0
			time.Sleep(time.Millisecond * 1)
		} else {
			task := this.queue[0]
			this.queue = this.queue[1:]
			this.mutex.Unlock()
			this.Handler(task)

			//fmt.Printf(" handler task  = %v ", task)

			count++
		}
		//处理10条数据就休眠
		if count%5 == 0 {
			//time.Sleep(time.Millisecond * 1)
			runtime.Gosched()
		}

	}
}

// 处理单条消息
func (this *QueueSlave) Handler(task *MessageTask) {

	defer func() {
		if p := recover(); p != nil {
			fmt.Printf(" 捕获到 panic: %v\n", p)
			fmt.Printf("堆栈信息:\n%s\n", debug.Stack())
		}
	}()

	//handler.RegistryInstance.RunHandler(task.Msg, task.SessionId)

	if this.MsgHandler.Fun == nil {
		//没有消息处理器，则默认输出信息
		fmt.Printf("  ###  handler info =  【id =%v, name = %v 】 \n	 msg =  %#v , 剩余数量 = %v  \n", this.Id, this.Name, task, len(this.queue))
	} else {
		this.MsgHandler.Fun(task)
	}
}

// 接受消息
func (this *QueueSlave) Accept(task *MessageTask) {
	defer this.mutex.Unlock()
	this.mutex.Lock()
	this.queue = append(this.queue, task)

	//fmt.Printf("  		接收到新task  %v   \n", task)

}
