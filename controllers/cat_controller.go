package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"github.com/beego/beego/v2/server/web"
)

type CatController struct {
	web.Controller
}

// ShowVotePage renders the initial voting page
func (c *CatController) ShowVotePage() {
	c.TplName = "catviewer.tpl" 
}

type VoteRequest struct {
	ImageID string `json:"image_id"`
	SubID   string `json:"sub_id,omitempty"`
	Value   int    `json:"value"`
}

type VoteResponse struct {
	StatusCode int
	Body       string
	Err        error
}
type FavoriteRequest struct {
    ImageID string `json:"image_id"`
    SubID   string `json:"sub_id"`
}

// CastVote handles the vote casting via the Cat API
func (c *CatController) CastVote() {
	// Parse the incoming JSON request body
	var voteReq VoteRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &voteReq); err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{"error": "Invalid request payload"}
		c.ServeJSON()
		return
	}

	// Get API key and URL from Beego configuration
	apiKey, err := web.AppConfig.String("catapi_key")
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Failed to read API key from config"}
		c.ServeJSON()
		return
	}

	catApiUrl, err := web.AppConfig.String("catapi_url")
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Failed to read Cat API URL from config"}
		c.ServeJSON()
		return
	}

	voteEndpoint := catApiUrl + "/votes"

	// Create a channel for the response
	respChan := make(chan VoteResponse)

	// Send the request asynchronously using Go routine
	go func() {
		client := &http.Client{}
		reqBody, _ := json.Marshal(voteReq)
		req, err := http.NewRequest("POST", voteEndpoint, bytes.NewBuffer(reqBody))
		if err != nil {
			respChan <- VoteResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
			return
		}
		req.Header.Set("x-api-key", apiKey)
		req.Header.Set("Content-Type", "application/json")

		// Make the API request
		resp, err := client.Do(req)
		if err != nil {
			respChan <- VoteResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			respChan <- VoteResponse{
				StatusCode: http.StatusInternalServerError,
				Err:        err,
			}
			return
		}

		respChan <- VoteResponse{
			StatusCode: resp.StatusCode,
			Body:       string(body),
			Err:        nil,
		}
	}()

	// Wait for the response from the Go routine
	resp := <-respChan

	// Check for errors or non-200 status codes
	if resp.Err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Failed to make API call"}
		c.ServeJSON()
		return
	}

	// Check if the voting was successful
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		c.Ctx.Output.SetStatus(resp.StatusCode)
		c.Data["json"] = map[string]string{"error": "Failed to cast vote"}
		c.ServeJSON()
		return
	}

	// Send a success response
	c.Data["json"] = map[string]string{"message": "Vote successfully cast!"}
	c.ServeJSON()
}

// FetchRandomImage fetches a random cat image from the API
func (c *CatController) FetchRandomImage() {
	apiKey, err := web.AppConfig.String("catapi_key")
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Failed to read API key from config"}
		c.ServeJSON()
		return
	}

	apiUrl := "https://api.thecatapi.com/v1/images/search"
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Error creating API request"}
		c.ServeJSON()
		return
	}
	req.Header.Set("x-api-key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Failed to fetch image from API"}
		c.ServeJSON()
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Failed to read API response"}
		c.ServeJSON()
		return
	}

	var images []map[string]interface{}
	err = json.Unmarshal(body, &images)
	if err != nil || len(images) == 0 {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Failed to parse image response"}
		c.ServeJSON()
		return
	}

	c.Data["json"] = map[string]interface{}{
		"url": images[0]["url"].(string),
		"id":  images[0]["id"].(string),
	}
	c.ServeJSON()
}

// AddToFavorites adds an image to the user's favorites
func (c *CatController) AddToFavorites() {
	var req FavoriteRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{"error": "Invalid request body"}
		c.ServeJSON()
		return
	}

	apiKey, err := web.AppConfig.String("catapi_key")
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Failed to read API key from config"}
		c.ServeJSON()
		return
	}

	// Corrected line here
	rawBody := fmt.Sprintf(`{"image_id": "%s", "sub_id": "%s"}`, req.ImageID, req.SubID)

	apiUrl := "https://api.thecatapi.com/v1/favourites"
	request, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer([]byte(rawBody)))
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Error creating API request"}
		c.ServeJSON()
		return
	}

	request.Header.Set("x-api-key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Failed to add image to favorites"}
		c.ServeJSON()
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Failed to read response"}
		c.ServeJSON()
		return
	}

	// Return the response from the Cat API
	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = map[string]string{"error": "Failed to parse response"}
		c.ServeJSON()
		return
	}

	c.Data["json"] = response
	c.ServeJSON()
}
