package models

import "github.com/astaxie/beego/orm"

func init() {
	orm.RegisterModel(new(KLine), new(Stock), new(StockInfo), new(Tester), new(MyStock))
}


