package main

import (
	"github.com/astaxie/beego"

	_ "github.com/astaxie/beego/session/redis"
	"gitlab.ecoworkinc.com/Subspace/web-console/helpers"
	_ "gitlab.ecoworkinc.com/Subspace/web-console/models"
	_ "gitlab.ecoworkinc.com/Subspace/web-console/routers"
)

func init() {
	beego.AddFuncMap("localTimeFmt", helpers.LocalTimeFmt)
	if beego.AppConfig.String("RunMode") == "dev" {
		beego.SetLevel(beego.LevelDebug)
		// beego.SetLogger("file", `{"filename":"public/dev.log"}`)
	}
	if beego.BConfig.WebConfig.Session.SessionProvider == "redis" {
		beego.BConfig.WebConfig.Session.SessionProviderConfig = beego.AppConfig.String("host") + ":6379,100,," + beego.AppConfig.String("session_redis_db_number")
	}
}

func main() {
	beego.Run()
}
