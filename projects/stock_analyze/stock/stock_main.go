package stock

import (
	"github.com/cheneylew/goutil/utils"
	"github.com/cheneylew/goutil/projects/stock_analyze/database"
	"time"
	"github.com/cheneylew/goutil/projects/stock_analyze/models"
	"fmt"
	"sort"
	"math"
	"strings"
)

var serverKLines []*models.KLine

func ClearCache()  {
	serverKLines = serverKLines[:0]
}


func  StockTestMain()  {
	CronMain()						//定时任务
	InitCache()
	Main_rsi()
	//DownloadTaskAddKLines()		//同步增量数据

	//uploadStocksCodeToDB()		//同步所有股票代码到数据库
	//downloadSHStockKLines()		//下载上证所有股票日K
	//downloadSZStockKLines()		//下载深证所有股票日K
	//downloadFaildStocks()			//下载失败的股票日K
	//downloadStockRealTimeInfo()	//五日增减仓数据
	//downloadStockInfo()			//下载股票信息，总市值等

	//AnalysMACD()
	//CalculateMACD()				//macd计算
	//ReCalculateMACD()				//根据已经计算过的macd，再次计算新增的
	//AnalysBuyWhat()				//分析买什么
	//Analys5MainInStocks()			//五日增仓
	//AnalysUpStock()				//分析走势向上的股票

	//date := database.DB.GetDateOfLastKLineWithStockID(588)
	//utils.JJKPrintln(date)

	//downloadDayKLine("600000")	//单独下载某只股票
	utils.JJKPrintln("股票分析结束, 启动web服务！")
}

func DownloadTaskAll()  {
	InitCache()

	uploadStocksCodeToDB()
	downloadSHStockKLines()		//下载上证所有股票日K
	downloadSZStockKLines()		//下载深证所有股票日K
	////downloadFaildStocks()			//下载失败的股票日K
	downloadStockRealTimeInfo()	//五日增减仓数据
	downloadStockInfo()			//下载股票信息，总市值等

	CalculateMACD()
}

func DownloadTaskAddKLines()  {
	InitCache()

	//uploadStocksCodeToDB()
	downloadSHStockKLines()		//下载上证所有股票日K
	downloadSZStockKLines()		//下载深证所有股票日K
	////downloadFaildStocks()			//下载失败的股票日K
	downloadStockRealTimeInfo()	//五日增减仓数据
	downloadStockInfo()			//下载股票信息，总市值等

	ReCalculateMACD()
}

func uploadStocksCodeToDB()  {
	UpdateOnlineCodesToDatabase()
}

// 增量下载股票的日K
func downloadDayKLine(code string)  {
	if len(serverKLines) == 0 {
		serverKLines = database.DB.GetKLineAll()
	}

	tmpStock := database.DB.GetStockWithCode(code)
	//增量计算,防止重复下载同步
	lastSyncTime := database.DB.GetDateOfLastKLineWithStockID(tmpStock.StockId)
	deltaSecnods := time.Now().Unix() - lastSyncTime.Unix()

	deltaDays := float64(deltaSecnods)/float64(3600*24)
	utils.JJKPrintln(deltaDays)
	var defaultDays float64 = 100
	if deltaDays > defaultDays {	// > 100天
		deltaDays = defaultDays
	} else if deltaDays >= 1 && deltaDays <= defaultDays {// [1,100]天
		deltaDays = math.Ceil(deltaDays)
	} else { // < 1天
		//是否为今天
		isDateEqual := utils.DateEqual(lastSyncTime, time.Now())
		if isDateEqual { //当天需要再次同步一下
			deltaDays = 1
		} else {	//非当天同步一次即可
			utils.JJKPrintln(fmt.Sprintf("%s have sync!", code))
			return
		}
	}

	//某只股票的K线
	tmpStockKLines := make([]*models.KLine, 0)
	for _, value := range serverKLines {
		if value.StockId == tmpStock.StockId {
			tmpStockKLines = append(tmpStockKLines, value)
		}
	}

	klines := GetStockDayKLine(tmpStock.CodeStr(),int64(deltaDays))
	if len(klines) == 0 {
		//tmpStock.SyncTime = time.Now()
		tmpStock.SyncOk = false
		database.DB.Orm.Update(tmpStock)
		utils.JJKPrintln("download failed ", code)
		return ;
	}
	for _, kline := range klines {
		kline.StockId = tmpStock.StockId

		isExist := false
		for _, sKLine := range tmpStockKLines {
			if sKLine.StockId == kline.StockId && sKLine.Date == kline.Date {
				isExist = true
				kline.KLineId = sKLine.KLineId
			}
		}

		if !isExist {
			//不存在，插入
			database.DB.Orm.Insert(kline)
		} else {
			y,m,d := kline.Date.Date()
			ty,tm,td := time.Now().Date()
			if y == ty && m == tm && d == td {
				//今天的股票 需要update
				a, err := database.DB.Orm.Update(kline)
				if err != nil {
					utils.JJKPrintln("update kline today failed!", err)
				} else {
					utils.JJKPrintln("update kline today ok!", a)
				}
			}
		}

	}

	tmpStock.SyncTime = time.Now()
	tmpStock.SyncOk = true
	database.DB.Orm.Update(tmpStock)
	utils.JJKPrintln("download kline ok ",code)

}

