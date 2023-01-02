package database

import (
	"fmt"
	"sync"

	"gorm.io/gorm"
)

type ConfTable struct {
	gorm.Model
	ChannelSize int    `json:"channelsize" gorm:"<-"`
	ReportSize  int    `json:"reportsize" gorm:"<-"`
	GlobalCode  []byte `json:"globalcode" gorm:"<-"`
}

type parmConfig struct {
	ChannelSize int
	ReportSize  int
	GlobalCode  []byte
}

type OptionConfig func(*parmConfig)

type Conf struct {
	db *gorm.DB
	sync.Mutex
}

func WithChannelSize(channelSize int) OptionConfig {
	return func(p *parmConfig) {
		p.ChannelSize = channelSize
	}
}

func (c *Conf) Rmv(parm ...OptionConfig) {

}

func (c *Conf) Find() (ConfTable, error) {
	info := ConfTable{}
	res := c.db.Where("name = ?", "sysconfig").First(&info)

	return info, res.Error
}

func (c *Conf) Update(parm ...OptionConfig) {
	p := &parmConfig{}
	for _, opt := range parm {
		opt(p)
	}

	info := ConfTable{}
	res := c.db.Where("name = ?", "sysconfig").First(&info)

	if p.ChannelSize != 0 {
		info.ChannelSize = p.ChannelSize
	}

	if p.ReportSize != 0 {
		info.ReportSize = p.ReportSize
	}

	if len(p.GlobalCode) != 0 {
		info.GlobalCode = p.GlobalCode
	}

	if res.Error == nil {
		res = c.db.Model(&ConfTable{}).Where("name = ?", "sysconfig").Updates(info)
	} else if res.Error == gorm.ErrRecordNotFound {
		res = c.db.Create(info)
	}

	if res.Error != nil {
		fmt.Println("conf update err", res.Error)
	}
}
