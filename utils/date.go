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

func DateEqual(t1, t2 time.Time) bool {
	y,m,d := t1.UTC().Date()
	y1,m1,d1 := t2.UTC().Date()
	return y == y1 && m == m1 && d == d1
}