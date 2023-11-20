package controllers

import (
	"strings"

	"github.com/astaxie/beego"
	"gitlab.ecoworkinc.com/Subspace/web-console/models"
)

type AboutController struct {
	beego.Controller
}

func (c *AboutController) Prepare() {
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

func (c *AboutController) Get() {
	c.Data["Title"], c.TplName, c.Layout, c.LayoutSections = RenderView("About", "about", "index.tpl")
	systemInfo, _ := models.GetSystemInfo()

	c.Data["VersionNumber"] = systemInfo.SubspaceVersion
	c.Data["BuildNumber"] = systemInfo.SubspaceBuildNumber

}
