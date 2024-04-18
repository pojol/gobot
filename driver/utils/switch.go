package utils

// from https://github.com/grpc/grpc-go/blob/master/internal/grpcsync/event.go

import (
	"sync/atomic"
)

// Switch 用于模拟开关行为的事件结构
type Switch struct {
	opened int32
	c      chan struct{}
}

// Open 开启
func (s *Switch) Open() {
	atomic.StoreInt32(&s.opened, 1)
	s.c <- struct{}{}
}

func (s *Switch) Close() {
	atomic.StoreInt32(&s.opened, 0)
}

// Done 事件触发
func (s *Switch) Done() <-chan struct{} {
	return s.c
}

// HasOpend 是否以开启
func (s *Switch) HasOpend() bool {
	return atomic.LoadInt32(&s.opened) == 1
}

// NewSwitch 返回一个新的开关事件
func NewSwitch() *Switch {
	return &Switch{c: make(chan struct{}, 1)}
}
