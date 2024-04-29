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

	// Interrupt 当card遇到err的时候是否中断整个程序 （默认为否
	Interrupt bool

	// 脚本路径
	ScriptPath string

	// 报告的次数限制
	ReportLimit int

	// 服务id
	ServiceID string

	//
	DBType string

	// cluster mode
	ClusterMode bool
}

// Option consul discover config wrapper
type Option func(*Parm)

func WithScriptPath(path string) Option {
	return func(c *Parm) {
		c.ScriptPath = path
	}
}

func WithServiceID(id string) Option {
	return func(c *Parm) {
		c.ServiceID = id
	}
}

func WithClusterMode(open bool) Option {
	return func(c *Parm) {
		c.ClusterMode = open
	}

}

func WithReportLimit(limit int) Option {
	return func(c *Parm) {
		c.ReportLimit = limit
	}
}

func WithDatabase(db string) Option {
	return func(c *Parm) {
		c.DBType = db
	}
}
