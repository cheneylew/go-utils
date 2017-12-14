package main

import (
	_ "github.com/cheneylew/goutil/stock_web_server/routers"
	"github.com/astaxie/beego"
	"github.com/cheneylew/goutil/stock_web_server/stock"
)

func init() {

}

func beegoRun()  {
	beego.Run()
}

func main() {
	//StockTestMain()
	stock.GetRealTimeStockInfo("sh600703")
	beegoRun()

}

