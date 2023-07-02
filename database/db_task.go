package database

import (
	"fmt"
	"sync"

	"gorm.io/gorm"
)

type TaskTable struct {
	gorm.Model
	ID          string `gorm:"<-"`
	Name        string `gorm:"<-"`
	TotalNumber int32  `gorm:"<-"`
	CurNumber   int32  `gorm:"<-"`
}

type Task struct {
	db *gorm.DB
	sync.Mutex
}

func CreateTask(mysqlptr *gorm.DB) *Task {
	b := &Task{
		db: mysqlptr,
	}

	err := b.db.AutoMigrate(&TaskTable{})
	if err != nil {
		panic(err)
	}

	return b
}

func (b *Task) New(tt TaskTable) {
	b.db.Model(&TaskTable{}).Create(&tt)
}

func (b *Task) List() ([]TaskTable, error) {
	lst := []TaskTable{}

	res := b.db.Find(&lst)

	return lst, res.Error
}

func (b *Task) Rmv(id string) error {
	var tt TaskTable
	return b.db.Model(&tt).Where("id = ?", id).Delete(&tt).Error
}

func (b *Task) Update(id string, cur int32) error {
	var tt TaskTable
	fmt.Println("update task", id, cur)
	if id == "" || cur < 0 {
		return nil
	}

	res := b.db.Model(&tt).Where("id = ?", id).First(&tt)
	if res.Error != nil {
		return res.Error
	}

	if cur > tt.TotalNumber {
		b.db.Model(&tt).Where("id = ?", id).Delete(&tt)
		return nil
	}

	return b.db.Model(&tt).Where("id = ?", id).Update("cur_number", cur).Error
}
