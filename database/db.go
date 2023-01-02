package database

import (
	"database/sql/driver"
	"encoding/json"
	"sync"

	"gorm.io/gorm"
)

type BehaviorInfo struct {
	gorm.Model
	Name       string `gorm:"<-"`
	Dat        []byte `gorm:"<-"`
	UpdateTime int64  `gorm:"<-"`
	Status     string `gorm:"<-"`
	TagDat     []byte `gorm:"<-"`
}

type BotTemplateConfig struct {
	gorm.Model

	Name string `gorm:"<-"`
	Tpl  []byte `gorm:"<-"`
}

type BotConfig struct {
	gorm.Model

	Name string `gorm:"<-"`
	Addr string `gorm:"<-"` // bot driver address
}

type ReportApiInfo struct {
	Api        string
	ReqNum     int
	ErrNum     int
	ConsumeNum int64

	ReqSize int64
	ResSize int64
}

type ReportApiArr []ReportApiInfo

func (p ReportApiArr) Value() (driver.Value, error) {
	return json.Marshal(p)
}

func (p *ReportApiArr) Scan(data interface{}) error {
	return json.Unmarshal(data.([]byte), &p)
}

type ReportInfo struct {
	gorm.Model
	ID         string
	Name       string
	BotNum     int
	ReqNum     int
	ErrNum     int
	Tps        int
	Dura       string
	BeginTime  int64
	ApiInfoLst ReportApiArr `gorm:"column:childrens;type:longtext"`
}

var (
	DefaultConfig = map[string]string{
		"global": `{"title":"Global","content":"\n--[[\n\tGlobal constant area, users can define some constants here; it is easy to call in other scripts\n]]--\n\nREMOTE = \"http://127.0.0.1:8888\"\n","closable":false, "prefab":false}`,
		"http":   `{"title":"HTTP","content":"\nlocal parm = {\n    body = {},    -- request body\n    timeout = \"10s\",\n    headers = {},\n}\n\nlocal url = REMOTE .. \"/group/methon\"\nlocal http = require(\"http\")\n\nfunction execute()\n    res, errmsg = http.post(url, parm)\n  \tif errmsg ~= nil then\n\t\tmeta.Err = errmsg\n    \treturn\n  \tend\n  \t\n  \tif res[\"status_code\"] ~= 200 then\n\t\tmeta.Err = \"post \" .. url .. \" http status code err \" .. res[\"status_code\"]\n  \t\treturn\n  \tend\n  \n  \tbody = json.decode(res[\"body\"])\n  \tmerge(meta, body.Body)\n\nend\n","closable":false, "prefab":true}`,
		"system": `{"channelsize":1024, "reportsize":100}`,
	}
)

type Cache struct {
	conf Conf
}

var db *Cache
var once sync.Once

func Create() *Cache {
	once.Do(func() {
		db = &Cache{}
	})

	return db
}

func init() {
	db = &Cache{}
}
