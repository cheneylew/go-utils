package routers

import (
	"github.com/cheneylew/goutil/stock_web_server/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
	beego.AutoRouter(&controllers.MainController{})
	beego.AutoRouter(&controllers.PhoneController{})
}
