package factory

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/pojol/braid-go"
	"github.com/pojol/braid-go/module"
	"github.com/pojol/braid-go/module/meta"
	"github.com/pojol/gobot/constant"
	"github.com/pojol/gobot/database"
	script "github.com/pojol/gobot/script/module"
)

type Report struct {
	rep         *database.ReportDetail
	recordcnt   int
	recordtotal int
	chann       module.IChannel
}

type FactoryReport struct {
	lst             []Report
	createbatchchan module.IChannel
	serviceid       string

	sync.Mutex
}

func BuildReport(id string) *FactoryReport {
	return &FactoryReport{
		serviceid: id,
	}
}

func (r *FactoryReport) watch() {
	var err error

	fmt.Println("watch bot.batch.create topic", "factory-"+r.serviceid)
	r.createbatchchan, err = braid.Topic("bot.batch.create").Sub(context.TODO(), "factory-"+r.serviceid)
	if err != nil {
		fmt.Println("factory.watch", err.Error())
		return
	}

	r.createbatchchan.Arrived(func(msg *meta.Message) error {
		info := &BotBatchInfo{}
		json.Unmarshal(msg.Body, info)

		if constant.GetServerState() == meta.EMaster {
			r.create(info.ID, info.Name, info.Num)
		}

		return nil
	})
}

func (r *FactoryReport) create(id, name string, cnt int) error {
	flag := false
	total := 0

	r.Lock()
	for k, v := range r.lst {
		if v.rep.ID == id {
			r.lst[k].recordtotal += cnt
			total = r.lst[k].recordtotal
			flag = true
			break
		}
	}
	r.Unlock()

	if flag {
		fmt.Println("add batch report topic", id, name, cnt, "=>", total)
	} else {
		fmt.Println("create batch report topic", id, name, cnt)

		batchchan, err := braid.Topic("batch.report").Sub(context.TODO(), "batch.report."+id)
		if err != nil {
			fmt.Println("batch watch err", err.Error())
			return err
		}

		batchchan.Arrived(func(msg *meta.Message) error {

			// 让主节点处理
			if constant.GetServerState() != meta.EMaster {
				return nil
			}

			batchreport := &BatchReport{}
			json.Unmarshal(msg.Body, batchreport)

			for k, v := range r.lst {
				if v.rep.ID == batchreport.ID {

					r.Lock()
					r.lst[k].Record(batchreport.Reports)

					if r.lst[k].recordcnt >= r.lst[k].recordtotal {
						v.Generate()
						r.lst[k].chann.Close()

						fmt.Println("batch", batchreport.ID, "record", r.lst[k].recordcnt, "succ")
					}
					r.Unlock()

					return nil
				}
			}

			return nil
		})

		r.lst = append(r.lst, Report{
			chann:       batchchan,
			recordtotal: cnt,
			rep: &database.ReportDetail{
				ID:        id,
				Name:      name,
				BeginTime: time.Now(),
				UrlMap:    make(map[string]*database.ApiDetail),
			},
		})

	}

	return nil
}

func (r *FactoryReport) Close() {
	r.createbatchchan.Close()
}

func (b *Report) Record(report []script.Report) {
	b.rep.BotNum++
	b.recordcnt++

	//fmt.Println("report", b.rep.ID, b.recordcnt, "=>", b.recordtotal)

	b.rep.ReqNum += len(report)
	for _, v := range report {
		if _, ok := b.rep.UrlMap[v.Api]; !ok {
			b.rep.UrlMap[v.Api] = &database.ApiDetail{}
		}

		b.rep.UrlMap[v.Api].ReqNum++
		b.rep.UrlMap[v.Api].AvgNum += int64(v.Consume)
		b.rep.UrlMap[v.Api].ReqSize += int64(v.ReqBody)
		b.rep.UrlMap[v.Api].ResSize += int64(v.ResBody)
		if v.Err != "" {
			b.rep.ErrNum++
			b.rep.UrlMap[v.Api].ErrNum++
		}
	}
}

func (b *Report) Generate() {

	fmt.Println("+--------------------------------------------------------------------------------------------------------+")
	fmt.Printf("Req url%-33s Req count %-5s Average time %-5s Body req/res %-5s Succ rate %-10s\n", "", "", "", "", "")

	arr := []string{}
	for k := range b.rep.UrlMap {
		arr = append(arr, k)
	}
	sort.Strings(arr)

	var reqtotal int64

	for _, sk := range arr {
		v := b.rep.UrlMap[sk]
		var avg string
		if v.AvgNum == 0 {
			avg = "0 ms"
		} else {
			avg = strconv.Itoa(int(v.AvgNum/int64(v.ReqNum))) + "ms"
		}

		succ := strconv.Itoa(v.ReqNum-v.ErrNum) + "/" + strconv.Itoa(v.ReqNum)

		reqsize := strconv.Itoa(int(v.ReqSize/1024)) + "kb"
		ressize := strconv.Itoa(int(v.ResSize/1024)) + "kb"

		reqtotal += int64(v.ReqNum)

		u, _ := url.Parse(sk)
		fmt.Printf("%-40s %-15d %-18s %-18s %-10s\n", u.Path, v.ReqNum, avg, reqsize+" / "+ressize, succ)

	}
	fmt.Println("+--------------------------------------------------------------------------------------------------------+")

	durations := int(time.Since(b.rep.BeginTime).Seconds())
	if durations <= 0 {
		durations = 1
	}

	qps := int(reqtotal / int64(durations))
	duration := strconv.Itoa(durations) + "s"

	b.rep.Tps = qps
	b.rep.Dura = duration
	fmt.Printf("robot : %d match to %d APIs req count : %d duration : %s qps : %d errors : %d\n", b.rep.BotNum, len(b.rep.UrlMap), b.rep.ReqNum, duration, qps, b.rep.ErrNum)

	database.GetReport().Append(*b.rep)
}
