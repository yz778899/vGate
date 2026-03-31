package coroutine

import (
	"container/list"
	"fmt"
	"runtime/debug"
	"sync"
	"time"
)

// 一个协程
type Coroutine struct {
	Name string
	Id   int
	//消息处理器
	MsgHandler Handler
	//消息队列
	queue list.List
	mutex sync.Mutex
}

// 消息结构
type V1Msg interface {
	//消息编号 取模以分配给子线程处理
	MsgSnId() int
}

func (this *Coroutine) init() {
	this.queue = *list.New()
}

// 启动协程，持续处理消息 会阻塞 应当以  go Start() 启动
func (this *Coroutine) Start() {

	count := 0
	for {

		this.mutex.Lock()
		// 出队 (Dequeue)
		front := this.queue.Front()

		if front == nil {
			this.mutex.Unlock()
			//fmt.Println("队列为空")
			count = 0
			time.Sleep(time.Millisecond * 1)
		} else {
			this.queue.Remove(front)
			this.mutex.Unlock()
			if msg, ok := front.Value.(V1Msg); ok {

				this.Handler(msg)

			} else {
				fmt.Printf("类型错误,消息并不是 V1Msg %v ", front.Value)
			}
			count++
		}

		//处理10条数据就休眠
		if count%10 == 0 {
			time.Sleep(time.Millisecond * 1)
		}

	}
}

// 处理单条消息
func (this *Coroutine) Handler(msg V1Msg) {

	defer func() {
		if p := recover(); p != nil {
			fmt.Printf(" 捕获到 panic: %v\n", p)
			fmt.Printf("堆栈信息:\n%s\n", debug.Stack())
		}
	}()
	if this.MsgHandler.Fun == nil {
		//没有消息处理器，则默认输出信息
		fmt.Printf("  ###  handler info =  【id =%v, name = %v 】 \n	 msg =  %#v , 剩余数量 = %v  \n", this.Id, this.Name, msg, this.queue.Len())
	} else {
		this.MsgHandler.Fun(msg)
	}
}

// 接受消息
func (this *Coroutine) Accept(msg V1Msg) {
	defer this.mutex.Unlock()
	this.mutex.Lock()
	this.queue.PushBack(msg)
	fmt.Printf("  		接收到新消息  当前未处理消息数量 = %v  \n", this.queue.Len())

}
