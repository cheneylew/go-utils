package controllers

import (
	"github.com/cheneylew/goutil/stock_web_server/database"
	"github.com/cheneylew/goutil/stock_web_server/stock"
	"github.com/cheneylew/goutil/utils"
	"github.com/cheneylew/goutil/stock_web_server/models"
	"time"
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
	c.TplName = "index.html"
}

func (c *MainController) Index() {
	c.TplName = "main.html"

	c.Data["Stocks"] = stock.AnalysStockInfo()
}

func (c *MainController) Test() {
	c.TplName = "main.html"
}

func (c *MainController) IsUp() {
	code := c.GetString("code")
	klines := database.DB.GetKLineAllForStockCode(code)
	isUp,_,_ := stock.KLineIsUp(klines)
	utils.JJKPrintln(isUp)
	c.TplName = "main.html"
}

func (c *MainController) Red() {
	c.TplName = "red.html"
	days := c.GetString("days", "15")

	results := stock.AnalysRedRate(utils.JKStrToInt(days))
	var s models.SortAnalysDayKLins
	for _, value := range results {
		rate := value.UpDownRateTotal*100
		value.UpDownRateTotal = rate
		if  rate > 15 {
			utils.JJKPrintln(value.UpDownRateTotal*100)
			s = append(s, value)
		}
	}

	c.Data["AnalysDayKLines"] = s
}

func (c *MainController) MainPower() {
	stock.InitCache()

	c.TplName = "mainpower.html"
	stocks := stock.AnalysStockInfo()
	c.Data["AnalysDayKLines"] = stocks
}

func (c *MainController) UpStock()   {
	stocks := database.DB.GetStockWithCodePrefix("00")
	stocks = append(stocks, database.DB.GetStockWithCodePrefix("60")...)
	var ss []*models.Stock
	for _, value := range stocks {
		klines := database.DB.GetKLineAllForStockCode(value.Code)
		up,_,_ := stock.KLineIsUp(klines)
		if up && klines[len(klines)-1].Date.After(time.Now().Add(-time.Hour*24*3)) {
			ss = append(ss, value)
		}
	}

	c.Data["Stocks"] = ss
	c.TplName = "main.html"
}

func (c *MainController) MainIn()   {
	stock.InitCache()
	stocks := stock.Analys5MainInStocks()
	c.Data["Stocks"] = stocks
	c.TplName = "main.html"
}

func (c *MainController) MainInThreedays()   {
	stock.InitCache()
	stocks := stock.Analys5MainInGreatThan3DaysStocks()
	c.Data["Stocks"] = stocks
	c.TplName = "main.html"
}

func (c *MainController) MainInTwodays()   {
	stock.InitCache()
	stocks := stock.Analys5MainInGreatThan2DaysStocks()
	c.Data["Stocks"] = stocks
	c.TplName = "main.html"
}

func (c *MainController) MainInEqualTwodays()   {
	stock.InitCache()
	stocks := stock.Analys5MainInEqual2DaysStocks()
	c.Data["Stocks"] = stocks
	c.TplName = "main.html"
}


func (c *MainController) MainOut()   {
	stock.InitCache()
	stocks := stock.Analys5MainOutStocks()
	c.Data["Stocks"] = stocks
	c.TplName = "main.html"
}


func (c *MainController) DownDays()   {
	stock.InitCache()
	stocks := stock.AnalysBuyWhat()
	c.Data["Stocks"] = stocks
	c.TplName = "main5.html"
}

func (c *MainController) AllDowloadTask()  {
	stock.DownloadTaskAll()
	c.Ctx.WriteString("all task finished!")
}

func (c *MainController) StockSortByChangeHand()  {
	var stocks []*models.Stock
	database.DB.Orm.Raw("SELECT * FROM stock.stock where flow_amount > 100 and change_hand_rate > 3.0 order by change_hand_rate desc limit 0, 10000;").QueryRows(&stocks)
	c.Data["Stocks"] = stocks
	c.TplName = "main5.html"
}

