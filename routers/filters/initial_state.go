package filter

import (
	"gitlab.ecoworkinc.com/Subspace/web-console/models"

	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

// InitialState control access to initial setup flow
var InitialStateRequired = func(ctx *context.Context) {
	beego.Debug("URL: ", ctx.Input.URL())
	if strings.HasPrefix(ctx.Input.URL(), "/init/") {
		if _, ok := ctx.Input.Session("init_process").(bool); !ok && !strings.Contains(ctx.Input.URL(), "setup_start") {
			beego.Debug("Incorrect initial setup sequence, redirecting to Onboarding Flow")
			ctx.Redirect(302, "/init/setup_start")
		}
		return
	}
	if q := models.GetAllUsers(1); len(q) == 0 {
		beego.Debug("No users found in subspace_users, redirecting to Onboarding Flow")
		ctx.Redirect(302, "/init/setup_start")
	}
}

var InitialStateDeny = func(ctx *context.Context) {
	if _, ok := ctx.Input.Session("init_process").(bool); ok {
		return
	}
	if q := models.GetAllUsers(1); len(q) > 0 {
		beego.Debug("Already done onboarding, redirect to dashboard.")
		ctx.Redirect(302, "/")
	}
}