func downloadSHStockKLines()  {
	ClearCache()
	shStocks := database.DB.GetStockWithCodePrefix("60")
	var params []interface{}
	for _, value := range shStocks {
		params = append(params, value)
	}
	utils.QueueTask(10, params, func(idx int, param interface{}) {
		stock := param.(*models.Stock)
		downloadDayKLine(stock.Code)
	})
	utils.JJKPrintln(len(shStocks))
}

func downloadSZStockKLines()  {
	ClearCache()
	shStocks := database.DB.GetStockWithCodePrefix("00")
	var params []interface{}
	for _, value := range shStocks {
		params = append(params, value)
	}
	utils.QueueTask(10, params, func(idx int, param interface{}) {
		stock := param.(*models.Stock)
		downloadDayKLine(stock.Code)
	})
	utils.JJKPrintln(len(shStocks))
}


func downloadFaildStocks()  {
	ClearCache()
	shStocks := database.DB.GetSyncFailedStocks()

	var params []interface{}
	for _, value := range shStocks {
		params = append(params, value)
	}
	utils.QueueTask(10, params, func(idx int, param interface{}) {
		stock := param.(*models.Stock)
		downloadDayKLine(stock.Code)
	})

	fstocks := database.DB.GetSyncFailedStocks()
	utils.JJKPrintln(len(fstocks))
}

func downloadStockRealTimeInfo()  {
	stks := allStocks
	utils.JJKPrintln(len(stks))

	var params []interface{}
	for _, value := range stks {
		params = append(params, value)
	}

	utils.QueueTask(10, params, func(idx int, param interface{}) {
		tmpStk := param.(*models.Stock)

		onlineInfos := GetRealTimeStockInfo(tmpStk.CodeStr())
		dbInfos := CCGetStockInfoWithStockId(tmpStk.StockId)
		if len(onlineInfos) == 0 {
			utils.JJKPrintln(fmt.Sprintf("%s failed", tmpStk.Code))
		} else {
			for _, value := range onlineInfos {
				//是否已存在
				isExist := false
				var existObj *models.StockInfo
				for _, dbInfo := range dbInfos {
					if utils.DateEqual(dbInfo.Date, value.Date) {
						isExist = true
						existObj = value
						existObj.StockInfoId = dbInfo.StockInfoId
						existObj.StockId = dbInfo.StockId
					}
				}
				if !isExist {
					value.StockId = tmpStk.StockId
					_,e := database.DB.Orm.Insert(value)
					if e != nil {
						utils.JJKPrintln(e)
					} else {
						utils.JJKPrintln(fmt.Sprintf("%s %s insert ok!", tmpStk.Code, value.Date))
					}
				} else {
					utils.JJKPrintln(fmt.Sprintf("%s %s have existed", tmpStk.Code, value.Date))
					if utils.DateEqual(value.Date, time.Now()) {
						utils.JJKPrintln("today update!")
						a, err := database.DB.Orm.Update(existObj)
						if err != nil {
							utils.JJKPrintln("%s %s update main today failed!", tmpStk.Code, value.Date, err)
						} else {
							utils.JJKPrintln("%s %s update main today ok!", tmpStk.Code, value.Date, a)
						}
					}
				}
			}
		}
	})

	InitCache()
}

