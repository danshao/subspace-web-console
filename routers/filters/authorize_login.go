package filter

import (
	"strings"

	"gitlab.ecoworkinc.com/Subspace/web-console/definition"
	"gitlab.ecoworkinc.com/Subspace/web-console/helpers"
	"gitlab.ecoworkinc.com/Subspace/web-console/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

var PATHS_CAN_ACCESS_BY_ANONYMOUS = []string{
	"/init/",
	"/sign_in",
	"/settings/dns_verify",
	"/password_recovery",
}

var AuthLogin = func(ctx *context.Context) {
	// If path is allow anonymous access
	if CanAccessByAnonymous(ctx.Input.URL()) {
		return
	}

	success := new(bool)

	// Check is enabled user
	userID, ok := ctx.Input.Session(definition.SESSION_USER_ID).(int)
	if !ok {
		models.WriteLog("user_session_timeout", helpers.NewSubspaceRawLog(ctx, success).String(), "", "", ctx.Input.IP(), 0, 0)
		beego.Debug("No session('userID') found, redirecting to Sign In")
		RedirectToSignIn(ctx)
		return
	}

	// TODO Change GetUserWithID return nil, error
	currentUser := models.GetUserWithID(userID)
	if !currentUser.Enabled {
		models.WriteLog("user_session_timeout", helpers.NewSubspaceRawLog(ctx, success).String(), "", "", ctx.Input.IP(), currentUser.Id, 0)
		beego.Debug("User is disabled, redirecting to Sign In")
		RedirectToSignIn(ctx)
		return
	}

	//TODO Check ACL using casbin with Role, Path+Regex, HTTP method. e.g. "admin", "/users/[0-9]+/profiles", "POST"
	//TODO see https://github.com/casbin/casbin
	// All path except those allow anonymous access, only allow admin access.
	if models.USER_ROLE_ADMIN != currentUser.Role {
		models.WriteLog("user_session_timeout", helpers.NewSubspaceRawLog(ctx, success).String(), "", "", ctx.Input.IP(), currentUser.Id, 0)
		beego.Debug("User is not admin, redirecting to Sign In")
		RedirectToSignIn(ctx)
		return
	}
}

// Resources
// https://github.com/yinwhm12/catw/blob/2fb25761f48a2aa1cc1fdb1c4a06a793fa1b5163/catw/filters/log_auth.go

func CanAccessByAnonymous(path string) bool {
	for _, allowedPath := range PATHS_CAN_ACCESS_BY_ANONYMOUS {
		if strings.HasPrefix(path, allowedPath) {
			return true
		}
	}
	return false
}

func RedirectToSignIn(ctx *context.Context) {
	ctx.Output.Session("origin-uri", ctx.Input.URL())
	ctx.Redirect(302, "/sign_in")
}
