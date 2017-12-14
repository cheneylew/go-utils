package controllers

import (
	"github.com/cheneylew/goutil/stock_web_server/database"
	"github.com/cheneylew/goutil/utils"

	"github.com/cheneylew/goutil/stock_web_server/stock"
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

	s := stock.GetStockDayKLine("sz000725",1)
	utils.JJKPrintln(s)
	a,e := database.DB.Orm.Insert(s[0])
	utils.JJKPrintln(a,e)
}