func downloadStockInfo()  {
	stocks := CCGetStockAll()
	var params []interface{}
	for _, value := range stocks {
		params = append(params, value)
	}
	utils.QueueTask(10, params, func(idx int, param interface{}) {
		stock := param.(*models.Stock)

		info := utils.HTTPGet(fmt.Sprintf("http://qt.gtimg.cn/q=%s", stock.CodeStr()))
		arr := strings.Split(info, "=")
		if len(arr) == 2 {
			str := utils.TrimChars(utils.TrimWiteSpace(arr[1]),`;"`)
			infoArr := strings.Split(str, "~")

			if len(infoArr) <= 48 {
				return
			}

			deltaMoney := utils.ToFloat64(infoArr[31]) 		//涨跌
			deltaMonyRate := utils.ToFloat64(infoArr[32]) 	//涨跌%  

			flowAmount := utils.ToFloat64(infoArr[44]) 		//流通市值
			totalAmount := utils.ToFloat64(infoArr[45]) 	//总市值
			changeHandRate := utils.ToFloat64(infoArr[38]) 	//换手率
			PERate := utils.ToFloat64(infoArr[39]) 			//市盈率 
			PBRate := utils.ToFloat64(infoArr[46]) 			//市净率 
			volAmount := utils.ToFloat64(infoArr[36]) 		//成交量（手）  
			volAmountMoney := utils.ToFloat64(infoArr[37]) 	//成交额（万)
			VolRate := utils.ToFloat64(infoArr[49]) 		//量比

			stock.FlowAmount = flowAmount
			stock.TotalAmount = totalAmount
			stock.ChangeHandRate = changeHandRate
			stock.PERate = PERate
			stock.PBRate = PBRate
			stock.VolAmount = volAmount
			stock.VolAmountMoney = volAmountMoney
			stock.DeltaMoney = deltaMoney
			stock.DeltaMoneyRate = deltaMonyRate
			stock.VolRate = VolRate


			database.DB.Orm.Update(stock)
			utils.JJKPrintln(fmt.Sprintf("%s change hand ok", stock.Code))
		}
	})
}

func AnalysResultWithCodeAndDays(code string, days int) *models.AnalysDayKLine {
	stock := database.DB.GetStockWithCode(code)
	result := &models.AnalysDayKLine {
		Stock:stock,
		Days:int64(days),
	}

	if len(serverKLines) == 0 {
		serverKLines = database.DB.GetKLineAll()
	}

	//lines := database.DB.GetKLineAllForStockCodeAndDays(code,days)
	date := time.Now().Add(-time.Hour * 24*time.Duration(days))
	var lines []*models.KLine
	for _, value := range serverKLines {
		if value.StockId == stock.StockId {
			if value.Date.After(date){
				lines = append(lines, value)
			}
		}
	}
	status := ""
	var lastKLine *models.KLine
	for _, value := range lines {
		rate := value.GetAddRate(lastKLine)
		if value.IsRed() {
			result.RedCount += 1
			status += fmt.Sprintf("red %.4f|",rate)
		} else {
			result.GreenCount += 1
			status += fmt.Sprintf("green %.4f|", rate)
		}
		if rate > 0 {
			result.UpCount += 1
			result.UpRateCount += rate
		} else {
			result.DownCount += 1
			result.DownRateCount += rate
		}
		lastKLine = value
	}

	result.UpDownRateTotal = result.UpRateCount+result.DownRateCount

	return result
}

func AnalysStockInfo() []*models.Stock {
	stocks := database.DB.GetStockWithCodePrefix("60")
	stocks = append(stocks, database.DB.GetStockWithCodePrefix("00")...)
	results := []*models.Stock{}
	for i:=0; i<len(stocks) ; i++ {
		infos := CCGetStockInfoWithStockId(stocks[i].StockId)
		if len(infos) > 0 {
			stocks[i].Infos = infos
			ok := false
			//主力减仓
			//if infos[0].MainTotal < 0 && infos[1].MainTotal < 0 && infos[2].MainTotal < 0 && infos[3].MainTotal < 0 && infos[4].MainTotal < 0 {
			if infos[0].MainTotal < 0 && infos[1].MainTotal < 0 && infos[2].MainTotal < 0  && infos[3].MainTotal > 0 {
			//主力增仓
			//if infos[0].MainTotal > 0 && infos[1].MainTotal > 0 && infos[2].MainTotal > 0 && infos[3].MainTotal > 0 && infos[4].MainTotal > 0 {
				ok = true
			}
			if ok {
				results = append(results, stocks[i])
			}
		}
	}

	results2 := []*models.Stock{}
	for _, value := range results {
		klines := database.DB.GetKLineAllForStock(value)
		if len(klines) > 2 {
			rate := klines[len(klines)-1].GetAddRate(klines[len(klines)-2])
			up,ucnt,dcnt := KLineIsUp(klines)
			utils.JJKPrintln(rate)
			if rate * 100 < -1 && up {
				results2 = append(results2, value)
				utils.JJKPrintln(value.Code,ucnt, dcnt)
			}
		}

	}

	return results2
}

