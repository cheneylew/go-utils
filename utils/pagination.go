package utils

import (
	"math"
	"fmt"
	"strings"
)

type pagination struct {
	PageNum int64
	PageLink string
}

// curPageNum 从1开始，1表示第一页
func Pagination(url string, curPageNum, perPageCount, totalCount int64) string {
	var leftSideCount int64 = 3
	totalPages := int64(math.Ceil(float64(totalCount)/ float64(perPageCount)))
	JJKPrintln(totalPages)
	if totalPages <= 1 {
		return ""
	}

	params := TemplateParams()
	params["HasPrev"] = curPageNum > 1
	params["HasNext"] = curPageNum < totalPages
	params["CurPageNum"] = curPageNum
	params["PrevPage"] = pagination{
		PageNum:curPageNum - 1,
		PageLink:getLink(url,curPageNum - 1,perPageCount),
	}
	params["NextPage"] = pagination{
		PageNum:curPageNum + 1,
		PageLink:getLink(url,curPageNum + 1,perPageCount),
	}

	var lefts []pagination
	for i:= curPageNum-leftSideCount; i< curPageNum ; i++ {
		if i > 0 && i <= totalPages {
			lefts = append(lefts, pagination{PageNum:i, PageLink:getLink(url,i,perPageCount)})
		}
	}
	params["Lefts"] = lefts

	var rights []pagination
	for i:= curPageNum+1; i<= curPageNum+leftSideCount ; i++ {
		if i > 0 && i <= totalPages {
			rights = append(rights, pagination{PageNum:i, PageLink:getLink(url,i,perPageCount)})
		}
	}
	params["Rights"] = rights

	html := `
		<div class="pagination">
		{{if .HasPrev}}
			<span class="previous">
				<a href="{{.PrevPage.PageLink}}">&lt;</a>
			</span>
		{{end}}
			{{range .Lefts}}
			<span class="number">
				<a href="{{.PageLink}}">{{.PageNum}}</a>
			</span>
			{{end}}
			<b class="active">{{.CurPageNum}}</b>
			{{range .Rights}}
			<span class="number">
				<a href="{{.PageLink}}">{{.PageNum}}</a>
			</span>
			{{end}}
		{{if .HasNext}}
			<span class="next">
				<a href="{{.NextPage.PageLink}}">&gt;</a>
			</span>
		{{end}}
		</div>
	`
	return Template(html, params)
}

func getLink(url string, curPageNum int64, perPageCount int64) string {
	urls := strings.Split(url, "?")
	if len(urls) == 2 {
		url = urls[0]
		var items []string
		tmItems := strings.Split(urls[1], "&")
		for _, value := range tmItems {
			if !strings.Contains(value, "offset") && !strings.Contains(value, "limit") {
				items = append(items, value)
			}
		}
		params := strings.Join(items, "&")
		url = url+"?"+params
	}

	if strings.Contains(url, "?") {
		return fmt.Sprintf("%s&offset=%d&limit=%d", url, (curPageNum - 1)*perPageCount, perPageCount)
	} else {
		return fmt.Sprintf("%s?offset=%d&limit=%d", url, (curPageNum - 1)*perPageCount, perPageCount)
	}
}
