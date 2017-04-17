package services

import (
	"fmt"
	"net/url"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/krufyliu/dkvgo/dkvgo-admin/models"
)

type Page struct {
	Total    int64 `json:"total"`
	Current  int   `json:"current"`
	PageSize int   `josn:"pageSize"`
}

var (
	o           orm.Ormer
	UserService *userService
	JobService  *jobService
)

func ormSetup() {
	dbHost := beego.AppConfig.String("db.host")
	dbType := beego.AppConfig.String("db.type")
	dbName := beego.AppConfig.String("db.name")
	dbPort := beego.AppConfig.String("db.port")
	dbUser := beego.AppConfig.String("db.user")
	dbPassword := beego.AppConfig.String("db.password")
	dbTimezone := beego.AppConfig.String("db.timezone")
	dbCharset := beego.AppConfig.String("db.charset")
	if dbHost == "" {
		dbHost = "localhost"
	}
	if dbPort == "" {
		dbPort = "3306"
	}
	if dbCharset == "" {
		dbCharset = "utf8"
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s",
		dbUser, dbPassword, dbHost, dbPort, dbName, dbCharset)
	if dbTimezone != "" {
		dsn = dsn + "&loc=" + url.QueryEscape(dbTimezone)
	} else {
		dsn = dsn + "&loc=Local"
	}
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", dbType, dsn)
	orm.RegisterModel(
		new(models.User),
		new(models.Job),
		new(models.JobState),
	)
	if beego.AppConfig.String("runmode") != "prod" {
		orm.Debug = true
	}

	o = orm.NewOrm()
}

func GetOrm() orm.Ormer {
	return o
}

func initServices() {
	UserService = &userService{}
	JobService = &jobService{}
}

func Init() {
	ormSetup()
	initServices()
}