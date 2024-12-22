package routers

import (
	"catviewer/controllers"
	"github.com/beego/beego/v2/server/web"
)

func init() {
	// Route to show the voting page
	web.Router("/cat/vote", &controllers.CatController{}, "get:ShowVotePage")
	
	// Route to handle voting (upvote/downvote and showing voted images)
	web.Router("/cat/vote", &controllers.CatController{}, "post:CastVote")
	web.Router("/cat/voted_pics", &controllers.CatController{}, "get:ShowVotedImages")

	// Route to handle saving and showin favorites
    web.Router("/cat/favorite", &controllers.CatController{}, "post:FavoriteImage")
	web.Router("/cat/fav_pics", &controllers.CatController{}, "get:ShowFavoriteImages")



}