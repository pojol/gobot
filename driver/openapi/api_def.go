package openapi

import (
	"github.com/pojol/gobot/driver/database"
	"github.com/pojol/gobot/driver/factory"
)

// ------------------------ File ------------------------

// file.list
type FileListRes struct {
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

// file.setTags
type SetBehaviorTagsReq struct {
	Name    string
	NewTags []string
}

type SetBehaviorTagsRes struct {
	Bots []behaviorInfo
}

// ------------------------ Config ------------------------

// config.sys.info
type ConfigGetSysInfoRes struct {
	ReportSize   int
	ChannelSize  int
	EnqueneDelay int
}

// config.sys.set
type ConfigSetSysInfoReq struct {
	ReportSize   int
	ChannelSize  int
	EnqueneDelay int
}

type ConfigSetSysInfoRes struct {
	ReportSize   int
	ChannelSize  int
	EnqueneDelay int
}

type ConfigRemoveReq struct {
	Name string
}

// ------------------------ Report ------------------------

type ReportRes struct {
	Info []database.ReportTable
}

// ------------------------ Bot ------------------------

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

// ------------------------ Debug ------------------------

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

// ------------------------ Prefab ------------------------

// prefab.list
type PrefabListReq struct {
}

type PrefabListRes struct {
	Lst []prefabInfo
}

// prefab.rmv
type PrefabRmvReq struct {
	Name string `json:"name"`
}

// prefab.setTags
type PrefabSetTagsReq struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

// ------------------------ runtime ------------------------
type RuntimeInfoReq struct {
	ID string
}

type RuntimeInfoRes struct {
	Msg string
}
