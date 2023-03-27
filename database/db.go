package database

import (
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DefaultConfig = map[string]string{
		"global": `{"title":"Global","content":"\n--[[\n\tGlobal constant area, users can define some constants here; it is easy to call in other scripts\n]]--\n\nREMOTE = \"http://127.0.0.1:8888\"\n","closable":false, "prefab":false}`,
		"http":   `{"title":"HTTP","content":"\nlocal parm = {\n    body = {},    -- request body\n    timeout = \"10s\",\n    headers = {},\n}\n\nlocal url = REMOTE .. \"/group/methon\"\nlocal http = require(\"http\")\n\nfunction execute()\n    res, errmsg = http.post(url, parm)\n  \tif errmsg ~= nil then\n\t\tmeta.Err = errmsg\n    \treturn\n  \tend\n  \t\n  \tif res[\"status_code\"] ~= 200 then\n\t\tmeta.Err = \"post \" .. url .. \" http status code err \" .. res[\"status_code\"]\n  \t\treturn\n  \tend\n  \n  \tbody = json.decode(res[\"body\"])\n  \tmerge(meta, body.Body)\n\nend\n","closable":false, "prefab":true}`,
		"system": `{"channelsize":1024, "reportsize":100}`,
	}
)

type Cache struct {
	conf     *Conf
	behavior *Behavior
	prefab   *Prefab
	report   *Report

	mysqlptr *gorm.DB
}

var db *Cache
var once sync.Once

func GetConfig() *Conf {
	return db.conf
}

func GetPrefab() *Prefab {
	return db.prefab
}

func GetReport() *Report {
	return db.report
}

func GetBehavior() *Behavior {
	return db.behavior
}

func Init(NoDBMode bool) *Cache {
	once.Do(func() {

		pwd := os.Getenv("MYSQL_PASSWORD")
		if pwd == "" {
			panic(errors.New("mysql password is not defined"))
		}

		name := os.Getenv("MYSQL_DATABASE")
		if name == "" {
			panic(errors.New("mysql database is not defined"))
		}

		host := os.Getenv("MYSQL_HOST")
		if host == "" {
			panic(errors.New("mysql host is not defined"))
		}

		user := os.Getenv("MYSQL_USER")
		if user == "" {
			panic(errors.New("mysql user is not defined"))
		}

		dsn := user + ":" + pwd + "@tcp(" + host + ")/" + name + "?charset=utf8&parseTime=True&loc=Local"

		var sqlptr *gorm.DB
		var err error

		if NoDBMode {
			sqlptr, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		} else {
			sqlptr, err = gorm.Open(mysql.New(mysql.Config{
				DSN:                       dsn,   // data source name
				DefaultStringSize:         256,   // default size for string fields
				DisableDatetimePrecision:  true,  // disable datetime precision, which not supported before MySQL 5.6
				DontSupportRenameIndex:    true,  // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
				DontSupportRenameColumn:   true,  // `change` when rename column, rename column not supported before MySQL 8, MariaDB
				SkipInitializeWithVersion: false, // auto configure based on currently MySQL version
			}), &gorm.Config{})
		}

		if err != nil {
			fmt.Println("open mysql err", err.Error())
			return
		}

		db = &Cache{
			mysqlptr: sqlptr,
			conf:     CreateConfig(sqlptr),
			prefab:   CreatePrefab(sqlptr),
			behavior: CreateBehavior(sqlptr),
			report:   CreateReport(sqlptr),
		}
	})

	return db
}

func init() {
	db = &Cache{}
}
