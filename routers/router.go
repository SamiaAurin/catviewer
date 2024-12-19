package routers

import (
	"catviewer/controllers"
	"github.com/beego/beego/v2/server/web"
)

func init() {
	// Route to show the vote page
	web.Router("/", &controllers.CatController{}, "get:ShowVotePage")

	// Route to handle vote casting
	web.Router("/cast_vote", &controllers.CatController{}, "post:CastVote")

	// Route to fetch a random cat image
	web.Router("/random_cat_image", &controllers.CatController{}, "get:FetchRandomImage")
    
    // Route to handle favoriting an image
	web.Router("/add_to_favorites", &controllers.CatController{}, "post:AddToFavorites")
}
