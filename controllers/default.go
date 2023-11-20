package controllers

import (
	"gitlab.ecoworkinc.com/Subspace/web-console/form"
	"github.com/astaxie/beego"
)

const FORM_DATA = "Form"

var (
	// DefaultTitle is the default prefix for the title tag of the html
	DefaultTitle = "Subspace | "

	// LayoutPath defines the locations of layout files for routing paths
	LayoutPath = map[string]map[string]string{
		"error": map[string]string{
			"Layout":  "layout/error/default.tpl",
			"Header":  "layout/error/_header.tpl",
			"Footer":  "layout/error/_footer.tpl",
			"Scripts": "layout/error/_scripts.tpl",
		},
		"main": map[string]string{
			"Layout":  "layout/main/default.tpl",
			"Header":  "layout/main/_header.tpl",
			"Nav":     "layout/main/_nav.tpl",
			"Footer":  "layout/main/_footer.tpl",
			"Sidebar": "layout/main/_sidebar.tpl",
			"Scripts": "layout/main/_scripts.tpl",
		},
		"init": map[string]string{
			"Layout":  "layout/init/default.tpl",
			"Header":  "layout/init/_header.tpl",
			"Footer":  "layout/init/_footer.tpl",
			"Scripts": "layout/init/_scripts.tpl",
		},
		"user": map[string]string{
			"Layout":  "layout/user/default.tpl",
			"Header":  "layout/user/_header.tpl",
			"Nav":     "layout/main/_nav.tpl",
			"Footer":  "layout/main/_footer.tpl",
			"Sidebar": "layout/main/_sidebar.tpl",
			"Scripts": "layout/user/_scripts.tpl",
		},
		"log": map[string]string{
			"Layout":  "layout/log/default.tpl",
			"Header":  "layout/log/_header.tpl",
			"Nav":     "layout/main/_nav.tpl",
			"Footer":  "layout/main/_footer.tpl",
			"Sidebar": "layout/main/_sidebar.tpl",
			"Scripts": "layout/log/_scripts.tpl",
		},
		"setting": map[string]string{
			"Layout":  "layout/setting/default.tpl",
			"Header":  "layout/setting/_header.tpl",
			"Nav":     "layout/main/_nav.tpl",
			"Footer":  "layout/main/_footer.tpl",
			"Sidebar": "layout/main/_sidebar.tpl",
			"Scripts": "layout/setting/_scripts.tpl",
		},
		"about": map[string]string{
			"Layout":  "layout/about/default.tpl",
			"Header":  "layout/about/_header.tpl",
			"Nav":     "layout/main/_nav.tpl",
			"Footer":  "layout/main/_footer.tpl",
			"Sidebar": "layout/main/_sidebar.tpl",
			"Scripts": "layout/about/_scripts.tpl",
		},
	}

	// ContentPath defines the content locations of layout files for routing paths
	ContentPath = map[string]string{
		"error":   "content/error/",
		"main":    "content/main/",
		"init":    "content/init/",
		"user":    "content/user/",
		"log":     "content/log/",
		"about":   "content/about/",
		"setting": "content/setting/",
	}
)

// RenderView returns the exact paths based on the parameters given
func RenderView(title, layout, content string) (string, string, string, map[string]string) {
	// this.Data["Title"] = DefaultTitle + "Error!"
	// layout := "error"
	// this.TplName = ContentPath[layout] + "page_404.html"
	// this.Layout = LayoutPath[layout]["Layout"]
	// this.LayoutSections = make(map[string]string)
	// this.LayoutSections = LayoutPath[layout]
	t := DefaultTitle + title
	tp := ContentPath[layout] + content
	l := LayoutPath[layout]["Layout"]
	s := make(map[string]string)
	s = LayoutPath[layout]
	return t, tp, l, s
}

func IsFormValid(f form.IForm) (bool, error) {
	return f.IsValid()
}

func ShowSuccessMessage(controller beego.Controller, msg string) {
	flash := beego.NewFlash()
	flash.Success(msg)
	flash.Store(&controller)
}

func ShowNoticeMessage(controller beego.Controller, msg string) {
	flash := beego.NewFlash()
	flash.Notice(msg)
	flash.Store(&controller)
}

func ShowWarningMessage(controller beego.Controller, msg string) {
	flash := beego.NewFlash()
	flash.Warning(msg)
	flash.Store(&controller)
}

func ShowErrorMessage(controller beego.Controller, msg string) {
	flash := beego.NewFlash()
	flash.Error(msg)
	flash.Store(&controller)
}