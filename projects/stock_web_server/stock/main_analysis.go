package stock

import (
	"github.com/cheneylew/goutil/utils"
	"math"
	"github.com/cheneylew/goutil/projects/stock_web_server/database"
	"github.com/cheneylew/goutil/projects/stock_web_server/models"
	"time"
	"github.com/panshiqu/dysms"
	"fmt"
	"github.com/astaxie/beego/config"
	"strings"
	"sort"
	"path"
)

var conf config.Configer

func valueWithKey(key string) string {
	if conf == nil {
		iniconf, _ := config.NewConfig("ini", "conf/app.conf")
		conf = iniconf

	}
	return conf.String(key)
}

func Main_rsi()  {
	if false {
		rsi, _ := calculateRSI("sz000651",50)
		utils.JJKPrintln(rsi)
	}

	if false {
		sendSMS("000651")
	}

	//rsi开始回调的股票
	if false {
		shStocks := database.DB.GetStockWithCodePrefix("60")
		for _, stock := range shStocks {
			_, lines := calculateRSI(stock.CodeStr(),50)
			if len(lines) > 2 {
				last1 := lines[len(lines)-1]
				last2 := lines[len(lines)-2]
				last3 := lines[len(lines)-3]
				//last4 := lines[len(lines)-4]
				isRecent := last1.Date.After(time.Now().Add(time.Hour*24*-7))
				amountOk := stock.FlowAmount > 300

				//rsi低位连涨两天，第三天9:45入手，rsi低于20靠谱点
				condition1 := last3.Rsi < 27.5 && last3.Rsi < last2.Rsi && last2.Rsi < last1.Rsi
				//rsi低位涨一天跌一天，再涨一天。跌的一天rsi高于低位
				condition2 := false //last4.Rsi < 27.5 && last4.Rsi < last3.Rsi && last3.Rsi > last2.Rsi && last2.Rsi < last1.Rsi && last4.Rsi < last2.Rsi
				if isRecent && amountOk && (condition1 || condition2) {
				//if isRecent && last2.Rsi < 25 && last2.Rsi < last1.Rsi{
				//if isRecent && last4.Rsi < 25 && last4.Rsi < last3.Rsi && last3.Rsi > last2.Rsi && last2.Rsi < last1.Rsi && last2.Rsi > last4.Rsi {//回来
					utils.JJKPrintln(fmt.Sprintf("code:%s", stock.Code), fmt.Sprintf(" rsi:%.2f", last1.Rsi), fmt.Sprintf(" 价格:%.2f", last1.ClosingPrice),  fmt.Sprintf(" 总市值：%.2f", stock.FlowAmount))
				}
			}
		}
	}

	//监控某只几只股票
	if false {
		observeStocks()
	}

	//测试胜率
	if false {
		codes := strings.Split("601952","|")
		for _, code := range codes {
			stock := database.DB.GetStockWithCode(code)
			_, lines := calculateRSI(stock.CodeStr(),400)
			utils.JJKPrintln(len(lines))
			if len(lines) > 2 {
				for i:=10; i<len(lines); i++ {
					last1 := lines[i-1]	//买入点
					last2 := lines[i-2]
					last3 := lines[i-3]
					last4 := lines[i-4]

					//rsi低位连涨两天，第三天9:45入手，rsi低于20靠谱点
					condition1 := last3.Rsi < 27.5 && last3.Rsi < last2.Rsi && last2.Rsi < last1.Rsi
					//rsi低位涨一天跌一天，再涨一天。跌的一天rsi高于低位
					condition2 := last4.Rsi < 27.5 && last4.Rsi < last3.Rsi && last3.Rsi > last2.Rsi && last2.Rsi < last1.Rsi && last4.Rsi < last2.Rsi
					if condition1 || condition2 {
						utils.JJKPrintln(stock.Code, last1.Date, last1.Rsi, last1.ClosingPrice,  stock.FlowAmount)
					}
				}
			}
		}
	}

	//macd统计方式
	if true {
		shStocks := database.DB.GetStockWithCodePrefix("60")
		var stocks []*models.Stock
		for _, stock := range shStocks {
			_, lines := calculateMACD(stock.CodeStr(), 400)
			if len(lines) > 2 {
				last1 := lines[len(lines)-1]
				//last5 := lines[len(lines)-6]
				isRecent := last1.Date.After(time.Now().Add(time.Hour*24*-7))
				//空方DIF趋势向上,发生在0轴附近
				//condition1 := last1.Dif > last5.Dif && math.Abs(last1.Dif)<0.05
				//多方DIF趋势向上,发生在多方，一直往上涨。
				//condition1 := last1.Dif > last5.Dif && last5.Dif>0 && lines[len(lines)-12].Dif > 0 && lines[len(lines)-24].Dif > 0
				//空方DIF趋势向上,发生在空方，一直往上涨。
				//condition1 := last1.Dif > last5.Dif && last1.Dif < -0.1 && last1.Dif > last1.Dea+0.1
				//多方DIF趋势向上,马上突破0轴，一直往上涨。
				//condition1 := last1.Dif > last5.Dif && last5.Dif<=-0.1 && last1.Dif >= -0.05 && last1.Dif <= 0
				//macd发生交叉
				//condition1 := math.Abs((last1.Dif-last1.Dea)) <= 0.05 && lines[len(lines)-9].Dif < -0.2
				//多方DIF趋势向上,逼空点
				//condition1 := lines[len(lines)-3].Bar>0 && lines[len(lines)-1].Bar>lines[len(lines)-2].Bar && lines[len(lines)-2].Bar > lines[len(lines)-3].Bar && lines[len(lines)-3].Bar < lines[len(lines)-4].Bar && lines[len(lines)-4].Bar < lines[len(lines)-5].Bar
				//多方DIF趋势向上,多方绿柱抽脚
				//condition1 := lines[len(lines)-1].Bar<(-0.13*2) && lines[len(lines)-1].Bar<0 && lines[len(lines)-1].Bar>lines[len(lines)-2].Bar&&lines[len(lines)-2].Bar<lines[len(lines)-3].Bar
				//多方DIF趋势向上，绿柱消失，出现第一根红柱子
				//condition1 := lines[len(lines)-1].Dif>0 && lines[len(lines)-1].Bar>0 && lines[len(lines)-2].Bar<0
				//红柱汤匙形态发射
				condition1 := true
				days := 6
				for i:=1;i<=days ; i++ {
					if lines[len(lines)-i].Bar<lines[len(lines)-i-1].Bar {
						condition1 = false
					}
				}
				if lines[len(lines)-days-1].Bar<0 {
					condition1 = false
				}

				if isRecent && len(lines) > 100 && condition1 {
					utils.JJKPrintln("code:",stock.Code," dif:", last1.Dif," dea:", last1.Dea)
					stocks = append(stocks, stock)
				}
			}
		}

		//按换手率活跃度排序
		sort.Slice(stocks, func(i, j int) bool {
			return stocks[i].ChangeHandRate < stocks[j].ChangeHandRate
		})
		for _, value := range stocks {
			utils.JJKPrintln(value.Code, value.ChangeHandRate)
		}
	}

	//KDJ选择下叉方式
	if false {
		shStocks := database.DB.GetStockWithCodePrefix("60")
		var stocks []*models.Stock
		for _, stock := range shStocks {
			_, lines := calculateKDJ(stock.CodeStr(), 400)
			if len(lines) > 2 {
				curDayStock := lines[len(lines)-1]
				yestodayStock := lines[len(lines)-2]
				condition1 := math.Abs(curDayStock.Kdj_k-curDayStock.Kdj_d) < 2 && curDayStock.Kdj_k < 45 && curDayStock.Kdj_k>=curDayStock.Kdj_d && yestodayStock.Kdj_k < curDayStock.Kdj_k
				isRecent := curDayStock.Date.After(time.Now().Add(time.Hour*24*-7))
				if isRecent && len(lines) > 100 && condition1 {
					utils.JJKPrintln("code:",stock.Code," k:", curDayStock.Kdj_k," d:", curDayStock.Kdj_d," 换手率:",stock.ChangeHandRate)
					stocks = append(stocks, stock)
				}
			}
		}

		//按换手率活跃度排序
		sort.Slice(stocks, func(i, j int) bool {
			return stocks[i].ChangeHandRate < stocks[j].ChangeHandRate
		})
		for _, value := range stocks {
			utils.JJKPrintln(value.Code, value.ChangeHandRate)
		}

		writeStocks(stocks)
	}

	//测试胜率KDJ
	if false {
		codes := strings.Split("600271","|")
		for _, code := range codes {
			stock := database.DB.GetStockWithCode(code)
			_, lines := calculateKDJ(stock.CodeStr(),400)
			utils.JJKPrintln(len(lines))
			if len(lines) > 2 {
				for i:=10; i<len(lines); i++ {
					curDayStock := lines[i-1]
					yestodayStock := lines[i-2]
					condition1 := math.Abs(curDayStock.Kdj_k-curDayStock.Kdj_d) < 2 && curDayStock.Kdj_k < 45 && yestodayStock.Kdj_k < curDayStock.Kdj_k
					if condition1 {
						utils.JJKPrintln(curDayStock.Date)
					}
				}
			}
		}
	}

	//测试胜率MACD底背离
	if false {
		codes := strings.Split("603920","|")
		for _, code := range codes {
			stock := database.DB.GetStockWithCode(code)
			_, lines := calculateMACD(stock.CodeStr(),400)
			length := len(lines)
			utils.JJKPrintln(length)
			days := 40
			if len(lines) > 2 {
				for i:=days; i<length; i++ {
					minPrice := 10000.0
					minSencondPrice := 10000.0
					maxPrice := 0.0
					var minDate time.Time
					//var maxDate time.Time
					var minSencondDate time.Time
					var minKLine *models.KLine
					var minSecKLine *models.KLine

					for j:=i-days; j<=i; j++ {
						if lines[j].MaxPrice > maxPrice {
							maxPrice = lines[j].MaxPrice
							//maxDate = lines[j].Date
						}

						if lines[j].MinPrice <= minPrice {
							minPrice = lines[j].MinPrice
							minDate = lines[j].Date
							minKLine = lines[j]
						}

						ok := false
						if math.Abs(float64(lines[j].Date.Unix()-minDate.Unix())) < 3600*24 *3 {
							ok = true
						}
						if lines[j].MinPrice <= minSencondPrice && lines[j].MinPrice > minPrice && ok {
							minSencondPrice = lines[j].MinPrice
							minSencondDate = lines[j].Date
							minSecKLine = lines[j]
						}
					}
					if minKLine != nil && minSecKLine != nil && math.Abs(float64(minKLine.Date.Unix()-minSecKLine.Date.Unix())) > 3600*24*10 {
						if minKLine.Dif > minSecKLine.Dif && minKLine.MinPrice < minSecKLine.MinPrice {
							utils.JJKPrintln(" minSec:", minSencondPrice, minSencondDate," min:",minPrice, minDate)
						}
					}
				}
			}
		}
	}
}

