package database

import (
	"github.com/cheneylew/goutil/utils/beego"
	_ "github.com/go-sql-driver/mysql"
	"github.com/cheneylew/goutil/stock_web_server/models"

	"github.com/cheneylew/goutil/utils/beego"
)

var DB DataBase

func init() {
	db := beego.InitRegistDB("cheneylew","12344321","47.91.151.207","3308","shadowsocks-servers")
	db.Orm.Using("default")
	DB = DataBase{
		BaseDataBase:*db,
	}
}

type DataBase struct {
	beego.BaseDataBase
}

func (db *DataBase)HelloWorld() {

}

func (db *DataBase)GetUser() *models.User {
	return nil
}
