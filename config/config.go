package config

import (
	"github.com/pojol/gobot/database"
	lua "github.com/yuin/gopher-lua"
)

// 系统配置项
type SystemCfg struct {
	ChannelSize int `json:"channelsize"`
	ReportSize  int `json:"reportsize"`
}

type PrefabInfo struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
	Code []byte
}

// 预制节点配置
type prefabCfg struct {
	Arr []PrefabInfo
}

type Conf struct {
	db database.IDatabase

	sys        SystemCfg
	globalCode []byte

	prefab prefabCfg
}

func GetPrefabs() []PrefabInfo {
	arr := []PrefabInfo{}
	return arr
}

func GetSystemConfig() SystemCfg {
	sys := SystemCfg{}
	return sys
}

func SetSystemConfig(cfg SystemCfg) {

}

func GetGlobalDefine() []byte {
	return []byte{}
}

func SetGlobalDefine(code []byte) {
	L := lua.NewState()
	_, err := L.LoadString(string(code))
	if err != nil {

	}
}

// GetReportSize 获取报告的最大数量
func GetReportSize() int {
	return 0
}

// GetChannelSize 获取同时运行的机器人数量限制
func GetChannelSize() int {
	return 0
}
