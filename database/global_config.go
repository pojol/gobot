package database

import (
	"encoding/json"
	"fmt"
)

type CodeTemplateInfo struct {
	Title string `json:"title"`
	Code  string `json:"content"`
}

type CodeTemplate struct {
	Lst []CodeTemplateInfo
}

func GetGlobalScript(db IDatabase) string {
	globalScript := ""
	temp := CodeTemplate{}
	tc, err := db.ConfigFind("config")
	if err != nil {
		fmt.Println("code template load err", err.Error())
	} else {
		err = json.Unmarshal(tc.Dat, &temp.Lst)
		if err != nil {
			fmt.Println("code template unmarshal err", err.Error())
		}
		for _, v := range temp.Lst {
			if v.Title == "Global" {
				globalScript = v.Code
				break
			}
		}
	}

	return globalScript
}