func AnalysRedRate(days int) models.SortAnalysDayKLins {
	shStocks := database.DB.GetStockWithCodePrefix("60")
	var all models.SortAnalysDayKLins
	for _, value := range shStocks {
		result := AnalysResultWithCodeAndDays(value.Code,days)
		all = append(all, result)
	}
	sort.Sort(all)

	var back models.SortAnalysDayKLins
	for _, value := range all {
		if value.RedCount + value.GreenCount > 3 {
			back = append(back, value)
		}
	}

	return back
}

func AnalysUpStock()  {
	stocks := database.DB.GetStockWithCodePrefix("60")
	var ss []*models.Stock
	for _, value := range stocks {
		klines := database.DB.GetKLineAllForStockCode(value.Code)
		up,_,_ := KLineIsUp(klines)
		if up && klines[len(klines)-1].Date.After(time.Now().Add(-time.Hour*24*3)) {
			utils.JJKPrintln(fmt.Sprintf("%s ok", value.Code))
			ss = append(ss, value)
		}
	}

	for _, value := range ss {
		utils.JJKPrintln(value.Code)
	}
}

func AnalysBuyWhat() []*models.Stock {
	InitCache()

	var stocks []*models.Stock
	for _, stock := range CCGetStockAll() {
		kls := CCGetKLinesWithCode(stock.Code, 6)
		if len(kls) == 6 {
			klines := CCGetKLinesWithCode(stock.Code, 20)
			if len(klines) >= 20 {
				isUp,_,_ := KLineIsUp(klines[:len(klines)-5])
				if utils.DateEqual(kls[5].Date, time.Now()) && isUp {
					stock.DeltaVal = (kls[1].ClosingPrice - kls[5].ClosingPrice)/kls[5].ClosingPrice
					stock.DeltaVal = -stock.DeltaVal * 100
					stocks = append(stocks, stock)
				}
			}
		}
	}

	sort.Slice(stocks, func(i, j int) bool {
		return stocks[i].DeltaVal < stocks[j].DeltaVal
	})

	s := ""
	for _, value := range stocks {
		s += fmt.Sprintf("%s %f\n", value.Code, value.DeltaVal)
	}

	//utils.FileWriteString("/Users/apple/Desktop/a.txt", s)
	return stocks
}

func AnalysBuyUp() []*models.Stock {
	InitCache()

	recentCount := 10
	var stocks []*models.Stock
	for _, stock := range CCGetStockAll() {
		klines := CCGetKLinesWithCode(stock.Code, 30)
		if len(klines) >= 30 {
			isUp,_,_ := KLineIsUp(klines[:len(klines)-recentCount])
			if isUp {
				stocks = append(stocks, stock)
			}
		}
	}

	sort.Slice(stocks, func(i, j int) bool {
		return stocks[i].ChangeHandRate > stocks[j].ChangeHandRate
	})

	return stocks
}

func Analys5MainInStocks() []*models.Stock {
	var stocks []*models.Stock
	utils.JJKPrintln(len(allStocks))
	for _, stock := range allStocks {
		infos := CCGetStockInfoWithStockId(stock.StockId)
		sort.Slice(infos, func(i, j int) bool {
			return infos[i].Date.Before(infos[j].Date)
		})

		if len(infos) == 5 {
			ok := true
			for _, info := range infos {
				if info.MainTotal < 100 {
					ok = false
				}
			}

			if ok {
				stocks = append(stocks, stock)
			}
		}
	}

	return stocks
}

func Analys5MainInGreatThan3DaysStocks() []*models.Stock {
	var stocks []*models.Stock
	utils.JJKPrintln(len(allStocks))
	for _, stock := range allStocks {
		infos := CCGetStockInfoWithStockId(stock.StockId)
		sort.Slice(infos, func(i, j int) bool {
			return infos[i].Date.Before(infos[j].Date)
		})

		if len(infos) == 5 {
			ok := false
			addConnt := 0
			for _, info := range infos {
				if info.MainTotal > 100 {
					addConnt += 1
				}
			}

			if addConnt >= 3 {
				ok = true
			}

			if ok {
				stocks = append(stocks, stock)
			}
		}
	}

	return stocks
}

