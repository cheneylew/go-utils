package routers

import (
	"github.com/cheneylew/goutil/beego_web_demo/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
	beego.AutoRouter(&controllers.MainController{})
}
