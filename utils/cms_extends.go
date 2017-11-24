package utils

import (
	"math"
	"fmt"
	"time"
)

func DateTimeStringToPickerHtml(fieldName string, timeString string) string  {
	return DateTimeToPickerHtml(fieldName, JKStringToTime(timeString))
}

func DateTimeToPickerHtml(fieldName string, time time.Time) string {
	publishDateHtml := ""
	d, h, m, ampm := JKDateTimeSplit(JKTimeToString(time))
	publishDateHtml += fmt.Sprintf(`<input name="%s" value="%s" style="width: 80px" class="text datetime datepick dp-applied" type="text">`, fieldName, d)
	publishDateHtml += fmt.Sprintf(`<select name="%s_hour">`, fieldName)
	for i := 1; i <= 12; i++ {
		if h == fmt.Sprintf("%02d", i) {
			publishDateHtml += fmt.Sprintf(`<option value="%02d" selected="selected">%02d</option>`, i, i)
		} else {
			publishDateHtml += fmt.Sprintf(`<option value="%02d">%02d</option>`, i, i)
		}
	}
	publishDateHtml += `</select>`
	publishDateHtml += fmt.Sprintf(`<select name="%s_minute">`, fieldName)
	for i := 0; i <= 59; i++ {
		if m == fmt.Sprintf("%02d", i) {
			publishDateHtml += fmt.Sprintf(`<option value="%02d" selected="selected">%02d</option>`, i, i)
		} else {
			publishDateHtml += fmt.Sprintf(`<option value="%02d">%02d</option>`, i, i)
		}
	}
	publishDateHtml += `</select>`
	publishDateHtml += fmt.Sprintf(`<select name="%s_ampm">`, fieldName)
	if ampm == "am" {
		publishDateHtml += `	<option value="am" selected="selected">am</option>
								<option value="pm" >pm</option>`
	} else {
		publishDateHtml += `	<option value="am">am</option>
								<option value="pm" selected="selected">pm</option>`
	}
	publishDateHtml += `</select>`

	return publishDateHtml
}

var (
	LeftRightPageCount int64 = 3
)

type Pagenation struct {
	TotalRows 	int64
	PerPageCount 	int64
	Offset 		int64
	CommonParams 	string
}

func (pagenation *Pagenation) CurPageNum() int64 {
	return int64(math.Ceil(float64(pagenation.Offset)/float64(pagenation.PerPageCount))+1)
}

func (pagenation *Pagenation) TotalPages() int64 {
	return int64(math.Ceil(float64(pagenation.TotalRows) / float64(pagenation.PerPageCount)))
}

func (pagenation *Pagenation) havePrevious() bool {
	if pagenation.CurPageNum()-1 > 0 {
		return true
	} else {
		return false
	}
}

func (pagenation *Pagenation) haveNext() bool {
	if pagenation.CurPageNum()+1 <= pagenation.TotalPages() {
		return true
	} else {
		return false
	}
}

func (pagenation *Pagenation) shouldShowFirstPage() bool {
	if pagenation.CurPageNum() - LeftRightPageCount > 1 {
		return true
	}

	return false
}

func (pagenation *Pagenation) shouldShowLastPage() bool {
	if pagenation.CurPageNum() + LeftRightPageCount < pagenation.TotalPages() {
		return true
	}

	return false
}

func PageDivUtil(baseURL string, total int64, offset int64, limit int64, showtotal bool, commonParams string) string {
	page := &Pagenation{
		TotalRows:total,		//total
		PerPageCount:limit,		//limit
		Offset:offset,			//start
		CommonParams:commonParams,	//公共参数传递
	}
	return PageDiv(baseURL, page, showtotal)
}

func PageDiv(curURL string, pagenation *Pagenation, showCount bool) string {
	if pagenation.TotalPages() <= 1 {
		return ""
	}

	str := ""
	str += `<div class="pagination">`
	if showCount {
		str += fmt.Sprintf(`<div class="dataset_stats">共 <b>%d</b> 条记录</div>`, pagenation.TotalRows)
	}
	str += `<div class="pagination">`

	//第一页
	if pagenation.shouldShowFirstPage() {
		str += fmt.Sprintf(
			`<span class="previous"><a href="%s/?limit=%d&offset=%d&%s">第1页</a></span>`,
			curURL,
			pagenation.PerPageCount,
			0,
			pagenation.CommonParams,
		)
	}

	//上一页
	if pagenation.havePrevious() {
		str += fmt.Sprintf(
			`<span class="previous"><a href="%s/?limit=%d&offset=%d&%s">&lt;</a></span>`,
			curURL,
			pagenation.PerPageCount,
			pagenation.Offset - pagenation.PerPageCount,
			pagenation.CommonParams,
		)
	}
	//当前页的前面几页
	prevFirstPageNum := pagenation.CurPageNum() - LeftRightPageCount
	if prevFirstPageNum <= 0 {
		prevFirstPageNum = 1
	}
	for i := prevFirstPageNum; i < pagenation.CurPageNum() ; i++ {
		str += fmt.Sprintf(
			`<span class="number"><a href="%s/?limit=%d&offset=%d&%s">%d</a></span>`,
			curURL,
			pagenation.PerPageCount,
			pagenation.Offset - pagenation.PerPageCount*(pagenation.CurPageNum() - i),
			pagenation.CommonParams,
			i,
		)
	}

	//当前页
	str += fmt.Sprintf(
		`<b class="active">%d</b>`,
		pagenation.CurPageNum(),
	)

	//当前页的后面几页
	nextLastPageNum := pagenation.CurPageNum() + LeftRightPageCount
	if nextLastPageNum > pagenation.TotalPages() {
		nextLastPageNum = pagenation.TotalPages()
	}
	for i := pagenation.CurPageNum()+1; i <= nextLastPageNum ; i++ {
		str += fmt.Sprintf(
			`<span class="number"><a href="%s/?limit=%d&offset=%d&%s">%d</a></span>`,
			curURL,
			pagenation.PerPageCount,
			pagenation.Offset - pagenation.PerPageCount*(pagenation.CurPageNum() - i),
			pagenation.CommonParams,
			i,
		)
	}

	//下一页
	if pagenation.haveNext() {
		str += fmt.Sprintf(
			`<span class="next"><a href="%s/?limit=%d&offset=%d&%s">&gt;</a></span>`,
			curURL,
			pagenation.PerPageCount,
			pagenation.Offset + pagenation.PerPageCount,
			pagenation.CommonParams,
		)
	}

	//最后一页
	if pagenation.shouldShowLastPage() {
		str += fmt.Sprintf(
			`<span class="previous"><a href="%s/?limit=%d&offset=%d&%s">第%d页</a></span>`,
			curURL,
			pagenation.PerPageCount,
			pagenation.TotalPages()*pagenation.PerPageCount - pagenation.PerPageCount,
			pagenation.CommonParams,
			pagenation.TotalPages(),
		)
	}


	str += `</div>`
	str += `</div>`

	return str
}

