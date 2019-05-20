package routers

import (
	"github.com/cheneylew/goutil/projects/stock_analyze/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
	beego.AutoRouter(&controllers.MainController{})
	beego.AutoRouter(&controllers.PhoneController{})
}
