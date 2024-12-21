package main

import (
	_ "catviewer/routers"
	"github.com/beego/beego/v2/server/web"
)

func main() {
	// Explicitly set the views path
	web.SetViewsPath("views")

	web.Run()
}
