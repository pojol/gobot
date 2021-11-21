package server

import (
	"github.com/pojol/gobot/database"
	"github.com/pojol/gobot/factory"
)

type Response struct {
	Code int
	Msg  string
	Body interface{}
}

// file.list
type behaviorInfo struct {
	Name   string
	Update int64
	Status string
	Tags   []string
	Desc   string
}

type BehaviorListRes struct {
	Bots []behaviorInfo
}

// file.remove
type FileRemoveReq struct {
	Name string
}
type FileRemoveRes struct {
	Bots []behaviorInfo
}

// file.get
type FindBehaviorReq struct {
	Name string
}

type FindBehaviorRes struct {
	Info database.BehaviorInfo
}

// file.setTags
type SetBehaviorTagsReq struct {
	Name    string
	NewTags []string
}

type SetBehaviorTagsRes struct {
	Bots []behaviorInfo
}

// get.report
type ReportApiInfo struct {
	Api        string
	ReqNum     int
	ErrNum     int
	ConsumeNum int64

	ReqSize int64
	ResSize int64
}

type ReportInfo struct {
	ID        string
	Name      string
	BotNum    int
	ReqNum    int
	ErrNum    int
	Tps       int
	Dura      string
	BeginTime string
	Apilst    []ReportApiInfo
}

type ReportRes struct {
	Info []ReportInfo
}

// bot.create
type RunRequest struct {
	Name string
	Num  int
}

type RunResponse struct {
}

// bot.list
type BotListResponse struct {
	Lst []factory.BatchInfo
}

// debug.step
type StepRequest struct {
	BotID string
}

type StepResponse struct {
	Prev       string
	Cur        string
	Blackboard string
}

// debug.create
type CreateDebugBotResponse struct {
	BotID string
}
