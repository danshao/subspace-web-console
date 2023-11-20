package controllers

import (
	"github.com/astaxie/beego"
)

// ErrorController is beegoController struct
type ErrorController struct {
	beego.Controller
}

// Error404 defines the 404 error page
func (c *ErrorController) Error404() {
	c.Data["Title"], c.TplName, c.Layout, c.LayoutSections = RenderView("Error!", "error", "page_404.html")
}

// Error500 defines the 500 error page
func (c *ErrorController) Error500() {
	c.Data["Title"], c.TplName, c.Layout, c.LayoutSections = RenderView("Error!", "error", "page_500.html")
}

// ErrorMaintenance the Maintenace page
func (c *ErrorController) ErrorMaintenance() {
	c.Data["Title"], c.TplName, c.Layout, c.LayoutSections = RenderView("Error!", "error", "page_maintenance.html")
}
