package main

import (
	"github.com/cheneylew/goutil/stock_web_server/stock"
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
	stock.UpdateOnlineCodesToDatabase()
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

	klines := stock.GetStockDayKLine(tmpStock.CodeStr(),int64(days))
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

func analysResultWithCodeAndDays(code string, days int) *models.AnalysDayKLine {
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
		} else {
			result.DownCount += 1
		}
		lastKLine = value
	}

	return result
}

func StockTestMain()  {
	//uploadStocksCodeToDB()
	//utils.JJKPrintln(len(database.DB.GetKLineAll()))
	//downloadSHStockKLines()
	//downloadSZStockKLines()
	//downloadFaildStocks()
	//downloadDayKLine("600196")
	//stock.GetRealTimeStockInfo("sh600703")


	if true {
		shStocks := database.DB.GetStockWithCodePrefix("60")
		var all models.SortAnalysDayKLins
		for _, value := range shStocks {
			result := analysResultWithCodeAndDays(value.Code,18)
			all = append(all, result)
		}
		sort.Sort(all)

		result := ""
		for _, value := range all {
			utils.JJKPrintln(value.Stock.Code, value.RedCount, value.GreenCount)
			if value.RedCount + value.GreenCount > 10 {
				result += fmt.Sprintf("%s %d %d\n", value.Stock.Code, value.RedCount, value.GreenCount)
			}
		}

		utils.FileWriteString("/Users/apple/Desktop/a.txt", result)
	}

	utils.JJKPrintln("end")
}
