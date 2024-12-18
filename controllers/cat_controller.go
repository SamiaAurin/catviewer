package controllers

import (
    "github.com/beego/beego/v2/server/web"
)

type CatController struct {
    web.Controller
}

func (c *CatController) Get() {
    // Load API Key and URL from the config
    apiKey, _ := web.AppConfig.String("catapi_key")
    apiURL, _ := web.AppConfig.String("catapi_url")

    // Pass the data to the template
    c.Data["APIKey"] = apiKey
    c.Data["APIURL"] = apiURL

    // Set the template to render
    c.TplName = "catviewer.tpl"
	if err := c.Render(); err != nil {
		c.Ctx.WriteString("Error rendering template: " + err.Error())
	}
	
}
