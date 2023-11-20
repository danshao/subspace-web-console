package controllers

import (
	"errors"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	uuid "github.com/satori/go.uuid"
	"gitlab.ecoworkinc.com/Subspace/softetherlib/softether"
	"gitlab.ecoworkinc.com/Subspace/subspace-utility/subspace/administration/backup"
	"gitlab.ecoworkinc.com/Subspace/subspace-utility/subspace/administration/restore"
	"gitlab.ecoworkinc.com/Subspace/subspace-utility/subspace/repository"
	"gitlab.ecoworkinc.com/Subspace/web-console/helpers"
	"gitlab.ecoworkinc.com/Subspace/web-console/helpers/mail"
	"gitlab.ecoworkinc.com/Subspace/web-console/helpers/system"
	"gitlab.ecoworkinc.com/Subspace/web-console/models"
	"bytes"
	"gitlab.ecoworkinc.com/Subspace/subspace-utility/subspace/administration"
	"fmt"
	"gitlab.ecoworkinc.com/Subspace/web-console/form"
)

type SettingController struct {
	beego.Controller
}

type BackupCallback struct {
	controller *SettingController
}

func (c BackupCallback) OnStart()  {}
func (c BackupCallback) OnCancel() {}
func (c BackupCallback) OnSuccess(yaml string) {
	success := new(bool)
	*success = true
	var operatorID int
	if c.controller.GetSession("userID") != nil {
		operatorID = c.controller.GetSession("userID").(int)
	}
	now := time.Now()
	models.UpdateBackupTime(now)
	models.WriteLog("config_backup", helpers.NewSubspaceRawLog(c.controller.Ctx, success).String(), "", "", c.controller.Ctx.Input.IP(), 0, operatorID)

	backupSnapshotRepo := repository.InitBackupSnapshotRepositoryWithHost(beego.AppConfig.String("host"))
	backupSnapshotRepo.SetLatestSnapshot(yaml)
	backupSnapshotRepo.SetSnapshot(now, yaml)

	// Write to disk
	backupSnapshotRepo.(repository.RedisBackupSnapshotRepository).Client.Save()
}

func (c BackupCallback) OnFail(e error) {}

type RestoreCallback struct {
	controller *SettingController
}

func (c RestoreCallback) OnStart() {
	beego.AppConfig.Set("maintenance_mode", "true")
}
func (c RestoreCallback) OnCancel() {
	beego.AppConfig.Set("maintenance_mode", "false")
}
func (c RestoreCallback) OnSuccess(yaml string) {
	beego.AppConfig.Set("maintenance_mode", "false")
	success := new(bool)
	*success = true
	var operatorID int
	if c.controller.GetSession("userID") != nil {
		operatorID = c.controller.GetSession("userID").(int)
	}
	models.WriteLog("config_restore", helpers.NewSubspaceRawLog(c.controller.Ctx, success).String(), "", "", c.controller.Ctx.Input.IP(), 0, operatorID)

	// Log out all users.
	cmd := exec.Command("redis-cli", "-h", beego.AppConfig.String("host"), "-n", beego.AppConfig.String("session_redis_db_number"), "flushall")
	err := cmd.Run()
	if err != nil {
		beego.Warn("redis flush db failed: ", err)
	}
}
func (c RestoreCallback) OnFail(e error) {
	beego.AppConfig.Set("maintenance_mode", "false")
}

func (c *SettingController) Prepare() {
	// navbar user info
	if c.GetSession("auth") != nil {
		if c.GetSession("auth").(bool) {
			c.Data["Username"] = c.GetSession("username")
			c.Data["UserLink"] = c.GetSession("userlink")
			if c.GetSession("role") != nil {
				c.Data["Role"] = "(" + strings.Title(c.GetSession("role").(string)) + ")"
			} else {
				c.Data["Role"] = "()"
			}
		}
	}
}

// SystemInfo - controller for the system settings page
func (c *SettingController) SystemInfo() {
	currentIp, err := system.RefreshCurrentServerIp()

	flash := beego.ReadFromRequest(&c.Controller)
	_ = flash
	c.Data["Title"], c.TplName, c.Layout, c.LayoutSections = RenderView("Settings", "setting", "system_info.tpl")
	c.LayoutSections["Scripts_Custom"] = "layout/setting/_scripts_custom.tpl"
	systemInfo, err := models.GetSystemInfo()
	systemInfo.Ip = currentIp
	if err != nil {
		beego.Error("Error: ", err)
	}
	c.Data["SystemInfo"] = systemInfo
}

