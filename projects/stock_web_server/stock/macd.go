package stock

import (
	"math"
	"github.com/cheneylew/goutil/projects/stock_web_server/database"
	"github.com/cheneylew/goutil/projects/stock_web_server/models"
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

func ReCalculateMACD()  {
	stocks := CCGetStockAll()
	var ps []interface{}
	for _, value := range stocks {
		ps = append(ps, value)
	}
	utils.QueueTask(10, ps, func(idx int, param interface{}) {
		stock := param.(*models.Stock)
		calMacdRecent(stock)
	})
}

func calMacdRecent(stock *models.Stock)  {
	klines := CCGetKLinesWithCode(stock.Code, 20)
	if len(klines) == 0 {
		return
	}
	if klines[0].Ema12 == 0 || klines[0].Ema26 == 0 {
		//panic("该股票没有计算过MACD，需要重头计算")
		return
	}

	var lastKLine *models.KLine
	for index, kline := range klines {
		if index == 0 {
			lastKLine = kline
		} else {
			nowEMA12 := lastKLine.Ema12 * 11.0 / 13.0 + kline.ClosingPrice*2.0/13.0
			nowEMA26 := lastKLine.Ema26 * 25.0 / 27.0 + kline.ClosingPrice*2.0/27.0
			nowDIF := nowEMA12 - nowEMA26
			nowDEA := lastKLine.Dea * 8.0 / 10.0 + nowDIF*2.0/10.0
			nowBAR := 2*(nowDIF - nowDEA)

			kline.Ema12 = nowEMA12
			kline.Ema26 = nowEMA26
			kline.Dif = nowDIF
			kline.Dea = nowDEA
			kline.Bar = nowBAR

			database.DB.Orm.Update(kline)
			lastKLine = kline
		}
	}
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


func calculateMACD(code string, count int64) (todayDIF float64, lines []*models.KLine)  {
	klines := GetStockDayKLine(code,count)
	if len(klines) <= 14 {
		return 0, nil
	}

	var lastEMA12 float64 = 0
	var lastEMA26 float64 = 0
	var lastDEA float64 = 0
	for index, kline := range klines {
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

			klines[index].Ema12 = nowEMA12
			klines[index].Ema26 = nowEMA26
			klines[index].Dif = nowDIF
			klines[index].Dea = nowDEA
			klines[index].Bar = nowBAR
		} else {
			lastEMA12 = kline.ClosingPrice
			lastEMA26 = kline.ClosingPrice
			lastDEA = 0
		}
	}

	return klines[len(klines)-1].Dif, klines
}
