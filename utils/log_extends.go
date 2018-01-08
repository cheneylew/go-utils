package utils

import (
	"fmt"
	"strings"
	"encoding/json"
	"github.com/astaxie/beego/logs"
)

func Log(a ...interface{})  {
	fmt.Println(a...)
}

func JJKPrintln(a ...interface{})  {
	loger := logs.GetLogger("")
	loger.Output(3, fmt.Sprintln(a...))
}

func Println(a ...interface{})  {
	Log(a...)
}

func JKCheckError(err error)  {
	if err != nil {
		Println(err)
	}
}

func LogA(a ...interface{})  {
	fmt.Println(a...)
}

func Pt(s ...interface{})  {
	if len(s) == 1 {
		fmt.Println(s[0])
	}else {
		for i:=0;i<len(s);i++ {
			if i == len(s)-1 {
				fmt.Print(s[i])
			}else {
				fmt.Print(s[i],",")
			}
		}
		fmt.Print("\n")
	}
}

func JKFmt(str string,s ...interface{})  {
	fmt.Printf(str,s[0])
}

func JKFormat(format string,a... interface{}) string {
	return fmt.Sprintf(format,a...)
}

func StringArrToInterfaceArr(strings []string) []interface{} {
	s := make([]interface{}, len(strings))
	for i, v := range strings {
		s[i] = v
	}
	return s;
}

func Trim(s string) string {
	return strings.Trim(strings.Trim(strings.Trim(s," "), string('\n')), string('\r'))
}

func JKJSON(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		JJKPrintln(err)
		return ""
	}

	return string(b)
}

/*
str: 	2017-04-07 18:38:00

return
date:	2017-04-07
hour: 	12
minute: 59
ampm: 	am/pm
 */
func JKDateTimeSplit(str string) (date string, hour string, minute string,ampm string) {
	datetime := Trim(str)
	strs := strings.Split(datetime, " ")
	if len(strs) == 2 {
		date := strs[0]
		time := strs[1]
		pices := strings.Split(time, ":")
		if len(pices) == 3 {
			h := pices[0]
			m := pices[1]
			ap := "am"
			if JKStrToInt(h) > 12 {
				ap = "pm"
				h = fmt.Sprintf("%02d", JKStrToInt(h) - 12)
			}

			return date, h, m, ap
		}

		return "", "", "", ""
	}
	return "", "", "", ""
}

