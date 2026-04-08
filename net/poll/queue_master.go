package pool

import (
	"math"
	"strconv"
)

// 队列主入口, 职责 : 接收任务  + 分发
type QueueMaster struct {
	Slave []*QueueSlave
	Name  string
	Id    int
}

// 业务处理器
type Handler struct {
	Fun func(task *MessageTask)
}

// 创建协程组
func NewCoroutineGroup(id int, name string, slaveNum int) *QueueMaster {

	if slaveNum < 1 {
		slaveNum = 1 //设置最小的线程数量=2
	}
	ms := QueueMaster{}
	ms.Id = id
	ms.Name = name

	ms.Slave = make([]*QueueSlave, slaveNum)

	for i := 0; i < slaveNum; i++ {
		t := QueueSlave{Id: i, Name: "slave-" + strconv.Itoa(i)}
		t.Init()
		ms.Slave[i] = &t
		go ms.Slave[i].Start()
	}

	return &ms
}

// 接受消息
func (this *QueueMaster) Accept(task *MessageTask) {
	//取模后子线程入队列
	var m = task.SlaveId() % len(this.Slave)
	if m < 0 {
		m = int(math.Abs(float64(m)))
	}

	slave := this.Slave[m]
	slave.Accept(task)
}

// 添加处理器
func (this *QueueMaster) Handler(fun func(task *MessageTask)) QueueMaster {

	if fun != nil {
		hdl := Handler{Fun: fun}
		for _, slave := range this.Slave {
			slave.MsgHandler = hdl
		}
	}
	return *this
}
