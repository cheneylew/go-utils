package controllers

import (
	"github.com/astaxie/beego"
	"github.com/cheneylew/goutil/utils"
	"github.com/cheneylew/goutil/iOS_tool_server/models"
	"strings"
	"net/url"
)

type BeegoConfig struct {
	LoginCheck bool
}

var SESSTION_KEY_USER string
var BEEGO_CONFIG BeegoConfig
var FILTER_PATHS  []string

func init() {
	// configs
	SESSTION_KEY_USER = "LOGINED_USER"
	FILTER_PATHS = append(FILTER_PATHS,"/user/login")
	FILTER_PATHS = append(FILTER_PATHS,"/user/regist")
	BEEGO_CONFIG = BeegoConfig{
		LoginCheck:false,
	}
}

type BaseController struct {
	beego.Controller
}

func (c *BaseController) Prepare() {
	c.Controller.Prepare()
	urlPath := c.Ctx.Request.URL.Path

	c.Data["Website"] = "site name"
	c.Data["Email"] = "site email"

	c.Layout = "layout.html"

	if BEEGO_CONFIG.LoginCheck {
		if !utils.Contain(FILTER_PATHS, urlPath) {
			if c.IsLogin() {
				user := c.GetLoginedUser()
				c.Data["User"] = user
			} else {
				c.RedirectWithURL("/user/login")
			}
		}
	}
}

func (c *BaseController) IsPost() bool {
	return c.Ctx.Request.Method == "POST"
}

func (c *BaseController) IsGet() bool {
	return c.Ctx.Request.Method == "GET"
}

func (c *BaseController) PostForm() url.Values {
	return c.Ctx.Request.PostForm
}

func (c *BaseController) PostFormWithKey(key string) []string {
	return c.Ctx.Request.PostForm[key]
}

func (c *BaseController) Finish() {
	c.Controller.Finish()
}

//通用
func (c *BaseController) RedirectWithURL(url string) {
	c.Redirect(url, 302)
}

//用户
func (c *BaseController) IsLogin() bool {
	return c.GetSession(SESSTION_KEY_USER) != nil
}

func (c *BaseController) SetLoginedUser(user models.User) {
	c.SetSession(SESSTION_KEY_USER,user)
}

func (c *BaseController) SetUserLogout() {
	c.SetSession(SESSTION_KEY_USER,nil)
}

func (c *BaseController) GetLoginedUser() *models.User {
	v, b := c.GetSession(SESSTION_KEY_USER).(models.User)
	if b {
		return &v
	}

	return nil
}

func (c *BaseController) Path(idx int) string {
	path := c.Ctx.Request.URL.Path
	results := strings.Split(strings.TrimPrefix(path,"/"), "/")
	if idx < len(results) {
		return results[idx]
	}

	return ""
}

func (c *BaseController) PathValue() string {
	return c.Path(2)
}

func (c *BaseController) PathValueInt() int {
	return utils.JKStrToInt(c.Path(2))
}



