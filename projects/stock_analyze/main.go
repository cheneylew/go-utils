package main

import (
	_ "github.com/cheneylew/goutil/projects/stock_analyze/routers"
	"github.com/astaxie/beego"
	"github.com/cheneylew/goutil/projects/stock_analyze/stock"
)

func init() {

}

func beegoRun()  {
	beego.Run()
}

func main() {
	stock.StockTestMain()
	beegoRun()
}

