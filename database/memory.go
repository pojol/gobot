package database

import (
	"fmt"
	"sync"
	"time"

	"github.com/pojol/gobot/bot"
)

type MemoryAdapter struct {
	behaviormap map[string]BehaviorInfo
	configmap   map[string]TemplateConfig
	reportlst   []ReportInfo

	sync.Mutex
}

const (
	Momory = "momory"
)

func init() {
	Register(&MemoryAdapter{}, Momory)
}

func (f *MemoryAdapter) Init() error {

	f.behaviormap = make(map[string]BehaviorInfo)
	f.configmap = make(map[string]TemplateConfig)

	f.configmap["config"] = TemplateConfig{
		Name: "config",
		Dat:  []byte(`[{"title":"Global","content":"\n--[[\n\tGlobal constant area, users can define some constants here; it is easy to call in other scripts\n]]--\n\nREMOTE = \"http://127.0.0.1:8888\"\n","key":"global","closable":false},{"title":"HTTP","content":"\nlocal parm = {\n    body = {},    -- request body\n    timeout = \"10s\",\n    headers = {},\n}\n\nlocal url = REMOTE .. \"/group/methon\"\nlocal http = require(\"http\")\n\nfunction execute()\n    res, errmsg = http.post(url, parm)\n  \tif errmsg ~= nil then\n\t\tmeta.Err = errmsg\n    \treturn\n  \tend\n  \t\n  \tif res[\"status_code\"] ~= 200 then\n\t\tmeta.Err = \"post \" .. url .. \" http status code err \" .. res[\"status_code\"]\n  \t\treturn\n  \tend\n  \n  \tbody = json.decode(res[\"body\"])\n  \tmerge(meta, body.Body)\n\nend\n","key":"http","closable":false}]`),
	}

	fmt.Println("memory init succ")

	return nil
}

func (f *MemoryAdapter) UpsetFile(name string, byt []byte) error {

	f.Lock()
	defer f.Unlock()

	if _, ok := f.behaviormap[name]; ok {
		info := f.behaviormap[name]
		info.Dat = byt
		f.behaviormap[name] = info
	} else {
		f.behaviormap[name] = BehaviorInfo{
			Name:       name,
			Dat:        byt,
			Status:     bot.BotStatusUnknow,
			UpdateTime: time.Now().Unix(),
		}
	}

	return nil
}

func (f *MemoryAdapter) DelFile(name string) error {

	f.Lock()
	defer f.Unlock()

	delete(f.behaviormap, name)

	return nil
}

func (f *MemoryAdapter) FindFile(name string) (BehaviorInfo, error) {

	if _, ok := f.behaviormap[name]; ok {
		return f.behaviormap[name], nil
	}

	return BehaviorInfo{}, fmt.Errorf("cant find behavior %v", name)
}

func (f *MemoryAdapter) GetAllFiles() ([]BehaviorInfo, error) {

	lst := []BehaviorInfo{}

	for k := range f.behaviormap {
		lst = append(lst, f.behaviormap[k])
	}

	return lst, nil
}

func (f *MemoryAdapter) UpdateState(name string, status string) error {

	f.Lock()
	defer f.Unlock()

	if _, ok := f.behaviormap[name]; ok {
		info := f.behaviormap[name]
		info.Status = status
		f.behaviormap[name] = info
	}

	return nil
}

func (f *MemoryAdapter) UpdateTags(name string, tags []byte) error {

	f.Lock()
	defer f.Unlock()

	if _, ok := f.behaviormap[name]; ok {
		info := f.behaviormap[name]
		info.TagDat = tags
		f.behaviormap[name] = info
	}

	return nil
}

func (f *MemoryAdapter) FindConfig(name string) (TemplateConfig, error) {
	info := TemplateConfig{}

	f.Lock()
	defer f.Unlock()

	if _, ok := f.configmap[name]; ok {
		return f.configmap[name], nil
	}

	return info, fmt.Errorf("cant find config %v", name)
}

func (f *MemoryAdapter) UpsetConfig(byt []byte) error {

	f.Lock()
	defer f.Unlock()

	name := "config"

	if _, ok := f.configmap[name]; ok {
		info := f.configmap[name]
		info.Dat = byt
		f.configmap[name] = info
	} else {
		f.configmap[name] = TemplateConfig{
			Name: name,
			Dat:  byt,
		}
	}

	return nil
}

func (f *MemoryAdapter) RemoveReport(id string) error {

	f.Lock()
	defer f.Unlock()

	for k, v := range f.reportlst {
		if v.ID == id {
			f.reportlst = append(f.reportlst[:k], f.reportlst[k+1:]...)
			break
		}
	}

	return nil
}

func (f *MemoryAdapter) AppendReport(info ReportInfo) error {

	f.Lock()
	defer f.Unlock()

	f.reportlst = append(f.reportlst, info)

	return nil
}

func (f *MemoryAdapter) GetReport() []ReportInfo {
	return f.reportlst
}
