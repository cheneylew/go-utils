package main

import (
	"github.com/astaxie/beego/config"
	"log"
	"fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/cheneylew/goutil/utils"
	"github.com/wulijun/go-php-serialize/phpserialize"
	"path"
	"time"
	"strings"
	"os"
)

var ini_conf config.Configer

func init() {
	var err error
	iniconf, err := config.NewConfig("ini", "db.ini")
	if err != nil {
		log.Fatal(err)
	}
	ini_conf = iniconf
	db_host := iniconf.String("mysql::db_host")
	db_user := iniconf.String("mysql::db_user")
	db_password := iniconf.String("mysql::db_password")
	db_port := iniconf.String("mysql::db_port")
	db_name := iniconf.String("mysql::db_name")

	dialect := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", db_user, db_password, db_host, db_port, db_name)
	orm.RegisterDataBase("default", "mysql", dialect, 30)
}

func write_log(log string)  {
	logPath := path.Join(utils.ExeDir(), "data.log")
	fd,_:=os.OpenFile(logPath,os.O_RDWR|os.O_CREATE|os.O_APPEND,0644)
	fd_time:=time.Now().Format("2006-01-02 15:04:05");
	fd_content:=strings.Join([]string{fd_time," ",log,"\n"},"")
	buf:=[]byte(fd_content)
	fd.Write(buf)
	fd.Close()
}

// addslashes() 函数返回在预定义字符之前添加反斜杠的字符串。
// 预定义字符是：
// 单引号（'）
// 双引号（"）
// 反斜杠（\）
func Addslashes(str string) string {
	tmpRune := []rune{}
	strRune := []rune(str)
	for _, ch := range strRune {
		switch ch {
		case []rune{'\\'}[0], []rune{'"'}[0], []rune{'\''}[0]:
			tmpRune = append(tmpRune, []rune{'\\'}[0])
			tmpRune = append(tmpRune, ch)
		default:
			tmpRune = append(tmpRune, ch)
		}
	}
	return string(tmpRune)
}

// stripslashes() 函数删除由 addslashes() 函数添加的反斜杠。
func Stripslashes(str string) string {
	dstRune := []rune{}
	strRune := []rune(str)
	strLenth := len(strRune)
	for i := 0; i < strLenth; i++ {
		if strRune[i] == []rune{'\\'}[0] {
			i++
		}
		dstRune = append(dstRune, strRune[i])
	}
	return string(dstRune)
}

func start_task()  {
	max_count,_ := ini_conf.Int64("conf::max_count")
	per_num,_ := ini_conf.Int64("conf::per_num")
	start_id,_ := ini_conf.Int64("conf::start_id")
	concur_num,_ := ini_conf.Int64("conf::concur_num")

	var finished_count int64 = 0
	for {
		per_count := multi_task(start_id, per_num, concur_num)
		finished_count += per_count
		start_id += per_num

		percentage := float64(finished_count)/float64(max_count)
		if percentage > 1 {
			percentage = 1
		}
		write_log(fmt.Sprintf("已完成%.2f%%", percentage*100))
		if finished_count > max_count {
			break
		}
	}
}

func multi_task(start_id,step,concur_num int64) int64 {
	sql_select := fmt.Sprintf("select id,data from oc_data_edit_log where id between %d and %d+%d;", start_id,start_id,step)
	o := orm.NewOrm()
	var maps []orm.Params
	num, err := o.Raw(sql_select).Values(&maps)
	var results[]map[string]string
	if err == nil && num > 0 {
		for _, item := range maps {
			id := item["id"].(string)
			data := item["data"]
			mp,e := phpserialize.Decode(data.(string))
			if e == nil {
				mp2 := make(map[string]interface{})
				for key, value := range mp.(map[interface{}]interface{}) {
					switch key := key.(type) {
					case string:
						mp2[key] = value
					}
				}
				if mp2["SERVER"] != nil || mp2["POST"] != nil  {
					delete(mp2,"SERVER")
					delete(mp2,"POST")
					mp3 := make(map[interface{}]interface{})
					for key, value := range mp2 {
						mp3[key] = value
					}
					newData, e := phpserialize.Encode(mp3)
					if e == nil {
						m := make(map[string]string)
						m[id] = newData
						results = append(results, m)
					}
				}
			}
		}
	}

	if len(results) > 0 {
		groups := make(map[string][]map[string]string)
		for key, value := range results {
			gid := int64(key)%concur_num
			gkey := fmt.Sprintf("%d", gid)
			groups[gkey] = append(groups[gkey], value)
		}

		ch := make(chan int)
		for _, value := range groups {
			go one_task(value, ch)
		}

		task_count := len(groups)
		count := 0
		for {
			select {
			case c := <-ch:
				count += c
			}
			if count == task_count {
				close(ch)
				break
			}
		}
	} else {
		//utils.JJKPrintln("No items need process!")
	}

	return num
}

func one_task(rows []map[string]string, ch chan int)  {
	//批量update语句
	conditions := ""
	ids := ""
	for index, row := range rows {
		for id, data := range row {
			conditions += fmt.Sprintf(" WHEN %s THEN '%s'", id, Addslashes(data))
			if index == len(rows)-1 {
				ids += fmt.Sprintf("%s", id)
			} else {
				ids += fmt.Sprintf("%s,", id)
			}
		}
	}
	sql := fmt.Sprintf("update oc_data_edit_log set data = CASE id %s END where id in (%s)", conditions, ids)

	o := orm.NewOrm()
	res, err := o.Raw(sql).Exec()
	if err == nil {
		num, _ := res.RowsAffected()
		fmt.Println("mysql row affected nums: ", num)
	} else {
		write_log(fmt.Sprintf("Error:%s SQL:%s", err.Error(), sql))
	}
	ch<-1
}

func main() {
	start_task()
	utils.JJKPrintln("All task finished!")
	write_log("All task finished!")
}

