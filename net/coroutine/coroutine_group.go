package coroutine

import (
	"math"
	"strconv"
)

// 协程组
type CoroutineGroup struct {
	Slave []*Coroutine
	Name  string
	Id    int
}

// 业务处理器
type Handler struct {
	Fun func(msg V1Msg)
}

// 创建协程组
func NewCoroutineGroup(id int, name string, slaveNum int) *CoroutineGroup {

	if slaveNum < 1 {
		slaveNum = 1 //设置最小的线程数量=2
	}
	ms := CoroutineGroup{}
	ms.Id = id
	ms.Name = name

	ms.Slave = make([]*Coroutine, slaveNum)

	for i := 0; i < slaveNum; i++ {
		t := Coroutine{Id: i, Name: "slave-" + strconv.Itoa(i)}
		ms.Slave[i] = &t
		go ms.Slave[i].Start()
	}

	return &ms
}

// 接受消息
func (this *CoroutineGroup) Accept(msg V1Msg) {
	//取模后子线程入队列
	var m = msg.MsgSnId() % len(this.Slave)
	if m < 0 {
		m = int(math.Abs(float64(m)))
	}

	slave := this.Slave[m]
	slave.Accept(msg)
}

// 添加处理器
func (this *CoroutineGroup) Handler(fun func(msg V1Msg)) CoroutineGroup {

	if fun != nil {
		hdl := Handler{Fun: fun}
		for _, slave := range this.Slave {
			slave.MsgHandler = hdl
		}
	}
	return *this
}
