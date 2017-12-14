package models

import (
	"time"
	"strings"
	"fmt"
)

type Response struct {
	Code int64
	Msg string
	Data map[string]map[string]interface{}
}

type SortAnalysDayKLins []*AnalysDayKLine

func (a SortAnalysDayKLins) Len() int           { return len(a) }
func (a SortAnalysDayKLins) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortAnalysDayKLins) Less(i, j int) bool {
	if a[i].RedCount < a[j].RedCount {
		return true
	}

	return false
}

type AnalysDayKLine struct {
	Stock *Stock
	Days int64
	RedCount int64
	GreenCount int64
	UpCount int64
	DownCount int64
}



type Stock struct {
	StockId int64       `orm:"pk"`
	Code string
	SyncTime time.Time
	SyncOk bool
}

func (s *Stock)CodeStr() string {
	if strings.HasPrefix(s.Code, "60") {
		return fmt.Sprintf("sh%s",s.Code)
	} else if strings.HasPrefix(s.Code, "00") || strings.HasPrefix(s.Code, "30") {
		return fmt.Sprintf("sz%s",s.Code)
	}

	return s.Code
}

type KLine struct {
	KLineId int64       `orm:"pk"`
	StockId int64
	OpeningPrice float64
	ClosingPrice float64
	MaxPrice float64
	MinPrice float64
	Date time.Time
	Vol float64 		//万手
	Type int		//1 日K 2 周K 3月K 4年K
}

func (k *KLine)IsRed() bool {
	if k.ClosingPrice >= k.OpeningPrice {
		return true
	}
	return false
}

func (k *KLine)GetAddRate(last *KLine) float64 {
	if last == nil {
		return 0
	}
	return (k.ClosingPrice - last.ClosingPrice)/last.ClosingPrice
}

type StockInfo struct {
	Qfqday [][]string
	Qt map[string][]string
	Prec string
	Version string
}