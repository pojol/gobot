package factory

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"sync"
	"time"

	"github.com/pojol/braid-go"
	"github.com/pojol/braid-go/module"
	"github.com/pojol/braid-go/module/meta"
	"github.com/pojol/gobot/driver/constant"
	"github.com/pojol/gobot/driver/database"
	script "github.com/pojol/gobot/driver/script/module"
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
					r.lst[k].Record(1, batchreport.Reports)

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
				BeginTime: time.Now().Unix(),
				ApiMap:    make(map[string]*database.ApiDetail),
			},
		})

	}

	return nil
}

func (r *FactoryReport) Close() {
	r.createbatchchan.Close()
}

func (b *Report) Record(botnum int32, report []script.Report) {
	b.rep.BotNum += botnum
	b.recordcnt++

	b.rep.MatchNum += int32(len(report))
	for _, v := range report {
		if _, ok := b.rep.ApiMap[v.MsgID]; !ok {
			b.rep.ApiMap[v.MsgID] = &database.ApiDetail{}
		}

		b.rep.ApiMap[v.MsgID].MatchCnt++
		if v.Err != "" {
			b.rep.ErrNum++
			b.rep.ApiMap[v.MsgID].ErrCnt++
		}
	}
}

func (b *Report) Generate() {

	fmt.Println("+--------------------------------------------------------------------------------------------------------+")
	fmt.Printf("Req url%-33s Req count %-5s\n", "", "")

	arr := []string{}
	for k := range b.rep.ApiMap {
		arr = append(arr, k)
	}
	sort.Strings(arr)

	var reqtotal int64

	for _, sk := range arr {
		v := b.rep.ApiMap[sk]

		reqtotal += int64(v.MatchCnt)

		u, _ := url.Parse(sk)
		fmt.Printf("%-40s %-15d\n", u.Path, v.MatchCnt)

	}
	fmt.Println("+--------------------------------------------------------------------------------------------------------+")

	fmt.Printf("robot : %d match to %d APIs req count : %d errors : %d\n", b.rep.BotNum, len(b.rep.ApiMap), b.rep.MatchNum, b.rep.ErrNum)

	database.GetReport().Append(*b.rep)
}
