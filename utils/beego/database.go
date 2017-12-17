package beego

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"

	"github.com/astaxie/beego/orm"
	"github.com/cheneylew/goutil/utils"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func DBUrl(user, password, host, port, dbName string) string {
	return fmt.Sprintf(`%s:%s@tcp(%s:%s)/%s?charset=utf8`, user, password, host, port, dbName)
}

func InitRegistDB(user,pwd,host,port,dbname string) *BaseDataBase {
	url := DBUrl(user,pwd,host,port,dbname)
	orm.RegisterDriver("mysql", orm.DRMySQL)
	err := orm.RegisterDataBase("default", "mysql", url)

	if err != nil {
		utils.JJKPrintln("========database can't connect! error:" + err.Error()+"========")
	} else {
		utils.JJKPrintln("========database connected success！========")
	}

	GDB, err := gorm.Open("mysql",url)
	if err != nil {
		utils.JJKPrintln(err)
	}
	//表名不加s，gorm默认加s
	GDB.SingularTable(true)

	return &BaseDataBase {
		Orm:orm.NewOrm(),
		GDB:GDB,
	}
}

type BaseDataBase struct {
	Orm orm.Ormer
	GDB *gorm.DB
}


