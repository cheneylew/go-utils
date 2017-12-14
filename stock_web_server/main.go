package main

import (
	_ "github.com/cheneylew/goutil/stock_web_server/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/cheneylew/goutil/utils"
	"github.com/cheneylew/goutil/stock_web_server/models"
)

func init() {

}

func beegoRun()  {
	beego.Run()
}

func main() {

	beegoRun()

}

