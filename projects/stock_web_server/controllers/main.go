package controllers

import (
	"github.com/cheneylew/goutil/projects/stock_web_server/database"
	"github.com/cheneylew/goutil/projects/stock_web_server/stock"
	"github.com/cheneylew/goutil/utils"
	"github.com/cheneylew/goutil/projects/stock_web_server/models"
	"time"
	"sort"
	"strings"
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
	c.Ctx.WriteString("hello world")
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

func (c *MainController) Up()   {
	stock.InitCache()
	stocks := stock.AnalysBuyUp()
	c.Data["Stocks"] = stocks
	c.TplName = "main5.html"
}

func (c *MainController) AllDowloadTask()  {
	stock.DownloadTaskAll()
	c.Ctx.WriteString("all task finished!")
}

func (c *MainController) AllDowloadTaskAdds()  {
	stock.DownloadTaskAddKLines()
	c.Ctx.WriteString("all task finished!")
}

func (c *MainController) StockSortByChangeHand()  {
	var stocks []*models.Stock
	database.DB.Orm.Raw("SELECT * FROM stock.stock where change_hand_rate > 3.0 order by change_hand_rate desc limit 0, 10000;").QueryRows(&stocks)
	c.Data["Stocks"] = stocks
	c.TplName = "main5.html"
}

func (c *MainController) MacdGold()   {
	stock.InitCache()
	stocks := stock.AnalysMACD(0)
	var ps []interface{}
	for _, value := range stocks {
		ps = append(ps, value)
	}

	newStocks := utils.Filter(ps, func(i interface{}, i2 int) bool {
		stock := i.(*models.Stock)
		return stock.ChangeHandRate > 1.0
	})

	sort.Slice(newStocks, func(i, j int) bool {
		//return newStocks[i].(*models.Stock).GreenBarCount < newStocks[j].(*models.Stock).GreenBarCount
		return newStocks[i].(*models.Stock).ChangeHandRate > newStocks[j].(*models.Stock).ChangeHandRate
	})

	c.Data["Stocks"] = newStocks
	c.TplName = "main5.html"
}

func (c *MainController) AnalysNewStocks()  {
	stock.InitCache()

	newStocks := stock.AnalysNewStocks()
	sort.Slice(newStocks, func(i, j int) bool {
		//return newStocks[i].(*models.Stock).GreenBarCount < newStocks[j].(*models.Stock).GreenBarCount
		return newStocks[i].ChangeHandRate > newStocks[j].ChangeHandRate
	})
	c.Data["Stocks"] = newStocks
	c.TplName = "main5.html"
}

func (c *MainController) Macd()  {
	//stock.InitCache()
	file := stock.GetObserverStocksFilePath()
	text := utils.FileReadAllString(file)
	lines := strings.Split(text,"\n")
	var stocks []*models.Stock
	for _, line := range lines {
		infos := strings.Split(line," ")
		if len(infos) > 1 {
			code := infos[0]
			stock := database.DB.GetStockWithCode(code)
			stocks = append(stocks, stock)
		}
	}

	c.Data["Stocks"] = stocks
	c.TplName = "main5.html"
}

func (c *MainController) Push()  {
	stock.PushNotification("hello!")
	c.Ctx.WriteString("have sent!")
}

