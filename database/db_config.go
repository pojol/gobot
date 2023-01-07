package database

import (
	"fmt"
	"sync"

	lua "github.com/yuin/gopher-lua"
	"gorm.io/gorm"
)

type ConfTable struct {
	gorm.Model
	Name        string `json:"name" gorm:"<-"`
	ChannelSize int    `json:"channelsize" gorm:"<-"`
	ReportSize  int    `json:"reportsize" gorm:"<-"`
	GlobalCode  []byte `json:"globalcode" gorm:"<-"`
}

type Conf struct {
	db *gorm.DB
	sync.Mutex
}

func CreateConfig(mysqlptr *gorm.DB) *Conf {
	c := &Conf{
		db: mysqlptr,
	}

	err := c.db.AutoMigrate(&ConfTable{})
	if err != nil {
		fmt.Println("migrate err", err.Error())
	}

	_, err = c.Get()
	if err == gorm.ErrRecordNotFound {
		c.db.Create(&ConfTable{
			Name:        "sysconfig",
			ChannelSize: 512,
			ReportSize:  100,
			GlobalCode: []byte(`
--[[
	Global constant area, users can define some constants here; it is easy to call in other scripts
]]--

REMOTE = "http://127.0.0.1:8888\"
			`),
		})
	}

	return c
}

func (c *Conf) Get() (ConfTable, error) {
	info := ConfTable{}
	res := c.db.Where("name = ?", "sysconfig").First(&info)

	if res.Error != nil {
		fmt.Println("conf find err", res.Error)
	}

	return info, res.Error
}

func (c *Conf) update(key string, val interface{}) {
	res := c.db.Model(&ConfTable{}).Where("name = ?", "sysconfig").Update(key, val)
	if res.Error != nil {
		fmt.Println("conf update err", res.Error)
	}
}

func (c *Conf) UpdateChannelSize(cs int) {
	c.Lock()
	defer c.Unlock()

	if cs <= 0 {
		fmt.Println("wrong input", cs)
	}

	_, err := c.Get()
	if err != nil {
		return
	}

	c.update("channel_size", cs)
}

func (c *Conf) UpdateReportSize(rs int) {
	c.Lock()
	defer c.Unlock()

	if rs <= 0 {
		fmt.Println("wrong input", rs)
		return
	}

	_, err := c.Get()
	if err != nil {
		return
	}

	c.update("report_size", rs)
}

func (c *Conf) UpdateGlobalDefine(code []byte) error {
	c.Lock()
	defer c.Unlock()

	if len(code) == 0 {
		fmt.Println("wrong input", string(code))
		return fmt.Errorf("wrong input %v", string(code))
	}

	L := lua.NewState()
	_, err := L.LoadString(string(code))
	if err != nil {
		return fmt.Errorf("parse script err %w", err)
	}

	_, err = c.Get()
	if err != nil {
		return err
	}

	c.update("global_code", code)
	return nil
}
