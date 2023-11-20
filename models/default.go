package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	var (
		host           = beego.AppConfig.String("host")
		databaseType   = "mysql"
		databaseDriver = orm.DRMySQL
		databaseName   = "default"
		databaseSource = "subspace:subspace@tcp(" + host + ":3306)/subspace?charset=utf8"
	)

	orm.Debug = true
	orm.RegisterDriver(databaseType, databaseDriver)
	orm.RegisterDataBase(databaseName, databaseType, databaseSource)
	orm.RegisterModel(new(User), new(Profile), new(ProfileSnapshot), new(Log), new(System))

}
