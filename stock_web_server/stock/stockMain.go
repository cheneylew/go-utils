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
	shStocks := database.DB.GetStockWithCodePrefix("60")
	for _, value := range shStocks {
		downloadDayKLine(value.Code)
	}
	utils.JJKPrintln(len(shStocks))
}

func downloadSZStockKLines()  {
	shStocks := database.DB.GetStockWithCodePrefix("00")
	for _, value := range shStocks {
		downloadDayKLine(value.Code)
	}
	utils.JJKPrintln(len(shStocks))
}


func downloadFaildStocks()  {
	for _, value := range database.DB.GetSyncFailedStocks() {
		downloadDayKLine(value.Code)
	}

	fstocks := database.DB.GetSyncFailedStocks()
	utils.JJKPrintln(len(fstocks))
}

func downloadStockRealTimeInfo()  {
	stks := database.DB.GetStockWithCodePrefix("00")
	for _, tmpStk := range stks {
		code := tmpStk.Code
		stk := database.DB.GetStockWithCode(code)
		info := GetRealTimeStockInfo(stk.CodeStr())
		if len(info) == 0 {
			utils.JJKPrintln(fmt.Sprintf("%s failed", code))
		} else {
			utils.JJKPrintln(info)
			for _, value := range info {
				value.StockId = stk.StockId
				_,e := database.DB.Orm.Insert(value)
				if e != nil {
					utils.JJKPrintln(e)
				}
			}
		}
	}
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
		infos := database.DB.GetStockInfoAllForStock(stocks[i])
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

func StockTestMain()  {
	CronMain()

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