// Mail - controller for the mail settings page
func (c *SettingController) Mail() {
	c.Data["Title"], c.TplName, c.Layout, c.LayoutSections = RenderView("Settings", "setting", "mail.tpl")
	c.LayoutSections["Scripts_Custom"] = "layout/setting/_scripts_custom.tpl"

	var (
		flash = beego.ReadFromRequest(&c.Controller)
		_     = flash
	)

	if strings.ToLower(c.Ctx.Request.Method) == "post" {
		p := form.SmtpSettingsForm{}
		// Parse form
		if err := c.ParseForm(&p); err != nil { // Check if form is parsed correctly
			msg := "Parse form input data error: " + err.Error()
			beego.Debug(msg)
			ShowErrorMessage(c.Controller, msg)
			c.Redirect(c.URLFor("SettingController.Mail"), 302)
			return
		}

		// Validate form
		if _, err := IsFormValid(&p); nil != err {
			msg := "Validate form input data error: " + err.Error()
			beego.Debug(msg)
			ShowErrorMessage(c.Controller, msg)
			c.Redirect(c.URLFor("SettingController.Mail"), 302)
			return
		}

		models.UpdateSMTP(p.SMTPHost, p.SMTPPort, p.Username, p.Password, p.SenderName, p.SenderEmail, p.Authentication)
		if !p.IsEmpty() {
			if success := mail.TestSMTP(models.GetUserWithID(c.GetSession("userID").(int)).Email); success {
				beego.Debug("SMTP Server setup successful.")
				ShowNoticeMessage(c.Controller, "SMTP server settings look good! Check your email for a confirmation.")
				c.Redirect(c.URLFor("SettingController.Mail"), 302)
			} else {
				beego.Debug("SMTP Server setup unsuccessful.")
				ShowErrorMessage(c.Controller, "There appears to be something wrong in your SMTP server settings. Check your parameters or your SMTP server and try again.")
				c.Redirect(c.URLFor("SettingController.Mail"), 302)
			}
		}
	}

	systemInfo, _ := models.GetSystemInfo()
	c.Data["SystemInfo"] = systemInfo
}

func (c *SettingController) BackupRestore() {
	c.Data["Title"], c.TplName, c.Layout, c.LayoutSections = RenderView("Settings", "setting", "backup_restore.tpl")
	c.LayoutSections["Scripts_Custom"] = "layout/setting/_scripts_custom.tpl"
	systemInfo, err := models.GetSystemInfo()
	if err != nil {
		beego.Error("Error: ", err)
	}
	c.Data["SystemInfo"] = systemInfo
}

func (c *SettingController) AjaxResponse() {

	beego.Debug("receiving " + c.Ctx.Request.Method + " event from ajax")
	var (
		host      = beego.AppConfig.String("host")
		dbURI     = "subspace:subspace@tcp(" + host + ":3306)/subspace?charset=utf8&parseTime=True&loc=Local"
		downloadBackupFileUri = "/configs/latest"
	)
	// 1. start backup / restore (POST)
	if strings.ToLower(c.Ctx.Request.Method) == "post" {
		option := c.GetString("backup_restore_option")
		beego.Debug("do: ", c.GetString("backup_restore_option"))

		if option == "backup" {
			backupController := backup.GetInstance()
			backupController.SetDatabaseUri(dbURI)

			backupController.SetCallback(BackupCallback{c})
			backupController.Start()
			c.Data["json"] = downloadBackupFileUri
		} else if option == "restore" {
			_, fileHeader, err := c.GetFile("upload_config")
			start := false
			if err != nil {
				beego.Debug("Error: ", err)
			} else {
				beego.Debug("File name: ", fileHeader.Filename)
				beego.Debug("File header: ", fileHeader.Header)

				filePath := "./public/restore-" + fileHeader.Filename
				c.SaveToFile("upload_config", filePath)

				restoreController := restore.GetInstance()
				restoreController.SetDatabaseUri(dbURI)

				restoreController.SetCallback(RestoreCallback{c})
				start = restoreController.Start(filePath)

				// TODO: put service into maintenace mode.
				// beego.AppConfig.Set("maintenance_mode", "true")
			}
			if start {
				c.Data["json"] = fileHeader.Filename
			} else {
				c.Data["json"] = "false"
			}
		} else {
			c.Data["json"] = "ok"
		}
		c.ServeJSON()
		return
	}

	// 2. ask status (GET)
	if strings.ToLower(c.Ctx.Request.Method) == "get" {
		var option string
		c.Ctx.Input.Bind(&option, "backup_restore_option")
		beego.Debug("ask status: ", option)

		if option == "backup" {
			backupController := backup.GetInstance()
			status := backupController.GetStatus()
			switch status.Step {
			case backup.IDLE:
				beego.Info("backup idle")
				c.Data["json"] = "backup idle"
			case backup.RUNNING:
				beego.Info("backup running")
				c.Data["json"] = "backup running"
			case backup.SUCCEED:
				beego.Info("backup succeeded")
				c.Data["json"] = "backup succeeded" // + status.Result
			case backup.FAILED:
				beego.Error("backup failed", status.Error)
				c.Data["json"] = "backup failed: " + status.Error.Error()
			case backup.CANCELING:
				beego.Warn("backup canceling")
				c.Data["json"] = "backup canceling"
			case backup.CANCELED:
				beego.Warn("backup canceled")
				c.Data["json"] = "backup canceled"
			case backup.UNKNOWN:
				beego.Error("something wrong")
				c.Data["json"] = "backup wrong"
			}
		} else if option == "restore" {
			c.Data["json"] = "restore ok"
			restoreController := restore.GetInstance()
			status := restoreController.GetStatus()
			switch status.Step {
			case restore.IDLE:
				beego.Info("restore idle")
				c.Data["json"] = "restore idle"
			case restore.RUNNING:
				beego.Info("restore running")
				c.Data["json"] = "restore running"
			case restore.SUCCEED:
				beego.Info("restore succeeded")
				c.Data["json"] = "restore succeeded" // + status.Result
			case restore.FAILED:
				beego.Error("restore failed", status.Error)
				c.Data["json"] = "restore failed: " + status.Error.Error()
			case restore.CANCELING:
				beego.Warn("restore canceling")
				c.Data["json"] = "restore canceling"
			case restore.CANCELED:
				beego.Warn("restore canceled")
				c.Data["json"] = "restore canceled"
			case restore.UNKNOWN:
				beego.Error("something wrong")
				c.Data["json"] = "restore wrong"
			}
		} else {
			c.Data["json"] = "ok"
		}
		c.ServeJSON()
	}
}

