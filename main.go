package main

import (
	_ "catviewer/routers"
	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	// Explicitly set the views path
	beego.SetViewsPath("views")

	beego.Run()
}
