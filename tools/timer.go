package tools

import "time"

type Timer struct {
	timer *time.Timer
}

func NewTimer() *Timer {
	return &Timer{}
}