func (c *SettingController) ReGenerateUUID() {
	var (
		err     error
		result  string
		warnMsg string
		success = new(bool)
	)

	*success = false

	uuid := uuid.NewV4()
	_, err = models.UpdateSystemUUID(uuid.String())
	if err != nil {
		result = "UUID was not updated. Error: " + err.Error()
		beego.Error("UUID was not updated. Error: ", err)
	} else {
		*success = true
		result = "UUID was successfully updated."
	}

	models.WriteLog("uuid_update", helpers.NewSubspaceRawLog(c.Ctx, success).String(), "", "", c.Ctx.Input.IP(), -1, c.GetSession("userID").(int))
	flash := beego.NewFlash()
	if err != nil {
		flash.Error(result)
	} else {
		flash.Notice(result)
	}
	if warnMsg != "" {
		flash.Warning(warnMsg)
	}
	flash.Store(&c.Controller)
	c.Redirect(c.URLFor("SettingController.SystemInfo"), 302)
}

func (c *SettingController) UpdateHost() {
	var (
		err     error
		result  string
		warnMsg string
		success = new(bool)
	)

	domainName := strings.TrimSpace(c.GetString("hostname"))

	*success = false
	re := regexp.MustCompile("^((([a-zA-Z]{1})|([a-zA-Z]{1}[a-zA-Z]{1})|([a-zA-Z]{1}[0-9]{1})|([0-9]{1}[a-zA-Z]{1})|([a-zA-Z0-9][a-zA-Z0-9-_]{1,61}[a-zA-Z0-9])).)+([a-zA-Z]{2,61})$")
	// beego.Debug("Match host name result: ", re.MatchString(s.Hostname))
	if re.MatchString(domainName) || domainName == "" {
		_, err = models.UpdateSystemHost(domainName)
		if err != nil {
			result = "Host name was not updated. Error: " + err.Error()
			beego.Error("Host name was not updated. Error: ", err)
		} else {
			*success = true
			result = "Host name was successfully updated."
			if helpers.DNSMatch(domainName) == false {
				warnMsg = "Hostname does not resolve to the IP of Subspace."
			}
			mail.SendSystemSettingsChangeEmail("Hostname", domainName)
		}
	} else {
		err = errors.New("hostname format is incorrect")
		result = "Hostname was not updated. Error: " + err.Error()
		beego.Error("Hostname was not updated. Error: ", err)
	}

	models.WriteLog("hostname_update", helpers.NewSubspaceRawLog(c.Ctx, success).String(), "", "", c.Ctx.Input.IP(), -1, c.GetSession("userID").(int))
	flash := beego.NewFlash()
	if err != nil {
		flash.Error(result)
	} else {
		flash.Notice(result)
	}
	if warnMsg != "" {
		flash.Warning(warnMsg)
	}
	flash.Store(&c.Controller)
	c.Redirect(c.URLFor("SettingController.SystemInfo"), 302)
}

