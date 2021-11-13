package database

import (
	"errors"
	"os"
	"sync"
	"time"

	"github.com/pojol/gobot/bot"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type BehaviorInfo struct {
	gorm.Model
	Name       string `gorm:"<-"`
	Dat        []byte `gorm:"<-"`
	UpdateTime int64  `gorm:"<-"`
	Status     string `gorm:"<-"`
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

type Database struct {
	db *gorm.DB
	sync.Mutex
}

var bf *Database

func New(pwd, name, host, user string) {

	if pwd == "" {
		pwd = os.Getenv("MYSQL_PASSWORD")
		if pwd == "" {
			panic(errors.New("mysql password is not defined"))
		}
	}

	if name == "" {
		name = os.Getenv("MYSQL_DATABASE")
		if name == "" {
			panic(errors.New("mysql database is not defined"))
		}
	}

	if host == "" {
		host = os.Getenv("MYSQL_HOST")
		if host == "" {
			panic(errors.New("mysql host is not defined"))
		}
	}

	if user == "" {
		user = os.Getenv("MYSQL_USER")
		if user == "" {
			panic(errors.New("mysql user is not defined"))
		}
	}

	dsn := user + ":" + pwd + "@tcp(" + host + ")/" + name + "?charset=utf8&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,   // data source name
		DefaultStringSize:         256,   // default size for string fields
		DisableDatetimePrecision:  true,  // disable datetime precision, which not supported before MySQL 5.6
		DontSupportRenameIndex:    true,  // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,  // `change` when rename column, rename column not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false, // auto configure based on currently MySQL version
	}), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&BehaviorInfo{}, &BotTemplateConfig{}, &BotConfig{})
	if err != nil {
		panic(err)
	}

	bf = &Database{
		db: db,
	}
}

func Get() *Database {
	return bf
}

func (f *Database) UpsetFile(name string, byt []byte) error {

	f.Lock()
	defer f.Unlock()

	var res *gorm.DB

	info := BehaviorInfo{
		Name:       name,
		Dat:        byt,
		Status:     bot.BotStatusUnknow,
		UpdateTime: time.Now().Unix(),
	}

	_, err := f.FindFile(name)
	if err == nil {

		res = f.db.Model(&BehaviorInfo{}).Where("name = ?", name).Updates(info)

	} else if err == gorm.ErrRecordNotFound {
		res = f.db.Create(&info)
	}

	return res.Error
}

func (f *Database) UpdateState(name string, status string) error {
	return f.db.Model(&BehaviorInfo{}).Where("name = ?", name).Update("Status", status).Error
}

func (f *Database) DelFile(name string) error {

	f.Lock()
	defer f.Unlock()

	info := BehaviorInfo{}
	res := f.db.Where("name = ?", name).Delete(&info)

	return res.Error
}

func (f *Database) FindFile(name string) (BehaviorInfo, error) {

	info := BehaviorInfo{}

	res := f.db.Where("name = ?", name).First(&info)

	return info, res.Error
}

func (f *Database) GetAllFiles() ([]BehaviorInfo, error) {

	lst := []BehaviorInfo{}

	result := f.db.Find(&lst)

	return lst, result.Error
}

func init() {
	New("", "", "", "")
}
