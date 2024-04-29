package database

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"net/url"
	"sync"

	"gorm.io/gorm"
)

type ApiDetail struct {
	MatchCnt int32
	ErrCnt   int32
}

type ReportDetail struct {
	ID   string
	Name string

	MatchNum int32
	ErrNum   int32

	BeginTime int64 // 队列的开始时间
	BotNum    int32 // 一个执行队列中机器人总数量

	ApiMap map[string]*ApiDetail
}

////////////////////////////////////////////////////////

type ReportApiInfo struct {
	Api    string
	ReqNum int
	ErrNum int
}

type ReportApiArr []ReportApiInfo

func (p ReportApiArr) Value() (driver.Value, error) {
	return json.Marshal(p)
}

func (p *ReportApiArr) Scan(data interface{}) error {
	return json.Unmarshal(data.([]byte), &p)
}

type ReportTable struct {
	gorm.Model
	ID         string
	Name       string
	BotNum     int32
	ReqNum     int32
	ErrNum     int32
	BeginTime  int64
	ApiInfoLst ReportApiArr `gorm:"column:childrens;type:longtext"`
}

type Report struct {
	db *gorm.DB
	sync.Mutex
}

func CreateReport(mysqlptr *gorm.DB) *Report {
	r := &Report{
		db: mysqlptr,
	}

	err := r.db.AutoMigrate(&ReportTable{})
	if err != nil {
		fmt.Println("migrate err", err.Error())
	}

	return r
}

func (r *Report) Append(info ReportDetail) error {

	ri := ReportTable{
		ID:        info.ID,
		Name:      info.Name,
		BotNum:    info.BotNum,
		ErrNum:    info.ErrNum,
		BeginTime: info.BeginTime,
	}
	for api, detail := range info.ApiMap {
		u, err := url.Parse(api)
		fmtapi := ""
		if err == nil {
			fmtapi = u.Path
		}
		ri.ApiInfoLst = append(ri.ApiInfoLst, ReportApiInfo{
			Api:    fmtapi,
			ReqNum: int(detail.MatchCnt),
			ErrNum: int(detail.ErrCnt),
		})
	}

	return r.db.Model(&ReportTable{}).Create(&ri).Error
}

func (r *Report) List() ([]ReportTable, error) {
	var lst []ReportTable

	res := r.db.Order("begin_time desc").Limit(100).Find(&lst)

	return lst, res.Error
}
