package factory

import (
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

	// Interrupt 当card遇到err的时候是否中断整个程序 （默认为否
	Interrupt bool

	// batchSize 批次大小（用于控制goroutine的并发数量（默认1024
	batchSize int

	// 脚本路径
	ScriptPath string

	// 报告的次数限制
	ReportLimit int
}

// Option consul discover config wrapper
type Option func(*Parm)

func WithScriptPath(path string) Option {
	return func(c *Parm) {
		c.ScriptPath = path
	}
}

func WithReportLimit(limit int) Option {
	return func(c *Parm) {
		c.ReportLimit = limit
	}
}
