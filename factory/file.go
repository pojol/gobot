package factory

import (
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type BehaviorInfo struct {
	gorm.Model
	Name       string `gorm:"<-"`
	Dat        []byte `gorm:"<-"`
	UpdateTime int64  `gorm:"<-"`
}

type BehaviorFile struct {
	db *gorm.DB
}

func NewBehaviorFile() *BehaviorFile {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       "gobot:gobot@tcp(db)/gobot?charset=utf8&parseTime=True&loc=Local", // data source name
		DefaultStringSize:         256,                                                               // default size for string fields
		DisableDatetimePrecision:  true,                                                              // disable datetime precision, which not supported before MySQL 5.6
		DontSupportRenameIndex:    true,                                                              // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,                                                              // `change` when rename column, rename column not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false,                                                             // auto configure based on currently MySQL version
	}), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&BehaviorInfo{})

	return &BehaviorFile{
		db: db,
	}
}

func (f *BehaviorFile) Upset(name string, byt []byte) error {

	var res *gorm.DB

	info := BehaviorInfo{
		Name:       name,
		Dat:        byt,
		UpdateTime: time.Now().Unix(),
	}

	_, err := f.Find(name)
	if err == nil {

		res = f.db.Model(&BehaviorInfo{}).Where("name = ?", name).Updates(info)

	} else if err == gorm.ErrRecordNotFound {
		res = f.db.Create(&info)
	}

	return res.Error
}

func (f *BehaviorFile) Del(name string) error {

	info := BehaviorInfo{}
	res := f.db.Where("name = ?", name).Delete(&info)

	return res.Error
}

func (f *BehaviorFile) Find(name string) (BehaviorInfo, error) {
	info := BehaviorInfo{}

	res := f.db.Where("name = ?", name).First(&info)

	return info, res.Error
}

func (f *BehaviorFile) All() ([]BehaviorInfo, error) {
	lst := []BehaviorInfo{}

	result := f.db.Find(&lst)

	return lst, result.Error
}