func (c *SettingController) UpdatePreSharedKey() {
	var (
		err     error
		result  string
		warnMsg string
		logType string

		success = new(bool)

		se = softether.SoftEther{IP: VPN_HOST, Password: SOFTETHER_MANGEMENT_PASSWORD, Hub: SOFTETHER_HUB}
	)

	preSharedKey := strings.TrimSpace(c.GetString("presharedkey"))

	*success = false
	logType = "preshared_key_update"
	re := regexp.MustCompile("^[a-zA-Z0-9]{1,9}$")
	if re.MatchString(preSharedKey) {
		if retCode := se.SetPreSharedKey(preSharedKey); retCode == 0 {
			if _, err = models.UpdatePreSharedKey(preSharedKey); err == nil {
				disconnectAllSessions(se, c)
				*success = true
				result = "Pre-shared key was successfully updated."
			} else {
				result = "Pre-shared key was not updated. Error: " + err.Error()
				beego.Error("Pre-shared key was not updated. Error: ", err)
			}
		} else {
			result = "Pre-shared key was not updated. Error: " + err.Error()
			beego.Error("Pre-shared key was not updated. Error: ", err)
		}
	} else {
		err = errors.New("Pre-shared key must be alphanumeric and 9 characters or less.")
		result = "Pre-shared key was not updated. Error: " + err.Error()
		beego.Error("Pre-shared key was not updated. Error: ", err)
	}

	models.WriteLog(logType, helpers.NewSubspaceRawLog(c.Ctx, success).String(), "", "", c.Ctx.Input.IP(), -1, c.GetSession("userID").(int))
	flash := beego.NewFlash()
	if err != nil {
		flash.Error(result)
	} else {
		flash.Notice(result)
	}
	if warnMsg != "" {
		flash.Warning(warnMsg)
	}
	flash.Store(&c.Controller)
	c.Redirect(c.URLFor("SettingController.SystemInfo"), 302)
}

func (c *SettingController) DownloadLatestConfig() {
	// Check user who access this url is authorized.
	operatorId := c.GetSession("userID").(int)
	operator := models.GetUserWithID(operatorId)
	if models.USER_ROLE_ADMIN != operator.Role {
		// Return error, not admin
		beego.Debug("Operator", operatorId, "is not admin and try to download latest backup file.")
		c.CustomAbort(401, "Only admin can download the backup file.")
	}

	// Fetch backup from Redis
	backupSnapshotRepo := repository.InitBackupSnapshotRepositoryWithHost(beego.AppConfig.String("host"))
	if content, err := backupSnapshotRepo.GetLatestSnapshot(); nil == err {
		cfg, _ := administration.ParseConfig([]byte(content))
		timeString := cfg.GetConfigCreateTime().Format("20060102-150405")
		fileName := fmt.Sprintf("backup-%s.config", timeString)

		// Response to client
		buf := &bytes.Buffer{}
		buf.WriteString(content)

		c.Ctx.Output.Header("Content-Type", "application/force-download")
		c.Ctx.Output.Header("Content-Disposition", "attachment;filename="+fileName)
		c.Ctx.Output.Header("Expires", "0")
		c.Ctx.Output.Header("Cache-Control", "no-cache, no-store, must-revalidate")

		c.Ctx.Output.Body(buf.Bytes())
	} else {
		beego.Error("Get latest backup config.", err.Error())
		c.CustomAbort(500, "Get latest backup config.")
	}
}

func (c *SettingController) DNSVerify() {
	hostname := c.GetString("hostname")

	if hostname != "" {
		c.Data["json"] = helpers.DNSMatch(hostname)
	} else {
		c.Data["json"] = false
	}
	c.ServeJSON()
}

func disconnectAllSessions(se softether.SoftEther, c *SettingController) (err error) {
	// 1a. Get all sessions from Redis
	var (
		serverIP       = beego.AppConfig.String("host")
		success        = new(bool)
		userID, _      = strconv.Atoi(c.Ctx.Input.Param(":id"))
	)
	sessionRepository := repository.InitSessionRepositoryWithHost(serverIP)
	sessions, err := sessionRepository.GetAllSessions()
	if err != nil {
		beego.Error("Error occurred while getting session details by user ID: ", err)
		return
	}

	*success = true
	// 1b. Loop through all sessions and disconnect
	for _, session := range sessions {
		if returnCode := se.DisconnectSession(session.SessionName); returnCode != 0 {
			*success = false
			err = errors.New("Error disconnecting session " + session.SessionName + ". Error: " + softether.Strerror(returnCode))
		}
		models.WriteLog("session_kicked", helpers.NewSubspaceRawLog(c.Ctx, success).String(), "", session.SessionName, c.Ctx.Input.IP(), userID, c.GetSession("userID").(int))
	}

	return
}
