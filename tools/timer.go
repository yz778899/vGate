package tools

import (
	"fmt"
	"math/rand"
	"time"
)

// 时间消耗
type TimeExpend struct {
	T *time.Time
}

// 时间消耗
func NewTimeExpend() *TimeExpend {
	now := time.Now()
	return &TimeExpend{T: &now}
}

// 重置时间
func (this *TimeExpend) Reset() *TimeExpend {
	now := time.Now()
	this.T = &now
	return this
}

// 计时并返回时长
func (this *TimeExpend) End() time.Duration {
	return time.Since(*this.T)
}

// 打印耗时信息 比如 ： EndAndPrint("消耗时长")
func (this *TimeExpend) EndAndPrint(tmplet string) {

	fmt.Printf(tmplet+" %v \n", formatDurationSmart(this.End()))
}

func formatDurationSmart(d time.Duration) string {
	switch {
	case d < time.Microsecond:
		return fmt.Sprintf("%d 纳秒", d.Nanoseconds())
	case d < time.Millisecond:
		return fmt.Sprintf("%.2f 微秒", float64(d.Nanoseconds())/1000)
	case d < time.Second:
		return fmt.Sprintf("%.2f 毫秒", float64(d.Milliseconds()))
	case d < time.Minute:
		return fmt.Sprintf("%.2f 秒", d.Seconds())
	case d < time.Hour:
		minutes := int(d.Minutes())
		seconds := int(d.Seconds()) % 60
		return fmt.Sprintf("%d 分 %d 秒", minutes, seconds)
	default:
		hours := int(d.Hours())
		minutes := int(d.Minutes()) % 60
		return fmt.Sprintf("%d 小时 %d 分", hours, minutes)
	}
}

// 使用例子
func SimpleTimeExpend() {
	t := NewTimeExpend()
	mem := make([]byte, 1024*1024*5)
	t.EndAndPrint("申请内存,消耗时长")
	t.Reset()
	for x := 0; x < len(mem); x++ {
		mem[x] = byte(rand.Intn(255))
	}
	t.EndAndPrint("	写入内存,消耗时长")
}
