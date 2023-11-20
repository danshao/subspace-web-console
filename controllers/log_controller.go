package controllers

import (
	"bytes"
	"encoding/csv"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"gitlab.ecoworkinc.com/Subspace/web-console/models"
)

type LogController struct {
	beego.Controller
}

type LogDetail struct {
	Id      int
	LogTime time.Time
	Type    string
	RawLog  string
}

var logTypeMap = map[int]string{
	-1:  "undefined",
	0:   "user_sign_in",
	1:   "user_sign_out",
	2:   "user_session_timeout",
	10:  "user_create",
	11:  "user_edit",
	12:  "user_disable",
	13:  "user_enable",
	14:  "user_delete",
	20:  "profile_create",
	21:  "profile_edit",
	22:  "profile_disable",
	23:  "profile_enable",
	24:  "profile_delete",
	25:  "profile_download_apple",
	26:  "profile_download_windows",
	30:  "sessions_connect",
	31:  "session_disconnect",
	32:  "session_kicked",
	33:  "session_auth_fail",
	40:  "hostname_update",
	41:  "uuid_update",
	42:  "preshared_key_update",
	100: "subspace_start",
	101: "subspace_stop",
	110: "vpn_server_start",
	111: "vpn_server_stop",
	120: "config_backup",
	121: "config_restore",
}

// Prepare
func (c *LogController) Prepare() {
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

// ListLogs
func (c *LogController) ListLogs() {
	c.Data["Title"], c.TplName, c.Layout, c.LayoutSections = RenderView("System Log", "log", "index.tpl")
	c.LayoutSections["Scripts_Custom"] = "layout/log/_scripts_custom.tpl"

	var (
		logList = models.ScanLog()
		logs    []LogDetail
	)

	for _, log := range logList {
		var l LogDetail
		l.Id = log.Id
		l.LogTime = log.LogTime
		l.Type = logTypeMap[log.Type]
		l.RawLog = log.RawLog
		logs = append(logs, l)
	}

	c.Data["LogList"] = logs
}

// DownloadLogs
func (c *LogController) DownloadLogs() {
	var (
		filename = "subspace_logs-" + time.Now().Format("20060102-150405") + ".csv"
		logList  = models.ScanLog()
	)

	// create io buffer and csv writer
	buf := &bytes.Buffer{}
	writer := csv.NewWriter(buf)

	// write data and send to buffer
	writer.Write([]string{"id", "time", "type", "message"})
	for i, log := range logList {
		err := writer.Write([]string{strconv.Itoa(i), log.LogTime.String(), logTypeMap[log.Type], log.RawLog})
		if err != nil {
			beego.Debug("Cannot write to file", err)
		}
	}
	writer.Flush()

	c.Ctx.Output.Header("Content-Type", "text/csv")
	c.Ctx.Output.Header("Content-Disposition", "attachment;filename="+filename)
	c.Ctx.Output.Header("Expires", "0")
	c.Ctx.Output.Header("Cache-Control", "no-cache, no-store, must-revalidate")

	c.Ctx.Output.Body(buf.Bytes())
}
