package routers

import (
	"github.com/astaxie/beego"
	"gitlab.ecoworkinc.com/Subspace/web-console/controllers"
	filter "gitlab.ecoworkinc.com/Subspace/web-console/routers/filters"
	"fmt"
)

func init() {
	// Static files
	beego.SetStaticPath("/css", "static/css")
	beego.SetStaticPath("/img", "static/img")
	beego.SetStaticPath("/js", "static/js")
	beego.SetStaticPath("/frontend", "static/bower_components/gentelella")
	beego.SetStaticPath("/profile", "public")
	beego.SetStaticPath("/log", "public")

	// Error handling
	beego.ErrorController(&controllers.ErrorController{})

	// Router: Dashboard
	beego.Router("/", &controllers.MainController{}, "get:Index")
	beego.Router("/password_recovery", &controllers.MainController{}, "get,post:ForgotPassword")
	beego.Router("/sign_in", &controllers.MainController{}, "get,post:SignIn")
	beego.Router("/sign_out", &controllers.MainController{}, "get:SignOut")
	beego.Router("/maintenance/?:param", &controllers.MainController{}, "get:Maintenance")

	// Router: Initial Setup
	setupNS := beego.NewNamespace("/init",
		beego.NSRouter("/setup_start", &controllers.InitController{}, "get,post:SetupStart"),
		beego.NSRouter("/create_admin", &controllers.InitController{}, "get,post:CreateAdmin"),
		beego.NSRouter("/setup_complete", &controllers.InitController{}, "get,post:SetupComplete"),
		beego.NSRouter("/clean_install", &controllers.InitController{}, "get,post:CleanInstall"),
		beego.NSRouter("/maintenance", &controllers.InitController{}, "get:Maintenance"),
	)

	// Router: User Management
	userNS := beego.NewNamespace("/users",
		beego.NSRouter("/", &controllers.UserController{}, "get:ListUsers"),
		beego.NSRouter("/:id([0-9]+)", &controllers.UserController{}, "get:UserInfo"),
		beego.NSRouter("/add", &controllers.UserController{}, "get,post:UserCreate"),
		beego.NSRouter("/:id([0-9]+)/edit", &controllers.UserController{}, "get,post:UserUpdate"),
		beego.NSRouter(fmt.Sprintf("/%s/enable", controllers.PATH_USER),
			&controllers.UserController{}, "patch,post:UserEnable"),
		beego.NSRouter(fmt.Sprintf("/%s/disable", controllers.PATH_USER),
			&controllers.UserController{}, "patch,post:UserDisable"),
		beego.NSRouter("/:id([0-9]+)/delete", &controllers.UserController{}, "get,post:UserDelete"),
		beego.NSRouter("/:id([0-9]+)/profile_add", &controllers.UserController{}, "get,post:ProfileCreate"),
		beego.NSRouter("/:id([0-9]+)/profile_info/:profile_id([0-9]+)", &controllers.UserController{}, "get,post:ProfileInfo"),
		beego.NSRouter("/:id([0-9]+)/profile_edit/:profile_id([0-9]+)", &controllers.UserController{}, "get,post:ProfileUpdate"),
		beego.NSRouter("/:id([0-9]+)/profile_delete/:profile_id([0-9]+)", &controllers.UserController{}, "get:ProfileDelete"),
		beego.NSRouter(fmt.Sprintf("/%s/profiles/%s/enable", controllers.PATH_USER, controllers.PATH_PROFILE),
			&controllers.UserController{}, "patch,post:ProfileEnable"),
		beego.NSRouter(fmt.Sprintf("/%s/profiles/%s/disable", controllers.PATH_USER, controllers.PATH_PROFILE),
			&controllers.UserController{}, "patch,post:ProfileDisable"),
		beego.NSRouter(fmt.Sprintf("/%s/profiles/%s/download/%s", controllers.PATH_USER, controllers.PATH_PROFILE, controllers.PATH_PLATFORM),
			&controllers.UserController{}, "get:ProfileDownload"),
		beego.NSRouter("/:id([0-9]+)/session_delete/:session_name", &controllers.UserController{}, "get:SessionDisconnect"),
	)

	logNS := beego.NewNamespace("/logs",
		beego.NSRouter("/", &controllers.LogController{}, "get:ListLogs"),
		beego.NSRouter("/download", &controllers.LogController{}, "get:DownloadLogs"),
	)

	settingNS := beego.NewNamespace("/settings",
		beego.NSRouter("/", &controllers.SettingController{}, "get:SystemInfo"),
		beego.NSRouter("/data_update/uuid_edit", &controllers.SettingController{}, "post:ReGenerateUUID"),
		beego.NSRouter("/data_update/hostname_edit", &controllers.SettingController{}, "post:UpdateHost"),
		beego.NSRouter("/data_update/presharedkey_edit", &controllers.SettingController{}, "post:UpdatePreSharedKey"),
		beego.NSRouter("/mail", &controllers.SettingController{}, "get,post:Mail"),
		beego.NSRouter("/backup_restore", &controllers.SettingController{}, "get:BackupRestore"),
		beego.NSRouter("/ajax_response", &controllers.SettingController{}, "get,post:AjaxResponse"),
		beego.NSRouter("/dns_verify", &controllers.SettingController{}, "post:DNSVerify"),
	)

	// Router: Backup config
	configNS := beego.NewNamespace("/configs",
		beego.NSRouter("/latest", &controllers.SettingController{}, "get:DownloadLatestConfig"),
	)

	aboutNS := beego.NewNamespace("/about",
		beego.NSRouter("/", &controllers.AboutController{}, "get:Get"),
	)

	beego.AddNamespace(setupNS)
	beego.AddNamespace(userNS)
	beego.AddNamespace(logNS)
	beego.AddNamespace(settingNS)
	beego.AddNamespace(configNS)
	beego.AddNamespace(aboutNS)

	// Filters
	beego.InsertFilter("/*", beego.BeforeRouter, filter.FilterMaintenance)
	beego.InsertFilter("/init/*", beego.BeforeRouter, filter.InitialStateDeny)
	beego.InsertFilter("/*", beego.BeforeRouter, filter.InitialStateRequired)
	beego.InsertFilter("/*", beego.BeforeRouter, filter.AuthLogin)
}
