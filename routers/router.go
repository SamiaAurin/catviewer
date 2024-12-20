package routers

import (
	"catviewer/controllers"
	"github.com/beego/beego/v2/server/web"
)

func init() {
	// Route to show the voting page
	web.Router("/cat/vote", &controllers.CatController{}, "get:ShowVotePage")
	
	// Route to handle voting (upvote/downvote)
	web.Router("/cat/vote", &controllers.CatController{}, "post:CastVote")
	
	// Route to handle saving favorites
	web.Router("/cat/favorite", &controllers.CatController{}, "post:SaveFavorite")

	// In router.go
    web.Router("/cat/favorites", &controllers.CatController{}, "get:GetFavorites")

}
