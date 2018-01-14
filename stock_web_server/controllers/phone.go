package controllers

import "github.com/cheneylew/goutil/utils"

type PhoneController struct {
	BaseController
}

func (c *PhoneController) Prepare() {
	c.BaseController.Prepare()
}

func (c *PhoneController) Finish() {
	c.Controller.Finish()
}

func (c *PhoneController) Test() {
	c.SetSecureCookie("secret", "cheneylew","this is value")
	cookie, err := c.Ctx.Request.Cookie("cheneylew")
	if err != nil {
		utils.JJKPrintln(err)
	} else {
		utils.JJKPrintln("name:",cookie.Name, "value:", utils.Base64Decode(cookie.Value))
	}
	c.TplName = "test.html"
}

func (c *PhoneController) Testa() {
	c.Ctx.WriteString("aaaaaaa")
	//c.Redirect("/phone/testb", 301)
}

func (c *PhoneController) Testb() {
	c.Ctx.WriteString("bbbbbbbb")
}
