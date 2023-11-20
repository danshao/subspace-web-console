package filter

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

var FilterMaintenance = func(ctx *context.Context) {
	if beego.AppConfig.String("maintenance_mode") == "true" {
		ctx.Abort(302, "Maintenance")
	}
}
