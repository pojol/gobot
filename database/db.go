package database

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
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

type TemplateConfig struct {
	gorm.Model
	Name string `gorm:"<-"`
	Dat  []byte `gorm:"<-"`
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

type IDatabase interface {
	Init() error

	UpsetFile(string, []byte) error
	DelFile(string) error
	FindFile(string) (BehaviorInfo, error)
	GetAllFiles() ([]BehaviorInfo, error)

	UpdateState(name string, status string) error
	UpdateTags(name string, tags []byte) error

	FindConfig(name string) (TemplateConfig, error)
	UpsetConfig(byt []byte) error

	RemoveReport(id string) error
	AppendReport(info ReportInfo) error
	GetReport() []ReportInfo
}

var registry = struct {
	sync.Mutex
	once    sync.Once
	dbpoint map[string]IDatabase
}{
	dbpoint: make(map[string]IDatabase),
}

func Register(component IDatabase, name string) {
	registry.Lock()
	defer registry.Unlock()

	if _, ok := registry.dbpoint[name]; !ok {
		registry.dbpoint[name] = component
	}
}

func Lookup(name string) IDatabase {
	if _, ok := registry.dbpoint[name]; ok {
		registry.once.Do(func() {
			err := registry.dbpoint[name].Init()
			if err != nil {
				panic(fmt.Errorf("loop up %v database fail %v", name, err.Error()))
			}
		})

		return registry.dbpoint[name]
	}

	return nil
}
