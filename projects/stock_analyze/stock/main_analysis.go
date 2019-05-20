package stock

import (
	"github.com/cheneylew/goutil/utils"
	"math"
	"github.com/cheneylew/goutil/projects/stock_analyze/database"
	"github.com/cheneylew/goutil/projects/stock_analyze/models"
	"time"
	"github.com/panshiqu/dysms"
	"fmt"
	"github.com/astaxie/beego/config"
	"strings"
	"sort"
	"path"
)

func valueWithKey(key string) string {
	conf, _:= config.NewConfig("ini", "conf/app.conf")
	return conf.String(key)
}

func Main_rsi()  {
	//DownloadTaskAddKLines()

	//监控某只几只股票
	if false {
		utils.JJKPrintln("开始监控股票...")
		//定时执行
		utils.CronJob("00 05 13 * * 1-5", func() {
			observeStocks()
		})
	}

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
			_, lines := calculateRSI(stock.Code,50)
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

	//测试胜率
	if false {
		codes := strings.Split("601952","|")
		for _, code := range codes {
			stock := database.DB.GetStockWithCode(code)
			_, lines := calculateRSI(stock.Code,400)
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
	if false {
		////下载换手率
		//InitCache()
		//downloadStockInfo()
		//utils.JJKPrintln("update stock infos ok!")
		//分析股票
		shStocks := database.DB.GetStockWithCodePrefix("60")
		shStocks = append(shStocks, database.DB.GetStockWithCodePrefix("00")...)
		winCount := 0
		var stocks []*models.Stock
		var winStocks []*models.Stock
		var failedStocks []*models.Stock
		for _, stock := range shStocks {
			_, lines := calculateMACD(stock.Code, 400)
			lines = calculateKDJWithLines(lines)
			if len(lines) > 20 {
				delta := 0
				deltaDays := 0
				last1 := lines[len(lines)-1-delta]
				//last2 := lines[len(lines)-2-delta]
				//last5 := lines[len(lines)-6-delta]
				//isRecent := true
				isRecent := last1.Date.After(time.Now().Add(time.Hour*24*(-7-time.Duration(delta))))
				//追涨停
				//rate := last1.GetAddRate(last2)
				//condition1 := rate >= 0.08 && rate <= 0.1
				//购买3天后，胜率45%。多方DIF趋势向上,发生在多方，一直往上涨。
				//condition1 := last1.Dif > last5.Dif && last5.Dif>0 && lines[len(lines)-12-delta].Dif > 0 && lines[len(lines)-24-delta].Dif > 0
				//购买3天后，胜率60%。多方,线上阴线买。暴涨的不能买。
				//rate := last1.GetAddRate(last2)
				//condition1 := last1.Dif>0 &&  rate > -0.035 && rate < -0.025 && last1.Dif > last5.Dif
				//空方DIF趋势向上,发生在多方
				condition1 := lines[len(lines)-1].Dif > lines[len(lines)-8].Dif && lines[len(lines)-8].Dif > lines[len(lines)-16].Dif && lines[len(lines)-1].Dif>0.1
				//空方DIF趋势向上,发生在0轴附近
				//condition1 := last1.Dif > last5.Dif && math.Abs(last1.Dif)<0.05
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
				//condition1 := true
				//days := 6
				//for i:=1;i<=days ; i++ {
				//	if lines[len(lines)-i-delta].Bar<lines[len(lines)-i-delta-1].Bar {
				//		condition1 = false
				//	}
				//}
				//if lines[len(lines)-days-delta-1].Bar<0 {
				//	condition1 = false
				//}
				//DIF>0,KDJ下叉。多方KDJ金叉。
				//condition1 := last1.Kdj_k>last1.Kdj_d && last2.Kdj_k<last2.Kdj_d && last1.Dif > -0.5
				//DIF>0,KDJ死叉
				//condition1 := last1.Kdj_k<last1.Kdj_d && last2.Kdj_k>last2.Kdj_d
				//中短线买入条件(DIF>5,明显向上突破0轴,MACD红柱发散,收阳线,MACD>5,明显突破0轴)
				//mcount := 0.0
				//count := 4
				//for i:=0; i<count; i++ {
				//	mcount += math.Abs(lines[len(lines)-1-delta-i].Bar-lines[len(lines)-1-delta-i-1].Bar)
				//}
				//deltaBar := mcount/float64(count)
				//condition1 := last1.Dif > deltaBar && last1.Bar>deltaBar && last1.IsRed() && last1.Bar>last2.Bar && last2.Bar>0 && last5.Dif<0
				//dif突破0轴线
				//condition1 := (last1.Dif>0 && last2.Dif<0) || (last1.Dif>0 && last2.Dif>0 && lines[len(lines)-2].Dif<0)
				if isRecent && len(lines) > 100 && condition1 {
					if last1.ClosingPrice < lines[len(lines)-1-deltaDays].ClosingPrice {
						winCount += 1
						winStocks = append(winStocks, stock)
					} else {
						failedStocks = append(failedStocks, stock)
					}
					utils.JJKPrintln("code:",stock.Code," dif:", last1.Dif," dea:", last1.Dea," kdj-k:", last1.Kdj_k," kdj-d:", last1.Kdj_d)
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
		if len(winStocks) > 0 {
			//utils.JJKPrintln("======================")
			//for _, value := range winStocks {
			//	utils.JJKPrintln(value.Code, value.ChangeHandRate)
			//}
			utils.JJKPrintln("胜率：",float64(winCount)/float64(len(stocks)), len(stocks), winCount)
		}
		writeStocks(stocks)
		//writeStocks(winStocks)
		//writeStocks(failedStocks)
	}

	//量价突破
	if false {
		shStocks := database.DB.GetStockWithCodePrefix("60")
		//shStocks = append(shStocks, database.DB.GetStockWithCodePrefix("00")...)
		var stocks []*models.Stock
		for _, stock := range shStocks {
			_, lines := calculateMACD(stock.Code, 400)
			lines = calculateKDJWithLines(lines)
			if len(lines) > 40 {
				last1 := lines[len(lines)-5]
				last2 := lines[len(lines)-6]
				isRecent := last1.Date.After(time.Now().Add(time.Hour*24*-(7+5)))
				condition1 := false

				//量价同时突破
				//if true {
				//	daysCount := 20
				//	maxPrice := 0.0
				//	maxVol := 0.0
				//	for i:=3; i<daysCount; i++ {
				//		tmp := lines[len(lines)-i]
				//		if maxPrice < tmp.ClosingPrice {
				//			maxPrice = tmp.ClosingPrice
				//		}
				//		if maxVol < tmp.Vol {
				//			maxVol = tmp.Vol
				//		}
				//	}
				//
				//	rate := last1.GetAddRate(last2)
				//	if last1.Vol > maxVol && last1.ClosingPrice > maxPrice && last1.IsRed() && rate>0 && rate <=0.05 {
				//		condition1 = true
				//	}
				//}

				//量先突破
				if true {
					//daysCount := 4
					//for i:=2; i<=daysCount; i++ {
					//	if lines[len(lines)-i].Vol < lines[len(lines)-i-1].Vol {
					//		condition1 = false
					//	}
					//}
					sortVal := (last1.Vol-last2.Vol)/last2.Vol
					stock.SortVal = sortVal
					if last1.Vol > last2.Vol && last1.IsRed() && last1.ClosingPrice>last2.ClosingPrice {
						condition1 = true
					}
				}
				if isRecent && len(lines) > 100 && condition1 {
					utils.JJKPrintln("code:",stock.Code," dif:", last1.Dif," dea:", last1.Dea," kdj-k:", last1.Kdj_k," kdj-d:", last1.Kdj_d)
					stocks = append(stocks, stock)
				}

			}
		}

		//按换手率活跃度排序
		sort.Slice(stocks, func(i, j int) bool {
			return stocks[i].SortVal < stocks[j].SortVal
		})
		for _, value := range stocks {
			utils.JJKPrintln(value.Code, value.ChangeHandRate, value.SortVal)
		}
	}

	//找下引线越长的股票
	if false {
		//下载换手率
		InitCache()
		//DownloadTaskAddKLines()

		days := 1
		shStocks := database.DB.GetStockWithCodePrefix("60")
		shStocks = append(shStocks, database.DB.GetStockWithCodePrefix("00")...)
		var stocks []*models.Stock
		var winStocks []*models.Stock
		for _, stock := range shStocks {
			_, lines := calculateMACD(stock.Code, 400)
			lines = calculateKDJWithLines(lines)
			if len(lines) > 40 {
				nextDay := lines[len(lines)-days]
				if days != 1 {
					nextDay = lines[len(lines)-days+2]
				}
				last1 := lines[len(lines)-days]
				last2 := lines[len(lines)-days-1]
				isRecent :=  last1.Date.After(time.Now().Add(time.Hour*24*-(7)))
				lineDownRate := 0.0
				//priceDownAll := 0.0
				if last1.OpeningPrice > last1.ClosingPrice {
					lineDownRate = (last1.ClosingPrice-last1.MinPrice)/last2.ClosingPrice
					//priceDownAll = last1.ClosingPrice-last1.MinPrice
				} else {
					lineDownRate = (last1.OpeningPrice-last1.MinPrice)/last2.ClosingPrice
					//priceDownAll = last1.OpeningPrice-last1.MinPrice
				}
				rate := last1.GetAddRate(last2)
				//rate1 := priceDownAll/math.Abs(last1.ClosingPrice-last1.OpeningPrice)
				condition1 := rate<0  && lineDownRate>0.02 && last1.IsRed()
				if isRecent && len(lines) > 100 && condition1 {
					if days != 1 {
						if nextDay.ClosingPrice > last1.ClosingPrice {
							winStocks = append(winStocks, stock)
						}
					}
					stock.SortVal = lineDownRate
					utils.JJKPrintln("code:",stock.Code," dif:", last1.Dif," dea:", last1.Dea," kdj-k:", last1.Kdj_k," kdj-d:", last1.Kdj_d)
					stocks = append(stocks, stock)
				}

			}
		}

		//按换手率活跃度排序
		sort.Slice(stocks, func(i, j int) bool {
			return stocks[i].SortVal < stocks[j].SortVal
		})
		for _, value := range stocks {
			utils.JJKPrintln(value.Code, value.ChangeHandRate, value.SortVal)
		}
		if days != 1 {
			utils.JJKPrintln("胜率:", float64(len(winStocks))/float64(len(stocks)), len(stocks), len(winStocks))
		}
	}

	//KDJ选择下叉方式
	if false {
		shStocks := database.DB.GetStockWithCodePrefix("60")
		shStocks = append(shStocks, database.DB.GetStockWithCodePrefix("00")...)
		var stocks []*models.Stock
		for _, stock := range shStocks {
			_, lines := calculateKDJ(stock.Code, 400)
			lines = calculateMACDWithLines(lines)
			if len(lines) > 20 {
				delta := 0
				curDayStock := lines[len(lines)-1-delta]
				yestodayStock := lines[len(lines)-2-delta]
				condition1 := curDayStock.Dif>0 && curDayStock.Kdj_k < 45 && curDayStock.Kdj_k>curDayStock.Kdj_d && yestodayStock.Kdj_k < yestodayStock.Kdj_d
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
			_, lines := calculateKDJ(stock.Code,400)
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
			_, lines := calculateMACD(stock.Code,400)
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

	//测试下引线越长胜率
	if false {
		codes := strings.Split("600887", "|")
		for _, code := range codes {
			stock := database.DB.GetStockWithCode(code)
			lines := GetStockDayKLine(stock.CodeStr(), 400)
			length := len(lines)
			utils.JJKPrintln(length)
			allCount := 0.0
			winCount := 0.0
			if len(lines) > 40 {
				for i:=3; i<len(lines)-5 ; i++ {
					days := i
					nextDay := lines[len(lines)-days+2]
					last1 := lines[len(lines)-days]
					last2 := lines[len(lines)-days-1]
					lineDownRate := 0.0
					//priceDownAll := 0.0
					if last1.OpeningPrice > last1.ClosingPrice {
						lineDownRate = (last1.ClosingPrice-last1.MinPrice)/last2.ClosingPrice
						//priceDownAll = last1.ClosingPrice-last1.MinPrice
					} else {
						lineDownRate = (last1.OpeningPrice-last1.MinPrice)/last2.ClosingPrice
						//priceDownAll = last1.OpeningPrice-last1.MinPrice
					}
					//rate1 := priceDownAll/math.Abs(last1.ClosingPrice-last1.OpeningPrice)
					condition1 := !last1.IsRed() && lineDownRate>0.04
					if condition1 {
						allCount += 1
						utils.JJKPrintln(last1.Date, lineDownRate)
						if nextDay.ClosingPrice > last1.ClosingPrice {
							winCount += 1
						} else {

						}
					}
				}
			}
			utils.JJKPrintln("胜率：",winCount/allCount, winCount, allCount)
		}
	}

	//量比排序
	if false {
		InitCache()
		downloadStockInfo()
		InitCache()

		stocks := CCGetStockAll()
		sort.Slice(stocks, func(i, j int) bool {
			return stocks[i].VolRate < stocks[j].VolRate
		})
		for _, value := range stocks {
			utils.JJKPrintln(value.Code, value.VolRate)
		}
	}

	breakThroughMA13WithMADown()
	diWeiZhangTing()
	junXianXiangShang()
}

func GetObserverStocksFilePath() string {
	return path.Join(utils.ExeDir(), "conf/stocks.txt")
}

func observeStocks()  {
	if true {
		mystocks := database.DB.GetMyStocks()
		//是否可以买入监控
		if true {
			//macd黄金交叉买入
			for _, value := range mystocks {
				code := value.Code
				stock := database.DB.GetStockWithCode(code)
				klines := GetStockDayKLine(stock.CodeStr(),400)
				lines := calculateKDJWithLines(klines)
				lines = calculateMACDWithLines(lines)
				if len(lines) > 2 {
					last1 := lines[len(lines)-1]
					last2 := lines[len(lines)-2]
					isRecent := last1.Date.After(time.Now().Add(time.Hour*24*-7))
					//kdj低位，黄金交叉
					condition1 := last1.Kdj_k>last1.Kdj_d && last2.Kdj_k<last2.Kdj_d && last1.Kdj_k < 50.0
					if isRecent && condition1 {
						utils.JJKPrintln(stock.Code, last1.Rsi, last1.ClosingPrice,  stock.FlowAmount)
						PushNotification(fmt.Sprintf("【买入提醒】：股票代码:%s KDJ发生黄金交叉，可以择机买入！", stock.Code))
					}
				}
			}
		}

		//卖出监控
		if true {
			//kdj死亡交叉或macd变绿
			for _, value := range mystocks {
				if value.IsBuy {
					code := value.Code
					stock := database.DB.GetStockWithCode(code)
					klines := GetStockDayKLine(stock.CodeStr(),400)
					lines := calculateKDJWithLines(klines)
					lines = calculateMACDWithLines(lines)
					if len(lines) > 2 {
						last1 := lines[len(lines)-1]
						last2 := lines[len(lines)-2]
						isRecent := last1.Date.After(time.Now().Add(time.Hour*24*-7))

						//kdj死叉判断
						if isRecent && last1.Kdj_k<last1.Kdj_d && last2.Kdj_k>last2.Kdj_d {
							utils.JJKPrintln(stock.Code, last1.Rsi, last1.ClosingPrice,  stock.FlowAmount)
							PushNotification(fmt.Sprintf("【卖出提醒】：%s:%s 发生KDJ死亡交叉，应该减仓！", value.Name, value.Code))
						}

						//macd由红变绿
						if isRecent && last1.Bar<0 && last2.Bar>0 {
							utils.JJKPrintln(stock.Code, last1.Rsi, last1.ClosingPrice,  stock.FlowAmount)
							PushNotification(fmt.Sprintf("【卖出提醒】：%s:%s macd柱状由红柱变为绿柱，应全部清仓！", value.Name, stock.Code))
						}

						//价格低于买入%5，清仓一半
						rate := (last1.ClosingPrice - value.BuyPrice)/value.BuyPrice
						utils.JJKPrintln(rate)
						if isRecent && rate < -0.05 && rate > -0.1 {
							utils.JJKPrintln(stock.Code, last1.Rsi, last1.ClosingPrice,  stock.FlowAmount)
							PushNotification(fmt.Sprintf("【卖出提醒】：%s:%s 跌幅超过买入价%%5，应减仓一半！",value.Name, stock.Code))
						}

						if isRecent && rate < -0.1 {
							utils.JJKPrintln(stock.Code, last1.Rsi, last1.ClosingPrice,  stock.FlowAmount)
							PushNotification(fmt.Sprintf("【卖出提醒】：%s:%s 跌幅超过买入价%%10，应全部清仓！",value.Name, stock.Code))
						}
					}
				}
			}
		}
	}


	//macd统计方式
	if false {
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
}

func writeStocks(stocks []*models.Stock)  {
	file := GetObserverStocksFilePath()
	text := ""
	for _, value := range stocks {
		text = fmt.Sprintf("%s%s %.2f\n", text,value.Code, value.ChangeHandRate)
	}

	utils.FileWriteString(file, text)
}

func sendSMS(stockCode string)  {
	mobile := "17602125152"
	if err := dysms.SendSms("LTAIwsROTsyzMq1a", "KIJzyH67zdskMBdjOIYKnDdv7xavX7", mobile, "爱编程", fmt.Sprintf(`{"code":"%s"}`, stockCode), "SMS_128535187"); err != nil {
		utils.JJKPrintln("dysms.SendSms", err)
	}
}

func breakThroughMA13WithMADown()  {
	//均线突破选股
	if true {
		//下载换手率等
		//downloadStockInfo()

		stocks := CCGetStockAll()
		var resultStocks []*models.Stock
		for _, stock := range stocks {
			lines := CCGetKLinesWithCode(stock.Code, 100)
			if len(lines) > 40 {
				//计算13日均线
				lines = caculateMA(lines, 13)

				isRecent := lines[len(lines)-1].Date.After(time.Now().Add(time.Hour*24*-5))
				startIndex := 4
				endIndex := 16
				startDay := lines[len(lines)-startIndex]
				endDay := lines[len(lines)-endIndex]
				//均线向下
				maIsDown := startDay.MaVal < endDay.MaVal
				//均线斜率
				maRate := (endDay.MaVal-startDay.MaVal)/startDay.MaVal
				//突破天数
				days := 1
				bigCondition := true
				littleCondition := true
				for i:=1; i<=days+3; i++ {
					if i<=days {
						if lines[len(lines)-i].ClosingPrice < lines[len(lines)-i].MaVal {
							bigCondition = false
						}
					} else {
						if lines[len(lines)-i].ClosingPrice > lines[len(lines)-i].MaVal {
							littleCondition = false
						}
					}
				}
				//均线突破
				if isRecent && maIsDown && bigCondition && littleCondition && maRate < 0.2 && startDay.MaVal < 50.0 {
					utils.JJKPrintln(stock.Code, lines[len(lines)-1].MaVal, maRate)
					stock.SortVal = maRate
					resultStocks = append(resultStocks, stock)
				}
			}
		}

		//按换手率活跃度排序
		sort.Slice(resultStocks, func(i, j int) bool {
			return resultStocks[i].SortVal < resultStocks[j].SortVal
		})
		for _, value := range resultStocks {
			utils.JJKPrintln(value.Code, " 跌破点数：",value.SortVal*100)
		}

		utils.JJKPrintln("筛选结果共：", len(resultStocks))
	}
}

func diWeiZhangTing()  {
	//找低位涨停，缩量回调，不破平台
	if false {
		//下载换手率
		//DownloadTaskAddKLines()
		InitCache()


		shStocks := database.DB.GetStockWithCodePrefix("60")
		shStocks = append(shStocks, database.DB.GetStockWithCodePrefix("00")...)
		var stocks []*models.Stock
		for _, stock := range shStocks {
			_, lines := calculateMACD(stock.Code, 400)
			lines = calculateKDJWithLines(lines)
			if len(lines) > 40 {
				//6~10天之间发生过涨停
				startDays := 4
				recentDays := 6
				var recentKline *models.KLine
				recentKline = lines[len(lines)-recentDays]
				isRecent :=  lines[len(lines)-1].Date.After(time.Now().Add(time.Hour*24*-(7)))
				condition1 := false
				count := 0
				for i:=startDays; i<= recentDays; i++ {
					currentDay := lines[len(lines)-i]
					prevDay := lines[len(lines)-i-1]
					rate := currentDay.GetAddRate(prevDay)
					if rate > 0.095 && currentDay.Dif < 0 {
						condition1 = true
						count += 1
					}
				}
				if isRecent && len(lines) > 100 && condition1 && count == 1 {
					stock.SortVal = stock.ChangeHandRate
					utils.JJKPrintln("code:",stock.Code," price:", recentKline.ClosingPrice)
					stocks = append(stocks, stock)
				}

			}
		}

		//按换手率活跃度排序
		sort.Slice(stocks, func(i, j int) bool {
			return stocks[i].SortVal < stocks[j].SortVal
		})
		for _, value := range stocks {
			utils.JJKPrintln("低位涨停：",value.Code, value.ChangeHandRate, value.SortVal)
		}
	}
}

//均线向上的股票
func junXianXiangShang()  {
	//均线突破选股
	if false {
		//下载换手率等
		//downloadStockInfo()

		stocks := CCGetStockAll()
		var resultStocks []*models.Stock
		for _, stock := range stocks {
			lines := CCGetKLinesWithCode(stock.Code, 100)
			if len(lines) > 40 {
				//计算13日均线
				lines = caculateMA(lines, 13)

				isRecent := lines[len(lines)-1].Date.After(time.Now().Add(time.Hour*24*-5))
				startIndex := 2
				endIndex := 6
				index := 14
				startDay := lines[len(lines)-startIndex]
				endDay := lines[len(lines)-endIndex]
				indexDay := lines[len(lines)-index]
				//均线向上
				maIsUP := startDay.MaVal > endDay.MaVal
				//均线斜率
				maRate := (startDay.MaVal-endDay.MaVal)/startDay.MaVal
				//均线突破
				if isRecent && maIsUP && startDay.MaVal < 50.0 && indexDay.MaVal < endDay.MaVal {
					utils.JJKPrintln(stock.Code, lines[len(lines)-1].MaVal, maRate)
					stock.SortVal = maRate
					resultStocks = append(resultStocks, stock)
				}
			}
		}

		//按换手率活跃度排序
		sort.Slice(resultStocks, func(i, j int) bool {
			return resultStocks[i].SortVal < resultStocks[j].SortVal
		})
		for _, value := range resultStocks {
			utils.JJKPrintln(value.Code, " 涨的点数：",value.SortVal*100)
		}

		utils.JJKPrintln("筛选结果共：", len(resultStocks))
	}
}