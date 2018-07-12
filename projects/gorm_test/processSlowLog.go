package main

import (
	"github.com/cheneylew/goutil/utils"
	"strings"
	"sort"
	"fmt"
)

type SQLLogParser struct {
	TimeStr string
	QueryTimeStr string
	SQLStr string

	QueryTime float64
}

func readSlowLog()  {
	lines := utils.FileReadLineString("/Users/dejunliu/Desktop/mysql-test-slow.log")
	var logs []SQLLogParser
	index := 0
	for index < len(lines) {
		line := lines[index]
		if strings.Contains(line, "# Time") {
			detailTimeStr := line
			queryTimeStr := lines[index+2]

			//拼接SQL
			sqlStr := ""
			i := index+3
			for i<len(lines) && !strings.Contains(lines[i],"# Time") {
				sqlStr += "\n"
				sqlStr += lines[i]
				i ++
			}

			queryTimeSlices := strings.Split(queryTimeStr, " ")
			parser := SQLLogParser{
				detailTimeStr,
				queryTimeStr,
				sqlStr,
				utils.JKStrToFloat64(queryTimeSlices[2])}
			logs = append(logs, parser)
		}
		index ++
	}

	sort.Slice(logs, func(i, j int) bool {
		return logs[i].QueryTime > logs[j].QueryTime
	})

	count := 0
	for _, value := range logs {
		if value.QueryTime > 4 {
			sql := strings.Replace(value.SQLStr,"","",-1)
			if strings.Contains(sql,"opc.") {
				utils.JJKPrintln(value.QueryTime)
				writeLog(fmt.Sprintf("====%f %s", value.QueryTime, sql))
				count++
			}
		}
	}

	utils.JJKPrintln("count:", count)
}