func Analys5MainInGreatThan2DaysStocks() []*models.Stock {
	var stocks []*models.Stock
	utils.JJKPrintln(len(allStocks))
	for _, stock := range allStocks {
		infos := CCGetStockInfoWithStockId(stock.StockId)
		sort.Slice(infos, func(i, j int) bool {
			return infos[i].Date.Before(infos[j].Date)
		})

		if len(infos) == 5 {
			ok := false
			addConnt := 0
			for _, info := range infos {
				if info.MainTotal > 100 {
					addConnt += 1
				}
			}

			if addConnt >= 2 {
				ok = true
			}

			if ok {
				stocks = append(stocks, stock)
			}
		}
	}

	return stocks
}

func Analys5MainInEqual2DaysStocks() []*models.Stock {
	var stocks []*models.Stock
	utils.JJKPrintln(len(allStocks))
	for _, stock := range allStocks {
		infos := CCGetStockInfoWithStockId(stock.StockId)
		sort.Slice(infos, func(i, j int) bool {
			return infos[i].Date.Before(infos[j].Date)
		})

		if len(infos) == 5 {
			ok := false
			addConnt := 0
			for _, info := range infos {
				if info.MainTotal > 100 {
					addConnt += 1
				}
			}

			if addConnt == 2 {
				ok = true
			}

			if ok {
				stocks = append(stocks, stock)
			}
		}
	}

	return stocks
}

func Analys5MainOutStocks() []*models.Stock {
	var stocks []*models.Stock
	utils.JJKPrintln(len(allStocks))
	for _, stock := range allStocks {
		infos := CCGetStockInfoWithStockId(stock.StockId)
		sort.Slice(infos, func(i, j int) bool {
			return infos[i].Date.Before(infos[j].Date)
		})

		if len(infos) == 5 {
			ok := true
			for _, info := range infos {
				if info.MainTotal > -100 {
					ok = false
				}
			}

			klines := CCGetKLinesWithCode(stock.Code,40)
			up,_,_ := KLineIsUp(klines)
			if ok && up {
				stocks = append(stocks, stock)
			}
		}
	}

	return stocks
}

func AnalysMACD(redCount int) []*models.Stock {
	var resStocks []*models.Stock
	stocks := CCGetStockAll()
	for _, stock := range stocks {
		count := 6
		redDays := redCount
		klines := CCGetKLinesWithCode(stock.Code, count)
		if len(klines) == count {
			if klines[0].Ema12 != 0 {
				leftOk := true
				rightOk := true
				greenBarCount := 0.0
				redBarCount := 0.0
				for i:=0; i< len(klines); i++ {
					if i>=(count-redDays) {
						if klines[i].Bar < 0 {
							rightOk = false
						} else {
							redBarCount += klines[i].Bar
						}
					} else {
						if klines[i].Bar > 0 {
							leftOk = false
						} else {
							greenBarCount += klines[i].Bar
						}
					}
				}
				if redDays > 0 {
					if leftOk && rightOk {
						stock.RedBarCount = redBarCount
						stock.GreenBarCount = greenBarCount

						utils.JJKPrintln(stock.Code, klines[0].Bar)
						resStocks = append(resStocks, stock)
					}
				} else {
					abs := math.Abs(klines[len(klines)-1].Bar -  klines[len(klines)-2].Bar)
					lastOne := klines[len(klines)-1]
					lastOneOk := lastOne.Bar >= - 100
					max := 0.02
					if lastOne.Bar >= 1 {
						max = math.Abs(lastOne.Bar)*0.01
					}
					if lastOneOk && leftOk && abs <= max {
						stock.RedBarCount = redBarCount
						stock.GreenBarCount = greenBarCount

						resStocks = append(resStocks, stock)
					}
				}
			}
		}
	}

	//5日增仓情况筛选,至少有两天增仓
	var filterStocks []*models.Stock
	for _, stock := range resStocks {
		var stockinfos []*models.StockInfo
		database.DB.Orm.Raw("SELECT * FROM stock.stock_info where stock_id = ? order by date asc limit 0, 5;", stock.StockId).QueryRows(&stockinfos)
		addCount := 0
		for _, info := range stockinfos {
			if info.MainTotal > 100 {
				addCount += 1
			}
		}
		//5天中有2天增仓
		if addCount >= 1 {
			filterStocks = append(filterStocks, stock)
		}
	}

	return filterStocks
}

func AnalysNewStocks() []*models.Stock  {

	lessThanDays := 20
	var resStocks []*models.Stock
	stocks := CCGetStockAll()
	for _, stock := range stocks {
		klines := CCGetKLinesWithCode(stock.Code,lessThanDays+1)
		if len(klines) < lessThanDays && len(klines) != 0 {
			resStocks = append(resStocks, stock)
		}
	}

	return resStocks
}