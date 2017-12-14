package main

import (
	_ "github.com/cheneylew/goutil/stock_web_server/routers"
	"github.com/astaxie/beego"
)

func init() {

}

func beegoRun()  {
	beego.Run()
}

func main() {
	StockTestMain()
	beegoRun()

}

