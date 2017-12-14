package stock

import (
	"strings"
	"fmt"
	"github.com/cheneylew/goutil/stock_web_server/models"
	"encoding/json"
	"github.com/jinzhu/now"
	"github.com/cheneylew/goutil/utils"
	"time"
	"github.com/cheneylew/goutil/stock_web_server/database"
)

const KLineDayUrl = "http://web.ifzq.gtimg.cn/appstock/app/fqkline/get?_var=kline_dayqfq&param=%s,day,,,%d,qfq&r=0.7273883641103628"

func GetKLineDayUrl(code string, count int64) string {
	return fmt.Sprintf(KLineDayUrl,code, count)
}

func GetStockDayKLine(code string, count int64) []*models.KLine {
	url := GetKLineDayUrl(code, count)
	utils.JJKPrintln(url)
	result := utils.HTTPGet(url)
	jsonStr := strings.Replace(result,"kline_dayqfq=","",-1)

	m := new(models.Response)
	err := json.Unmarshal([]byte(jsonStr), &m)
	if err != nil {
		utils.JJKPrintln(err)
		return nil
	}
	infos := m.Data[code]

	dayKLineArray, ok := infos["qfqday"].([]interface{})
	if !ok {
		dayKLineArray, ok = infos["day"].([]interface{})
	}
	dayKLineResult := make([]*models.KLine, 0)
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
			date = date.Add(time.Hour *8)
			s := &models.KLine{
				Date:date,
				OpeningPrice:utils.ToFloat64(dayValues[1]),
				ClosingPrice:utils.ToFloat64(dayValues[2]),
				MaxPrice:utils.ToFloat64(dayValues[3]),
				MinPrice:utils.ToFloat64(dayValues[4]),
				Vol:utils.ToFloat64(dayValues[5]),
				Type:1,
			}
			dayKLineResult = append(dayKLineResult,s)
		}
	}
	return dayKLineResult
}

func GetCodesOnline() []string {
	html := utils.HTTPGet("http://quote.eastmoney.com/stocklist.html")
	codes := utils.RegexpFindAll(html,`\d{6}`)
	return codes
}

func UpdateOnlineCodesToDatabase()  {
	codes := GetCodesOnline()
	codeModels := make([]*models.Stock,0)
	for _, code := range codes {
		aStock := &models.Stock{
			Code:code,
		}
		codeModels = append(codeModels, aStock)
	}
	utils.JJKPrintln("start upload")
	for _, value := range codeModels {
		database.DB.Orm.Insert(value)
	}
	utils.JJKPrintln("upload finished")
}




