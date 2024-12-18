package routers

import (
	"catviewer/controllers"
	beego "github.com/beego/beego/v2/server/web"
)

func init() {
    // Main route
    beego.Router("/", &controllers.MainController{})
    
    // Route for CatController
    beego.Router("/cats", &controllers.CatController{})
}