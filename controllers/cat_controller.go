package controllers

import (
	"encoding/json"
	"github.com/beego/beego/v2/server/web"
	"net/http"
	"bytes"
	"io/ioutil"
)

type CatController struct {
	web.Controller
}

// Handle the GET request for displaying the page
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

// Handle voting (Up or Down)
func (c *CatController) Vote() {
	// Get the vote data from the request body
	var voteData map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &voteData)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		c.Data["json"] = map[string]string{"error": "Invalid request data"}
		c.ServeJSON()
		return
	}

	// Extract vote parameters
	imageId := voteData["image_id"].(string)
	subId := voteData["sub_id"].(string)
	value := int(voteData["value"].(float64))

	// Prepare the vote request payload
	payload := map[string]interface{}{
		"image_id": imageId,
		"sub_id":   subId,
		"value":    value,
	}

	// Encode the payload to JSON
	payloadData, err := json.Marshal(payload)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Failed to encode vote data"}
		c.ServeJSON()
		return
	}

	// Get API key and URL from the config
	apiKey, _ := web.AppConfig.String("catapi_key")
	apiURL, _ := web.AppConfig.String("catapi_url")

	// Set the headers
	req, err := http.NewRequest("POST", apiURL+"/votes", bytes.NewBuffer(payloadData))
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Failed to create request"}
		c.ServeJSON()
		return
	}
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("Content-Type", "application/json")

	// Send the request to the Cat API
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Failed to send request to Cat API"}
		c.ServeJSON()
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Failed to read Cat API response"}
		c.ServeJSON()
		return
	}

	// Check if the vote was successful (adjust based on response structure)
	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Failed to parse Cat API response"}
		c.ServeJSON()
		return
	}

	// Check if the response indicates success (depending on Cat API's response structure)
	if _, ok := response["id"]; ok {
		c.Data["json"] = map[string]bool{"success": true}
	} else {
		c.Data["json"] = map[string]string{"error": "Vote not successful"}
	}

	c.ServeJSON()
}
