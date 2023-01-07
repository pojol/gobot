package database

import (
	"fmt"
	"sync"

	"gorm.io/gorm"
)

type PrefabTable struct {
	gorm.Model
	Name string `gorm:"<-"`
	Code []byte `gorm:"<-"`
	Tags []byte `gorm:"<-"`
}

type Prefab struct {
	db *gorm.DB
	sync.Mutex
}

func CreatePrefab(mysqlptr *gorm.DB) *Prefab {
	p := &Prefab{
		db: mysqlptr,
	}

	err := p.db.AutoMigrate(&PrefabTable{})
	if err != nil {
		fmt.Println("migrate err", err.Error())
	}

	return p
}

func (p *Prefab) List() ([]PrefabTable, error) {
	lst := []PrefabTable{}

	res := p.db.Find(&lst)

	return lst, res.Error
}

func (p *Prefab) Find(name string) (PrefabTable, error) {
	t := PrefabTable{}

	res := p.db.Where("name = ?", name).First(&t)

	return t, res.Error
}

func (b *Prefab) Upset(name string, code []byte) {
	var t PrefabTable
	var err error
	var res *gorm.DB

	t, err = b.Find(name)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			b.db.Model(&PrefabTable{}).Create(&PrefabTable{
				Name: name,
				Code: code,
			})
		} else {
			fmt.Println("behavior upset err", err.Error())
		}
	} else {
		t.Code = code
		res = b.db.Model(&PrefabTable{}).Where("name = ?", name).Updates(&t)
		if res.Error != nil {
			fmt.Println("behavior upset update err", res.Error)
		}
	}
}

func (p *Prefab) Rmv(name string) error {

	info := PrefabTable{}
	res := p.db.Where("name = ?", name).Delete(&info)

	return res.Error
}

func (p *Prefab) UpdateTags(name string, tags []byte) error {
	return p.db.Model(&PrefabTable{}).Where("name = ?", name).Update("tags", tags).Error
}
