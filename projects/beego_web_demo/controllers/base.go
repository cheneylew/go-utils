package controllers

import (
	"github.com/cheneylew/goutil/utils"
	beego2 "github.com/cheneylew/goutil/utils/beego"
)

var BEEGO_CONFIG beego2.BeegoConfig
var FILTER_PATHS  []string

func init() {
	// configs
	FILTER_PATHS = append(FILTER_PATHS,"/user/login")
	FILTER_PATHS = append(FILTER_PATHS,"/user/regist")
	BEEGO_CONFIG = beego2.BeegoConfig{
		LoginCheck:false,
	}

}

type BaseController struct {
	beego2.BBaseController
}

func (c *BaseController) Prepare() {
	c.BBaseController.Prepare()
	urlPath := c.Ctx.Request.URL.Path

	c.Data["Website"] = "site name"
	c.Data["Email"] = "site email"

	c.Layout = "layout.html"

	if BEEGO_CONFIG.LoginCheck {
		if !utils.Contain(FILTER_PATHS, urlPath) {
			if c.IsLogin() {
				c.Data["User"] = c.GetUser()
			} else {
				c.RedirectWithURL("/user/login")
			}
		}
	}
}


