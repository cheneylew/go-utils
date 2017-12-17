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
	"sort"
	"math"
)

const KLineDayUrl = "http://web.ifzq.gtimg.cn/appstock/app/fqkline/get?_var=kline_dayqfq&param=%s,day,,,%d,qfq&r=0.7273883641103628"
const KLineWeekUrl = "http://web.ifzq.gtimg.cn/appstock/app/fqkline/get?_var=kline_dayqfq&param=%s,week,,,%d,qfq&r=0.7273883641103628"
const KRealTimeUrl = "http://qt.gtimg.cn/q=ff_%s"

func GetKLineDayUrl(code string, count int64) string {
	return fmt.Sprintf(KLineDayUrl,code, count)
}

func GetKLineWeekUrl(code string, count int64) string {
	return fmt.Sprintf(KLineWeekUrl,code, count)
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

func GetStockWeekKLine(code string, count int64) []*models.KLine {
	url := GetKLineWeekUrl(code, count)
	result := utils.HTTPGet(url)
	jsonStr := strings.Replace(result,"kline_dayqfq=","",-1)

	m := new(models.Response)
	err := json.Unmarshal([]byte(jsonStr), &m)
	if err != nil {
		utils.JJKPrintln(err)
		return nil
	}
	infos := m.Data[code]

	dayKLineArray, ok := infos["qfqweek"].([]interface{})
	if !ok {
		dayKLineArray, ok = infos["week"].([]interface{})
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
				Type:2,
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
		_, err := database.DB.Orm.Insert(value)
		if err != nil {
			utils.JJKPrintln(err)
		}
	}
	utils.JJKPrintln("upload finished")
}

func GetRealTimeStockInfo(code string) []*models.StockInfo {
	url := fmt.Sprintf(KRealTimeUrl,code)
	s := utils.HTTPGet(url)
	s1 := utils.RegexpFindAll(s,`".+"`)
	var result []*models.StockInfo
	if len(s1) > 0 {
		infos := strings.Split(strings.Replace(s1[0],`"`,"",-1),"~")

		info := &models.StockInfo{
			MainIn:utils.ToFloat64(infos[1]),
			MainOut:utils.ToFloat64(infos[2]),
			MainTotal:utils.ToFloat64(infos[3]),
			RetailIn:utils.ToFloat64(infos[5]),
			RetailOut:utils.ToFloat64(infos[6]),
			RetailTotal:utils.ToFloat64(infos[7]),
			Date:utils.StrToDateTime(infos[13]),
		}

		info1 := strToStockInfo(infos[14])
		info2 := strToStockInfo(infos[15])
		info3 := strToStockInfo(infos[16])
		info4 := strToStockInfo(infos[17])

		result = append(result, info, info1, info2, info3, info4)
	}
	return result
}

func strToStockInfo(str string) *models.StockInfo {
	arr := strings.Split(str,"^")
	return &models.StockInfo{
		Date:utils.StrToDateTime(arr[0]),
		MainIn:utils.ToFloat64(arr[1]),
		MainOut:utils.ToFloat64(arr[2]),
		MainTotal:utils.ToFloat64(arr[1])-utils.ToFloat64(arr[2]),
	}
}

func KLineIsUp(klines []*models.KLine) (up bool,upCnt int, downCnt int) {
	durationDays := 30
	if len(klines) > durationDays {
		klines = klines[len(klines)-durationDays:]
	} else if len(klines) == 0 || len(klines) < 8 {
		return false, 0, 0
	}

	//升序排列
	sortKLines := models.SortKLine{}
	for _, value := range klines {
		sortKLines = append(sortKLines, value)
	}
	sort.Sort(sortKLines)

	isUp := false
	var last *models.KLine
	upcount := 0
	downCount := 0
	for _, kline := range sortKLines {
		if last != nil {
			if kline.ClosingPrice > last.ClosingPrice {
				upcount += 1
			} else {
				downCount += 1
			}
		}

		last = kline
	}

	parts := 2
	step := math.Floor(float64(len(sortKLines))/float64(parts))
	lastIndx := 0
	lineOk := true
	for i := 1; i< parts ; i++ {
		idx := i * int(step)
		if sortKLines[lastIndx].ClosingPrice > sortKLines[idx].ClosingPrice {
			lineOk = false
			break
		}

		lastIndx = idx
	}

	firstOne := sortKLines[0]
	lastOne := sortKLines[len(sortKLines)-1]
	if sortKLines[lastIndx].ClosingPrice > lastOne.ClosingPrice {
		lineOk = false
	}
	deltaRate := (lastOne.ClosingPrice-firstOne.ClosingPrice)/firstOne.ClosingPrice
	deltaRate = deltaRate*100;
	if deltaRate > 6.0 && lineOk {
		isUp = true
	}

	return isUp,upcount, downCount
}


