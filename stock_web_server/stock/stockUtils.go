package stock

import (
	"strings"
	"fmt"
	"github.com/cheneylew/goutil/stock_web_server/models"
	"encoding/json"
	"github.com/jinzhu/now"
	"github.com/cheneylew/goutil/utils"
)

const KLineDayUrl = "http://web.ifzq.gtimg.cn/appstock/app/fqkline/get?_var=kline_dayqfq&param=%s,day,,,%d,qfq&r=0.7273883641103628"

func GetKLineDayUrl(code string, count int64) string {
	return fmt.Sprintf(KLineDayUrl,code, count)
}

func GetStockDayKLine(code string, count int64) []*models.KLineDay {

	result := utils.HTTPGet(GetKLineDayUrl(code, count))
	jsonStr := strings.Replace(result,"kline_dayqfq=","",-1)


	m := new(models.Response)
	err := json.Unmarshal([]byte(jsonStr), &m)
	if err != nil {
		utils.JJKPrintln(err)
		return nil
	}
	infos := m.Data[code]

	dayKLineArray, ok := infos["qfqday"].([]interface{})
	dayKLineResult := make([]*models.KLineDay, 0)
	if ok {
		for _, value := range dayKLineArray {
			oneday, onedayOk := value.([]interface{})
			dayValues := []string{}
			if onedayOk {
				for _, tv := range oneday {
					dayValues = append(dayValues,utils.ToString(tv))
				}
			}
			date, _ := now.Parse(dayValues[0])
			s := &models.KLineDay{
				Date:date,
				OpeningPrice:utils.ToFloat64(dayValues[1]),
				ClosingPrice:utils.ToFloat64(dayValues[2]),
				MaxPrice:utils.ToFloat64(dayValues[3]),
				MinPrice:utils.ToFloat64(dayValues[4]),
			}
			dayKLineResult = append(dayKLineResult,s)
		}
	}

	return dayKLineResult
}



