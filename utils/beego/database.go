package beego

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"

	"github.com/astaxie/beego/orm"
	"github.com/cheneylew/goutil/utils"

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

	//创建模型表结构
	orm.RunSyncdb("default",false,true)

	return &BaseDataBase {
		Orm:orm.NewOrm(),
	}
}

type BaseDataBase struct {
	Orm orm.Ormer
}

func (db *BaseDataBase)DBBaseTableCount(tablename string) int64 {
	a, err := db.Orm.QueryTable(tablename).Count()
	if err != nil {
		return 0
	}

	return a
}


