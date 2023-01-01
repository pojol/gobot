package server

import (
	"github.com/pojol/gobot/config"
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

type SetSystemConfigReq struct {
	Cfg config.SystemCfg
}

type SetConfigReq struct {
	Name string `json:"name"`
	Dat  string `json:"dat"`
}

type ReportRes struct {
	Info []database.ReportInfo
}

//
type ConfigGetListInfoRes struct {
	Lst []string
}

type ConfigRemoveReq struct {
	Name string
}

// bot.batch
type BotBatchCreateRequest struct {
	Name string
	Num  int
}

type BotBatchCreateResponse struct {
}

// bot.run
type BotRunRequest struct {
	Name string
}

type BotRunResponse struct {
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
	Blackboard string
	ThreadInfo string
}

// debug.create
type CreateDebugBotResponse struct {
	BotID      string
	Blackboard string
	ThreadInfo string
}

type ConfigGetInfoResponse struct {
	Lst []database.TemplateConfig
}
