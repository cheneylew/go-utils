package main

import (
	"github.com/tealeg/xlsx"
	"github.com/cheneylew/goutil/utils"
	"strings"
	"sort"
	"fmt"
)

type SQLExcelRow struct {
	SQL string
	DBName string
	UseTimeRate float64
	UseTime float64
	BackRows int64
	ExcuteTimes int64
	LE3ms int64
	LE10ms int64
	LE1s int64
	GE1s int64
	RowID int
}

func mainSQLExcel()  {
	excelFileName := "/Users/dejunliu/Downloads/SQL统计-7.2(1).xlsx"
	xlFile, err := xlsx.OpenFile(excelFileName)
	if err != nil {
		utils.JJKPrintln(err)
		return
	}

	var sqls []SQLExcelRow
	for _, sheet := range xlFile.Sheets {
		for index, row := range sheet.Rows {
			if index != 0 {
				item := SQLExcelRow{
					row.Cells[0].String(),
					row.Cells[1].String(),
				utils.JKStrToFloat64(strings.Replace(row.Cells[2].String(),"%","",-1)),
					utils.JKStrToFloat64(row.Cells[3].String()),
					utils.JKStrToInt64(row.Cells[4].String()),
					utils.JKStrToInt64(row.Cells[5].String()),
					utils.JKStrToInt64(row.Cells[6].String()),
					utils.JKStrToInt64(row.Cells[7].String()),
					utils.JKStrToInt64(row.Cells[8].String()),
					utils.JKStrToInt64(row.Cells[9].String()),
					index+1,
				}

				if false {
					if strings.Contains(item.DBName,"opc") {
						sqls = append(sqls, item)
					}
				} else {
					sqls = append(sqls, item)
				}
			}
		}
	}

	if false {
		sort.Slice(sqls, func(i, j int) bool {
			return sqls[i].UseTime > sqls[j].UseTime
		})
	}

	for i:=0; i<26 ; i++ {
		sql := sqls[i]
		text := sql.SQL[0:20]
		text = strings.ToLower(text)
		utils.JJKPrintln(sql.RowID,text)
		if isInPHPProject(text) {
			utils.JJKPrintln(fmt.Sprintf("rowId:%d \tdbName:%s \tuseTime:%.2fms \t\tbackRows:%d words:%s", sql.RowID, sql.DBName, sql.UseTime, sql.BackRows, text))
		}
	}

	//utils.JJKPrintln("共计count:",len(sqls))
	//for _, sql := range sqls {
	//	text := sql.SQL[0:20]
	//	if isInPHPProject(text) {
	//		utils.JJKPrintln(fmt.Sprintf("rowId:%d \tdbName:%s \tuseTime:%.2fms \t\tbackRows:%d words:%s", sql.RowID, sql.DBName, sql.UseTime, sql.BackRows, text))
	//	}
	//}
}

func isInPHPProject(words string) bool {
	text := utils.ExecShell(fmt.Sprintf("cd /Users/dejunliu/Desktop/ehys/php/opc; pt -i \"%s\"", words))
	return len(text) > 0
}
