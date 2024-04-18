package database

import (
	"fmt"
	"sync"
	"time"

	"gorm.io/gorm"
)

type BehaviorTable struct {
	gorm.Model
	Name       string `gorm:"<-"`
	File       []byte `gorm:"<-"`
	UpdateTime int64  `gorm:"<-"`
	Status     string `gorm:"<-"`
	Tags       []byte `gorm:"<-"`
}

type Behavior struct {
	db *gorm.DB
	sync.Mutex
}

func CreateBehavior(mysqlptr *gorm.DB) *Behavior {
	b := &Behavior{
		db: mysqlptr,
	}

	err := b.db.AutoMigrate(&BehaviorTable{})
	if err != nil {
		fmt.Println("migrate err", err.Error())
	}

	return b
}

func (b *Behavior) List() ([]BehaviorTable, error) {
	lst := []BehaviorTable{}

	res := b.db.Find(&lst)

	return lst, res.Error
}

func (b *Behavior) Find(name string) (BehaviorTable, error) {
	t := BehaviorTable{}

	res := b.db.Where("name = ?", name).First(&t)

	return t, res.Error
}

func (b *Behavior) Upset(name string, file []byte) {
	var t BehaviorTable
	var err error
	var res *gorm.DB

	t, err = b.Find(name)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			b.db.Model(&BehaviorTable{}).Create(&BehaviorTable{
				Name:       name,
				File:       file,
				UpdateTime: time.Now().Unix(),
				Status:     "unknow",
			})
		} else {
			fmt.Println("behavior upset err", err.Error())
		}
	} else {
		t.UpdateTime = time.Now().Unix()
		t.File = make([]byte, len(file))
		copy(t.File, file)
		res = b.db.Model(&BehaviorTable{}).Where("name = ?", name).Updates(&t)
		if res.Error != nil {
			fmt.Println("behavior upset update err", res.Error)
		}
	}

}

func (b *Behavior) Rmv(name string) error {

	info := BehaviorTable{}
	res := b.db.Where("name = ?", name).Delete(&info)

	return res.Error
}

func (b *Behavior) UpdateTags(name string, tags []byte) error {
	return b.db.Model(&BehaviorTable{}).Where("name = ?", name).Update("tags", tags).Error
}

func (b *Behavior) UpdateStatus(name string, status string) error {
	return b.db.Model(&BehaviorTable{}).Where("name = ?", name).Update("status", status).Error
}