func GetObserverStocksFilePath() string {
	return path.Join("/Users/dejunliu/Desktop", "stocks.txt")
}

func observeStocks()  {
	stocks := valueWithKey("stocks")
	//定时执行
	utils.CronJob("00 05 15 * * 1-5", func() {
		//rsi低位监控，发现股票rsi处于低位，发送短信提示
		if true {
			codes := strings.Split(stocks,"|")
			for _, code := range codes {
				stock := database.DB.GetStockWithCode(code)
				_, lines := calculateRSI(stock.CodeStr(),50)
				if len(lines) > 2 {
					last1 := lines[len(lines)-1]
					last2 := lines[len(lines)-2]
					last3 := lines[len(lines)-3]
					isRecent := last1.Date.After(time.Now().Add(time.Hour*24*-7))
					//rsi低位连涨两天，第三天9:45入手，rsi低于20靠谱点
					if isRecent && last3.Rsi < 27.5 && last3.Rsi < last2.Rsi && last2.Rsi < last1.Rsi {
						utils.JJKPrintln(stock.Code, last1.Rsi, last1.ClosingPrice,  stock.FlowAmount)
						sendSMS(stock.Code)
					}
				}
			}
		}

		//macd统计方式
		if true {
			//下载股票信息，总市值等
			downloadStockInfo()
			//计算macd
			shStocks := database.DB.GetStockWithCodePrefix("60")
			shStocks = append(shStocks, database.DB.GetStockWithCodePrefix("00")...)
			var stocks []*models.Stock
			for _, stock := range shStocks {
				_, lines := calculateMACD(stock.CodeStr(), 400)
				if len(lines) > 2 {
					last1 := lines[len(lines)-1]
					last5 := lines[len(lines)-6]
					isRecent := last1.Date.After(time.Now().Add(time.Hour*24*-7))
					//空方DIF趋势向上,发生在0轴附近
					//condition1 := last1.Dif > last5.Dif && math.Abs(last1.Dif)<0.05
					//多方DIF趋势向上,发生在多方，一直往上涨。
					//condition1 := last1.Dif > last5.Dif && last5.Dif>0
					//多方DIF趋势向上,马上突破0轴，一直往上涨。
					condition1 := last1.Dif > last5.Dif && last5.Dif<=-0.1 && last1.Dif >= -0.05 && last1.Dif <= 0
					//macd发生交叉
					//condition1 := math.Abs((last1.Dif-last1.Dea)) <= 0.05 && lines[len(lines)-9].Dif < -0.2
					if isRecent && len(lines) > 100 && condition1 {
						utils.JJKPrintln("code:",stock.Code," dif:", last1.Dif," dea:", last1.Dea)
						stocks = append(stocks, stock)
					}
				}
			}

			//按换手率活跃度排序
			sort.Slice(stocks, func(i, j int) bool {
				return stocks[i].ChangeHandRate > stocks[j].ChangeHandRate
			})

			writeStocks(stocks)
		}
	})
}

func writeStocks(stocks []*models.Stock)  {
	file := GetObserverStocksFilePath()
	text := ""
	for _, value := range stocks {
		text = fmt.Sprintf("%s%s %.2f\n", text,value.Code, value.ChangeHandRate)
		utils.JJKPrintln(value.Code, value.ChangeHandRate)
	}

	utils.FileWriteString(file, text)
}

func sendSMS(stockCode string)  {
	mobile := "17602125152"
	if err := dysms.SendSms("LTAIwsROTsyzMq1a", "KIJzyH67zdskMBdjOIYKnDdv7xavX7", mobile, "爱编程", fmt.Sprintf(`{"code":"%s"}`, stockCode), "SMS_128535187"); err != nil {
		utils.JJKPrintln("dysms.SendSms", err)
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
func calculateRSI(code string, count int64) (todayRSI float64, lines []*models.KLine)  {
	klines := GetStockDayKLine(code,count)
	return klines[len(klines)-1].Rsi, lines
}

func calculateRSIWithLines(klines []*models.KLine) []*models.KLine  {
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

				klines[i].Rsi = RSI
			}
		}
	}

	return klines
}
