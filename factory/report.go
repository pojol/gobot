package factory

import (
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/pojol/gobot/database"
)

type urlDetail struct {
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

	UrlMap map[string]*urlDetail
}

type Report struct {
	Arr   []database.ReportInfo
	Limit int32
}

func NewReport(limit int32) *Report {

	if limit == 0 {
		panic(errors.New("report limit cannot be zero!"))
	}

	return &Report{
		Limit: limit,
		Arr:   database.Get().GetReport(),
	}
}

func (r *Report) Append(info ReportDetail) error {

	var err error

	if len(r.Arr) >= int(r.Limit) {
		err = database.Get().RemoveReport(r.Arr[0].ID)
		if err != nil {
			fmt.Println("append report remove limit", err.Error())
		}
		r.Arr = r.Arr[1:]
	}

	ri := database.ReportInfo{
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
		ri.ApiInfoLst = append(ri.ApiInfoLst, database.ReportApiInfo{
			Api:        fmtapi,
			ReqNum:     detail.ReqNum,
			ConsumeNum: int64(detail.AvgNum / int64(detail.ReqNum)),
			ReqSize:    detail.ReqSize,
			ResSize:    detail.ResSize,
			ErrNum:     detail.ErrNum,
		})
	}

	err = database.Get().AppendReport(ri)
	if err != nil {
		fmt.Println("append report", err.Error())
		return err
	}
	r.Arr = append(r.Arr, ri)

	return nil
}

func (r *Report) Info() []database.ReportInfo {

	desc := []database.ReportInfo{}

	if len(r.Arr) > 0 {
		for i := len(r.Arr) - 1; i >= 0; i -= 1 {
			desc = append(desc, r.Arr[i])
		}
	}

	return desc
}
