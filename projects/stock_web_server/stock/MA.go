package stock

import "github.com/cheneylew/goutil/projects/stock_web_server/models"

func caculateMA(lines []*models.KLine, days int) []*models.KLine {
	for i:=days; i<len(lines) ; i++ {
		total := 0.0
		for j:=i; j>i-days ; j-- {
			total += lines[j].ClosingPrice
		}
		lines[i].MaVal = total/float64(days)
	}
	return lines
}
