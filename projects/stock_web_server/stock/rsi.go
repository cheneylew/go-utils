package stock

import (
	"github.com/cheneylew/goutil/utils"
	"math"
	"github.com/cheneylew/goutil/projects/stock_web_server/database"
	"github.com/cheneylew/goutil/projects/stock_web_server/models"
	"time"
)

func Main_rsi()  {
	if false {
		rsi, _ := calculateRSI("sz000651")
		utils.JJKPrintln(rsi)
	}

	//rsi开始回调的股票
	if false {
		shStocks := database.DB.GetStockWithCodePrefix("60")
		for _, value := range shStocks {
			_, lines := calculateRSI(value.CodeStr())
			if len(lines) > 2 {
				lastOne := lines[len(lines)-1]
				lastSecondOne := lines[len(lines)-2]
				isRecent := lastOne.Date.After(time.Now().Add(time.Hour*24*-3))
				if lastSecondOne.Rsi < 20 && value.FlowAmount > 50 && lastOne.Rsi > lastSecondOne.Rsi && isRecent {
					utils.JJKPrintln(lastOne.Rsi, lastOne.ClosingPrice, value.Code, value.FlowAmount)
				}
			}
		}
	}

}

/*
The very first calculations for average gain and average loss are simple 14-period averages.
First Average Gain = Sum of Gains over the past 14 periods / 14.
First Average Loss = Sum of Losses over the past 14 periods / 14
The second, and subsequent, calculations are based on the prior averages and the current gain loss:
Average Gain = [(previous Average Gain) x 13 + current Gain] / 14.
Average Loss = [(previous Average Loss) x 13 + current Loss] / 14.
http://stockcharts.com/school/doku.php?id=chart_school:technical_indicators:relative_strength_index_rsi
 */
func calculateRSI(code string) (todayRSI float64, lines []*models.KLine)  {
	todayRsiValue := 0.0

	klines := GetStockDayKLine(code,50)
	//N日RSI =N日内收盘涨幅的平均值/(N日内收盘涨幅均值+N日内收盘跌幅均值) ×100%
	dayNum := 6
	prevAverageGain := 0.0
	prevAverageLoss := 0.0
	for i:=0; i<len(klines) ; i++  {
		if i>=dayNum {
			lastIndex := i-dayNum
			if lastIndex == 0 { //第一次求平均
				riseCount := 0.0
				fallCount := 0.0

				prevIndex := lastIndex
				for j:= lastIndex+1; j<= i;j++ {
					deltaPrice := klines[j].ClosingPrice-klines[prevIndex].ClosingPrice
					if deltaPrice > 0 {
						riseCount += deltaPrice
					} else {
						fallCount += math.Abs(deltaPrice)
					}
					prevIndex = j
				}

				AverageGain := riseCount/float64(dayNum)
				AverageLoss := fallCount/float64(dayNum)
				//RS := AverageGain/AverageLoss
				//RSI := 100.0-(100.0/(1+RS))

				prevAverageGain = AverageGain
				prevAverageLoss = AverageLoss

				//utils.JJKPrintln(RSI,klines[i].Date)
			} else { //根据前一次推算
				riseCount := 0.0
				fallCount := 0.0

				deltaPrice := klines[i].ClosingPrice-klines[i-1].ClosingPrice
				if deltaPrice > 0 {
					riseCount = prevAverageGain*float64((dayNum-1))+deltaPrice
					fallCount = prevAverageLoss*float64(dayNum-1)+0
				} else {
					riseCount = prevAverageGain*float64(dayNum-1)+0
					fallCount = prevAverageLoss*float64(dayNum-1)+math.Abs(deltaPrice)
				}

				AverageGain := riseCount/float64(dayNum)
				AverageLoss := fallCount/float64(dayNum)
				RS := AverageGain/AverageLoss
				RSI := 100.0-(100.0/(1+RS))

				prevAverageGain = math.Abs(AverageGain)
				prevAverageLoss = math.Abs(AverageLoss)

				todayRsiValue = RSI
				klines[i].Rsi = todayRsiValue
			}
		}
	}

	return todayRsiValue, klines
}
