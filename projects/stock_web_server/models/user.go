package models

import (
	"github.com/cheneylew/goutil/utils/beego"
	"github.com/astaxie/beego/orm"
	"github.com/cheneylew/goutil/utils"
)

type User struct {
	beego.BBaseUser
}

func init() {
	utils.JJKPrintln("==========regist models============")
	orm.RegisterModel(new(User))
}
