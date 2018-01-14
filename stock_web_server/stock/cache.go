package stock

import (
	"github.com/cheneylew/goutil/stock_web_server/models"
	"github.com/cheneylew/goutil/stock_web_server/database"
	"sort"
	"time"
)

var allKLines []*models.KLine
var allStocks []*models.Stock
var allStockInfos []*models.StockInfo

func InitCache()  {
	allKLines = database.DB.GetKLineAll()

	allStocks = database.DB.GetStockWithCodePrefix("60")
	allStocks = append(allStocks, database.DB.GetStockWithCodePrefix("00")...)

	allStockInfos = database.DB.GetStockInfoAll()
}

func CCGetStockAll() []*models.Stock {
	return allStocks
}

func CCGetStockInfoAll() []*models.StockInfo {
	return allStockInfos
}

func CCGetStockInfoWithStockId(stockId int64) []*models.StockInfo {
	var infos []*models.StockInfo
	for _, value := range allStockInfos {
		if value.StockId == stockId {
			infos = append(infos, value)
		}
	}

	sort.Slice(infos, func(i, j int) bool {
		return infos[i].Date.Before(infos[j].Date)
	})

	return infos
}

func CCGetStock(code string) *models.Stock {
	var stock *models.Stock
	for _, value := range allStocks {
		if value.Code == code {
			stock = value
		}
	}

	return stock
}

func CCGetKLinesWithCode(code string, recentCount int) []*models.KLine {
	var kl []*models.KLine
	stock := CCGetStock(code)
	for _, value := range allKLines {
		if value.StockId == stock.StockId {
			kl = append(kl, value)
		}
	}

	sort.Slice(kl, func(i, j int) bool {
		return kl[i].Date.Before(kl[j].Date)
	})

	if len(kl) < recentCount {
		return nil
	}
	return kl[len(kl)-recentCount:]
}

func CCGetKLineWithCodeAndDate(code string, date time.Time ) *models.KLine {
	var kl []*models.KLine
	stock := CCGetStock(code)
	for _, value := range allKLines {
		if value.StockId == stock.StockId {
			kl = append(kl, value)
		}
	}

	var res *models.KLine
	for _, value := range kl {
		if value.Date == date {
			res = value
			break
		}
	}

	return res
}

