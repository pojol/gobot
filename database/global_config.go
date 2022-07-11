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

func GetGlobalScript(db IDatabase) []string {
	globalScript := []string{}

	cfglst, err := db.ConfigList()
	if err != nil {
		fmt.Println("get config list err", err.Error())
		return globalScript
	}

	for _, v := range cfglst {

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
