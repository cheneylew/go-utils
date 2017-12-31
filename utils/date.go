package utils

import (
	"time"
	"fmt"
	"github.com/jinzhu/now"
)

/*
layout := "2006-01-02T15:04:05.000Z"
str := "2014-11-12T11:45:26.371Z"
t, err := time.Parse(layout, str)

if err != nil {
	fmt.Println(err)
}
fmt.Println(t)
 */
func StrToDateTime(str string) time.Time {
	layout := "20060102"
	date, _ := time.Parse(layout, str)
	return date
}

func DateEqual(t1, t2 time.Time) bool {
	y,m,d := t1.UTC().Date()
	y1,m1,d1 := t2.UTC().Date()
	return y == y1 && m == m1 && d == d1
}

func ValuesToDateTime(date, hour, minutes, ampm string) time.Time {
	finalHour := hour
	if ampm != "am" {
		finalHour = fmt.Sprintf("%d", JKStrToInt64(hour)+12)
	}

	datetimeStr := fmt.Sprintf("%s %s:%s:00", date, finalHour, minutes)
	datetime, _ := now.Parse(datetimeStr)
	return datetime
}