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

    // Route for handling votes (Up or Down)
    // POST method to handle the vote
    beego.Router("/cats/vote", &controllers.CatController{}, "post:Vote")
}
