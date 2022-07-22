package database

import (
	"encoding/json"
	"fmt"
)

type CodeTemplateInfo struct {
	Title  string `json:"title"`
	Code   string `json:"content"`
	Prefab bool   `json:"prefab"`
}

type CodeTemplate struct {
	Lst []CodeTemplateInfo
}

type SystemInfo struct {
	ChannelSize int32 `json:"channelsize"`
	ReportSize  int32 `json:"reportsize"`
}

func GetGlobalScript(db IDatabase) []string {
	globalScript := []string{}

	cfglst, err := db.ConfigList()
	if err != nil {
		fmt.Println("get config list err", err.Error())
		return globalScript
	}

	for _, v := range cfglst {

		if v == "system" { // 关键字
			continue
		}

		cfg, err := db.ConfigFind(v)
		if err != nil {
			fmt.Println("find config err", err.Error())
			continue
		}

		info := CodeTemplateInfo{}
		err = json.Unmarshal(cfg.Dat, &info)
		if err != nil {
			fmt.Println("config unmarshal err", err.Error())
			continue
		}

		if !info.Prefab { // global script
			globalScript = append(globalScript, info.Code)
		}
	}

	return globalScript
}

func GetSystemParm(db IDatabase) SystemInfo {

	sysinfo := SystemInfo{
		ReportSize:  100,
		ChannelSize: 128,
	}

	cfg, err := db.ConfigFind("system")
	if err != nil {
		fmt.Println("find config err", err.Error())
		goto EXT
	}

	err = json.Unmarshal(cfg.Dat, &sysinfo)
	if err != nil {
		fmt.Println("config unmarshal err", err.Error())
		goto EXT
	}

EXT:
	return sysinfo
}
