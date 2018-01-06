package stock

import (
	"github.com/cheneylew/goutil/utils"
)

func CronMain()  {
	//周一到周五，23:00:00执行
	utils.CronJob("00 00 23 * * 1-5", func() {
		downloadSHStockKLines()
		downloadSZStockKLines()
		downloadStockRealTimeInfo()
	})
	//每天23:00:00执行
	utils.CronJob("00 00 23 * * ?", func() {
		uploadStocksCodeToDB()
		downloadStockRealTimeInfo()
	})
}
