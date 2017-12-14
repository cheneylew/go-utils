package utils

import "time"

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
