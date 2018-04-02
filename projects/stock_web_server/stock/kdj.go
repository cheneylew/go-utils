package stock

import (
	"github.com/cheneylew/goutil/projects/stock_web_server/models"
)

func calculateKDJ(code string, count int64) (todayKDJ float64, lines []*models.KLine)  {
	klines := GetStockDayKLine(code,count)
	dayNum := 9
	lastVK := 0.0
	lastVD := 0.0
	for i:=0; i<len(klines) ; i++  {
		if i>=dayNum {
			//RSVt＝(Ct－L9)／(H9－L9)＊100
			HP := 0.0
			LP := 0.0
			startIndex := i-dayNum+1
			for j:= startIndex; j<=i ; j++ {
				if j == startIndex {
					HP = klines[j].MaxPrice
					LP = klines[j].MinPrice
				} else {
					if klines[j].MaxPrice > HP {
						HP = klines[j].MaxPrice
					}
					if klines[j].MinPrice < LP {
						LP = klines[j].MinPrice
					}
				}
			}
			rsvt := (klines[i].ClosingPrice-LP)/(HP-LP)*100
			curVK := (2.0/3.0)*lastVK+(1.0/3.0)*rsvt
			curVD := (2.0/3.0)*lastVD+(1.0/3.0)*curVK
			curVJ := 3*curVK-2*curVD

			klines[i].Kdj_k = curVK
			klines[i].Kdj_d = curVD
			klines[i].Kdj_j = curVJ

			lastVK = curVK
			lastVD = curVD
		} else {
			lastVK = 50.0
			lastVD = 50.0
		}
	}

	return 0, klines
}

