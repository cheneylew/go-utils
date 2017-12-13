package controllers

import (
	"github.com/cheneylew/goutil/stock_web_server/database"
)

type MainController struct {
	BaseController
}

func (c *MainController) Prepare() {
	c.BaseController.Prepare()
}

func (c *MainController) Finish() {
	c.Controller.Finish()
}

func (c *MainController) Get() {
	c.RedirectWithURL("/main/index")
}

func (c *MainController) Index() {
	database.DB.HelloWorld()
	c.TplName = "main.html"
}

