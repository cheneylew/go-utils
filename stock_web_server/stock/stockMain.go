package stock

import (
	"github.com/cheneylew/goutil/utils"
	"github.com/cheneylew/goutil/stock_web_server/database"
	"time"
	"github.com/cheneylew/goutil/stock_web_server/models"
	"fmt"
	"sort"
	"math"
)

var serverKLines []*models.KLine

func ClearCache()  {
	serverKLines = serverKLines[:0]
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
	lastSyncTime := tmpStock.SyncTime
	secnods := time.Now().Add(-time.Hour*8).Unix() - lastSyncTime.Unix()
	ly,lm,ld := lastSyncTime.Date()
	y,m,d := time.Now().Date()
	isSameDate := false
	if ly == y && lm == m && ld == d {
		isSameDate = true
	}

	days := float64(secnods)/float64(3600*24)
	if days > 100 || !tmpStock.SyncOk {
		days = 100
	} else if days >= 1 && days <= 100 {
		days = math.Ceil(days)
	} else {
		if !isSameDate {
			days = 1
		} else {
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

	klines := GetStockDayKLine(tmpStock.CodeStr(),int64(days))
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
			if sKLine.StockId == kline.StockId && sKLine.Date == kline.Date.Add(-time.Hour * 8) {
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
	utils.JJKPrintln("download ok ",code)

}

func downloadSHStockKLines()  {
	ClearCache()
	shStocks := database.DB.GetStockWithCodePrefix("60")
	for _, value := range shStocks {
		downloadDayKLine(value.Code)
	}
	utils.JJKPrintln(len(shStocks))
}

func downloadSZStockKLines()  {
	ClearCache()
	shStocks := database.DB.GetStockWithCodePrefix("00")
	for _, value := range shStocks {
		downloadDayKLine(value.Code)
	}
	utils.JJKPrintln(len(shStocks))
}


func downloadFaildStocks()  {
	ClearCache()
	for _, value := range database.DB.GetSyncFailedStocks() {
		downloadDayKLine(value.Code)
	}

	fstocks := database.DB.GetSyncFailedStocks()
	utils.JJKPrintln(len(fstocks))
}

func downloadStockRealTimeInfo()  {
	stks := allStocks
	utils.JJKPrintln(len(stks))
	for _, tmpStk := range stks {
		onlineInfos := GetRealTimeStockInfo(tmpStk.CodeStr())
		dbInfos := CCGetStockInfoWithStockId(tmpStk.StockId)
		if len(onlineInfos) == 0 {
			utils.JJKPrintln(fmt.Sprintf("%s failed", tmpStk.Code))
		} else {
			utils.JJKPrintln(onlineInfos)
			for _, value := range onlineInfos {
				//是否已存在
				isExist := false
				var existObj *models.StockInfo
				for _, dbInfo := range dbInfos {
					if utils.DateEqual(dbInfo.Date.Add(time.Hour*8), value.Date.Add(time.Hour*8)) {
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
						utils.JJKPrintln("insert ok!")
					}
				} else {
					utils.JJKPrintln("have existed")
					if utils.DateEqual(value.Date, time.Now()) {
						utils.JJKPrintln("today update!")
						a, err := database.DB.Orm.Update(existObj)
						if err != nil {
							utils.JJKPrintln("update main today failed!", err)
						} else {
							utils.JJKPrintln("update main today ok!", a)
						}
					}
				}
			}
		}
	}

	InitCache()
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

func AnalysBuyWhat() []*models.Stock {
	InitCache()
	//for _, stock := range CCGetStockAll() {
	//	kls := CCGetKLinesWithCode(stock.Code, 5)
	//	if len(kls) == 5 {
	//		if !kls[4].Date.Before(time.Now().Add(-time.Hour*24*3)) {
	//			if kls[0].IsRed() && kls[1].IsRed() &&  !kls[2].IsRed() &&  !kls[3].IsRed() &&  !kls[4].IsRed() {
	//				if kls[4].ClosingPrice < kls[0].ClosingPrice {
	//					rate := kls[4].GetAddRate(kls[3])
	//					if rate < -0.02 {
	//						utils.JJKPrintln(stock.Code)
	//					}
	//				}
	//			}
	//		}
	//	}
	//}

	var stocks []*models.Stock
	for _, stock := range CCGetStockAll() {
		kls := CCGetKLinesWithCode(stock.Code, 6)
		if len(kls) == 6 {
			klines := CCGetKLinesWithCode(stock.Code, 20)
			if len(klines) >= 20 {
				isUp,_,_ := KLineIsUp(klines[:len(klines)-5])
				if utils.DateEqual(kls[5].Date.Add(time.Hour*8), time.Now()) && isUp {
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

func StockTestMain()  {
	CronMain()
	InitCache()
	//AnalysBuyWhat()

	//uploadStocksCodeToDB()
	//utils.JJKPrintln(len(database.DB.GetKLineAll()))
	//downloadSHStockKLines()
	//downloadSZStockKLines()
	//downloadFaildStocks()
	//downloadDayKLine("600196")
	//stock.GetRealTimeStockInfo("sh600703")
	//downloadStockRealTimeInfo()
	
	//stocks := database.DB.GetStockWithCodePrefix("00")
	//var ss []*models.Stock
	//for _, value := range stocks {
	//	klines := database.DB.GetKLineAllForStockCode(value.Code)
	//	up,_,_ := KLineIsUp(klines)
	//	if up && klines[len(klines)-1].Date.After(time.Now().Add(-time.Hour*24*3)) {
	//		utils.JJKPrintln("====")
	//		ss = append(ss, value)
	//	}
	//}
	//
	//for _, value := range ss {
	//	utils.JJKPrintln(value.Code)
	//}

	//klines := database.DB.GetKLineAllForStockCode("600019")
	//up,_,_ := KLineIsUp(klines)
	//utils.JJKPrintln(up)

	utils.JJKPrintln("end")
}
