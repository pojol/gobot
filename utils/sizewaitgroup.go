package utils

import (
	"context"
	"fmt"
	"math"
	"sync"
)

// SizeWaitGroup 用于控制goroutine的并发数量
type SizeWaitGroup struct {
	Size    int
	wg      sync.WaitGroup
	blockch chan struct{}
}

// 创建一个固定大小的拥塞队列
func New(size int) SizeWaitGroup {
	if size <= 0 || size > math.MaxInt16 {
		panic(fmt.Errorf("not allow size %v", size))
	}

	return SizeWaitGroup{
		Size:    size,
		wg:      sync.WaitGroup{},
		blockch: make(chan struct{}, size),
	}
}

// Add 队列+1
func (swg *SizeWaitGroup) Add() {
	swg.EnterWithContext(context.Background())
}

func (swg *SizeWaitGroup) EnterWithContext(ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	case swg.blockch <- struct{}{}:
		break
	}

	swg.wg.Add(1)
}

// Done 队列-1
func (swg *SizeWaitGroup) Done() {
	<-swg.blockch
	swg.wg.Done()
}

// Wait 阻塞等待队列全部完成
func (swg *SizeWaitGroup) Wait() {
	swg.wg.Wait()
}
