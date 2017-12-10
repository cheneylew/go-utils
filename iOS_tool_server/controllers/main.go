package controllers

import (
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
	c.TplName = "main.html"
}

