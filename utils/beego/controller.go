package beego

import (
	"github.com/astaxie/beego"
	"github.com/cheneylew/goutil/utils"
	"strings"
	"net/url"
)

var SESSTION_KEY_USER string

func init() {
	SESSTION_KEY_USER = "LOGINED_USER"
}

type BBaseController struct {
	beego.Controller
}

func (c *BBaseController) IsPost() bool {
	return c.Ctx.Request.Method == "POST"
}

func (c *BBaseController) IsGet() bool {
	return c.Ctx.Request.Method == "GET"
}

func (c *BBaseController) PostForm() url.Values {
	return c.Ctx.Request.PostForm
}

func (c *BBaseController) PostFormWithKey(key string) []string {
	return c.Ctx.Request.PostForm[key]
}

func (c *BBaseController) Finish() {
	c.Controller.Finish()
}

//通用
func (c *BBaseController) RedirectWithURL(url string) {
	c.Redirect(url, 302)
}

//用户
func (c *BBaseController) IsLogin() bool {
	return c.GetSession(SESSTION_KEY_USER) != nil
}

func (c *BBaseController) SaveUser(user interface{}) {
	c.SetSession(SESSTION_KEY_USER,user)
}

func (c *BBaseController) UserLogout() {
	c.SetSession(SESSTION_KEY_USER,nil)
}

func (c *BBaseController) GetUser() interface{} {
	return c.GetSession(SESSTION_KEY_USER)
}

func (c *BBaseController) Path(idx int) string {
	path := c.Ctx.Request.URL.Path
	results := strings.Split(strings.TrimPrefix(path,"/"), "/")
	if idx < len(results) {
		return results[idx]
	}

	return ""
}

func (c *BBaseController) PathValue() string {
	return c.Path(2)
}

func (c *BBaseController) PathValueInt() int {
	return utils.JKStrToInt(c.Path(2))
}




