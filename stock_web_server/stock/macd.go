package stock

import (
	"math"
	"github.com/cheneylew/goutil/stock_web_server/database"
	"github.com/cheneylew/goutil/stock_web_server/models"
	"github.com/cheneylew/goutil/utils"
)

func CalculateMACD()  {
	stocks := CCGetStockAll()
	var ps []interface{}
	for _, value := range stocks {
		ps = append(ps, value)
	}
	utils.QueueTask(10, ps, func(idx int, param interface{}) {
		stock := param.(*models.Stock)
		macd(stock)
	})
}

func macd(stock *models.Stock)  {
	onlineKLines := GetStockDayKLine(stock.CodeStr(), math.MaxInt64)

	var lastEMA12 float64 = 0
	var lastEMA26 float64 = 0
	var lastDEA float64 = 0
	for index, kline := range onlineKLines {
		if index > 0 {
			nowEMA12 := lastEMA12 * 11.0 / 13.0 + kline.ClosingPrice*2.0/13.0
			nowEMA26 := lastEMA26 * 25.0 / 27.0 + kline.ClosingPrice*2.0/27.0
			nowDIF := nowEMA12 - nowEMA26
			nowDEA := lastDEA * 8.0 / 10.0 + nowDIF*2.0/10.0
			nowBAR := 2*(nowDIF - nowDEA)
			//utils.JJKPrintln(fmt.Sprintf("%s,%f,%f,dif=%f,dea=%f,macd_bar=%f", kline.Date,nowEMA12, nowEMA26, nowDIF, nowDEA, nowBAR))

			lastEMA12 = nowEMA12
			lastEMA26 = nowEMA26
			lastDEA = nowDEA

			dbKLine := CCGetKLineWithCodeAndDate(stock.Code, kline.Date)
			if dbKLine != nil {
				dbKLine.Ema12 = nowEMA12
				dbKLine.Ema26 = nowEMA26
				dbKLine.Dif = nowDIF
				dbKLine.Dea = nowDEA
				dbKLine.Bar = nowBAR
				database.DB.Orm.Update(dbKLine)
			}
		} else {
			lastEMA12 = kline.ClosingPrice
			lastEMA26 = kline.ClosingPrice
			lastDEA = 0
		}
	}
}
