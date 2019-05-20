package main

import (
	"github.com/cheneylew/goutil/utils"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"fmt"
	"sort"
)

var myDB *gorm.DB

type TableInfo struct {
	DBName string
	TableName string
	TableSize int64
}

func mainSQLSlow() {
	//mysqlOnline()
	//slowSql()
	//readSlowLog()
	mainSQLExcel()
}

/*
mysqldumpslow -s c -t 2 /Users/dejunliu/Desktop/mysql-test-slow.log
这会输出记录次数最多的10条SQL语句，其中：
  ● -s, 是表示按照何种方式排序，c、t、l、r分别是按照记录次数、时间、查询时间、返回的记录数来排序，ac、at、al、ar，表示相应的倒叙；
  ● -t, 是top n的意思，即为返回前面多少条的数据；
  ● -g, 后边可以写一个正则匹配模式，大小写不敏感的；
 */
func slowSql()  {

	//result := utils.ExecShell("mysqldumpslow -s c -t 3 -g \"opc.\" /Users/dejunliu/Desktop/mysql-test-slow.log")
	result := utils.ExecShell("mysqldumpslow -s c -t 10 -g \"opc.\" /Users/dejunliu/Desktop/mysql-test-slow.log")
	utils.JJKPrintln(result)

	//得到访问次数最多的10个SQL
	//result := utils.ExecShell("mysqldumpslow -s c -t 10 /Users/dejunliu/Desktop/mysql-test-slow.log")
	//utils.JJKPrintln(result)

	//得到返回记录集最多的10个SQL。
	//result := utils.ExecShell("mysqldumpslow -s r -t 10 /Users/dejunliu/Desktop/mysql-test-slow.log")
	//utils.JJKPrintln(result)


	//得到按照时间排序的前10条里面含有左连接的查询语句。。
	//result := utils.ExecShell("mysqldumpslow -s t -t 10 -g \"left join\" /Users/dejunliu/Desktop/mysql-test-slow.log")
	//utils.JJKPrintln(result)
}

func initDB(dbName string)  {
	db, err := gorm.Open("mysql", fmt.Sprintf("root:ehsy2016@tcp(118.178.135.2:3306)/%s?charset=utf8&parseTime=True&loc=Local", dbName))
	if err != nil {
		utils.JJKPrintln(err)
		return
	} else {
		utils.JJKPrintln(fmt.Sprintf("%s database connected!", dbName))
	}
	myDB = db
}

func getOneColum(sql string) []string {
	var results []string
	rows, err := myDB.Raw(sql).Rows()

	if rows != nil {
		defer rows.Close()

		for rows.Next() {
			var a string
			rows.Scan(&a)
			results = append(results, a)
		}
	}

	if err != nil {
		utils.JJKPrintln(err)
		results = append(results, "0")
		return results
	}

	return results;
}

func removeString(slices []string, val string) []string {
	var results []string
	for _, value := range slices {
		if value != val {
			results = append(results, value)
		}
	}

	return results
}

func getDBs() []string {
	rs := getOneColum("show databases;")
	rs = removeString(rs, "information_schema")
	rs = removeString(rs, "front_collection_error")
	rs = removeString(rs, "performance_schema")
	rs = removeString(rs, "mysql")
	rs = removeString(rs, "test")
	return rs
}

func getTables() []string {
	return getOneColum("show tables;")
}

func getTableCount(tableName string)  {
}

func mysqlOnline()  {
	initDB("opc")
	dbNames := getDBs()
	var tableinfos []TableInfo
	for _, dbName := range dbNames {
		initDB(dbName)

		tables := getTables()

		for _, tablename := range tables {
			//results := getOneColum(fmt.Sprintf("select count(*) from %s", tablename))
			//count := results[0]
			count := ""
			tableinfos = append(tableinfos, TableInfo{dbName, tablename,utils.JKStrToInt64(count)})
		}

		myDB.Close()
	}

	sort.Slice(tableinfos, func(i, j int) bool {
		return tableinfos[i].TableSize < tableinfos[j].TableSize
	})

	for _, value := range tableinfos {
		utils.JJKPrintln(value.DBName,value.TableName, value.TableSize)
	}
}
