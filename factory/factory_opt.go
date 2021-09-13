package factory

import (
	"net/http"
	"time"
)

// 机器人的运行模式
const (
	FactoryModeStatic   = "static"
	FactoryModeIncrease = "increase"
)

// Parm 机器人工厂可配置参数
type Parm struct {
	// lifeTime 工厂的生命周期
	//
	// 默认值 1分钟
	lifeTime time.Duration

	// frameRate 机器人工厂的运行频率，（每秒创建多少个机器人
	//
	// 默认值 1s
	frameRate time.Duration

	// tickCreateNum 机器人工厂每个频率创建的数量
	//
	// 默认值 1
	tickCreateNum int

	// mode 机器人工厂的运行模式
	//
	// FactoryModeStatic 静态模式，这种模式将只会执行第一帧，通常用于一次性运行若干机器人
	//
	// FactoryModeIncrease 累增模式，这种模式下会按频率不断创建机器人，并在生命周期到时销毁改机器人
	//
	// 默认值 FactoryModeStatic
	mode string

	// Interrupt 当card遇到err的时候是否中断整个程序 （默认为否
	Interrupt bool

	// pickMode 策略选取模式
	pickMode string

	// matchUrl 匹配路由列表
	matchUrl []string

	// client http client
	//
	// 如果没有调用 WithClient factory会创建一个默认的client
	client *http.Client

	// batchSize 批次大小（用于控制goroutine的并发数量（默认1024
	batchSize int

	//
	md interface{}
}

// Option consul discover config wrapper
type Option func(*Parm)
