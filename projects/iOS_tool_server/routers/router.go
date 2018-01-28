package routers

import (
	"github.com/cheneylew/goutil/projects/iOS_tool_server/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
	beego.AutoRouter(&controllers.MainController{})
}
