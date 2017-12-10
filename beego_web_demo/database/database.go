package database

import "github.com/cheneylew/goutil/utils/beego"
import _ "github.com/go-sql-driver/mysql"

var DB DataBase

func init() {
	db := beego.InitRegistDB("cheneylew","12344321","47.91.151.207","3308","shadowsocks-servers")
	DB = DataBase{
		BaseDataBase:*db,
	}
}

type DataBase struct {
	beego.BaseDataBase
}

func (db *DataBase)HelloWorld() {

}
