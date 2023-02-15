package database

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"net/url"
	"sync"
	"time"

	"gorm.io/gorm"
)

type ApiDetail struct {
	ReqNum int
	ErrNum int
	AvgNum int64

	ReqSize int64
	ResSize int64
}

type ReportDetail struct {
	ID     string
	Name   string
	BotNum int
	ReqNum int
	ErrNum int
	Tps    int
	Dura   string

	BeginTime time.Time

	UrlMap map[string]*ApiDetail
}

////////////////////////////////////////////////////////

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

type ReportTable struct {
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
		ReqNum:    info.ReqNum,
		ErrNum:    info.ErrNum,
		Tps:       info.Tps,
		Dura:      info.Dura,
		BeginTime: info.BeginTime.Unix(),
	}
	for api, detail := range info.UrlMap {
		u, err := url.Parse(api)
		fmtapi := ""
		if err == nil {
			fmtapi = u.Path
		}
		ri.ApiInfoLst = append(ri.ApiInfoLst, ReportApiInfo{
			Api:        fmtapi,
			ReqNum:     detail.ReqNum,
			ConsumeNum: int64(detail.AvgNum / int64(detail.ReqNum)),
			ReqSize:    detail.ReqSize,
			ResSize:    detail.ResSize,
			ErrNum:     detail.ErrNum,
		})
	}

	return r.db.Model(&ReportTable{}).Create(&ri).Error
}

func (r *Report) List() ([]ReportTable, error) {
	var lst []ReportTable

	res := r.db.Find(&lst).Limit(100)

	return lst, res.Error
}
