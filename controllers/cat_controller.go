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


// ShowVotePage renders the initial voting page with a random cat image.
func (c *CatController) ShowVotePage() {
	// Fetch a random cat image using a Go channel
	imageURL, imageID := fetchRandomImage()

	// Pass the image URL and ID to the template
	c.Data["ImageURL"] = imageURL
	c.Data["ImageID"] = imageID
	c.TplName = "catviewer.tpl"
}

// fetchRandomImage uses a Go channel to fetch a random image from TheCatAPI.
func fetchRandomImage() (string, string) {
	// Create a channel to communicate between the Go routine and the main thread
	ch := make(chan struct {
		url string
		id  string
		err error
	})

	// Start a Go routine to fetch the image asynchronously
	go func() {
		apiKey, err := web.AppConfig.String("catapi_key")
		if err != nil {
			ch <- struct {
				url string
				id  string
				err error
			}{err: fmt.Errorf("failed to read API key from config")}
			return
		}

		apiUrl := "https://api.thecatapi.com/v1/images/search"
		req, err := http.NewRequest("GET", apiUrl, nil)
		if err != nil {
			ch <- struct {
				url string
				id  string
				err error
			}{err: fmt.Errorf("error creating API request: %v", err)}
			return
		}
		req.Header.Set("x-api-key", apiKey)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			ch <- struct {
				url string
				id  string
				err error
			}{err: fmt.Errorf("failed to fetch image from API: %v", err)}
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			ch <- struct {
				url string
				id  string
				err error
			}{err: fmt.Errorf("failed to read API response: %v", err)}
			return
		}

		var images []map[string]interface{}
		err = json.Unmarshal(body, &images)
		if err != nil || len(images) == 0 {
			ch <- struct {
				url string
				id  string
				err error
			}{err: fmt.Errorf("failed to parse image response")}
			return
		}

		imageURL := images[0]["url"].(string)
		imageID := images[0]["id"].(string)
		ch <- struct {
			url string
			id  string
			err error
		}{url: imageURL, id: imageID}
	}()

	// Wait for the result from the Go routine
	result := <-ch
	if result.err != nil {
		// Log the error and render an error page
		fmt.Println(result.err)
		return "", ""
	}
	return result.url, result.id
}

// CastVote handles the voting action (upvote or downvote) for a cat image.
func (c *CatController) CastVote() {
	voteValue := c.GetString("vote")
	if voteValue != "1" && voteValue != "-1" {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{"error": "Invalid vote value"}
		c.ServeJSON()
		return
	}

	imageID := c.GetString("image_id")
	if imageID == "" {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{"error": "Missing image ID"}
		c.ServeJSON()
		return
	}

	// Call TheCatAPI to cast the vote
	go castVoteToAPI(imageID, voteValue)

	// Redirect to the same page after the vote
	c.Redirect("/cat/vote", http.StatusFound)
}

// castVoteToAPI sends the vote to TheCatAPI
func castVoteToAPI(imageID string, voteValue string) {
	apiKey, err := web.AppConfig.String("catapi_key")
	if err != nil {
		fmt.Println("Failed to read API key from config:", err)
		return
	}

	voteEndpoint := "https://api.thecatapi.com/v1/votes"
	voteRequest := map[string]interface{}{
		"image_id": imageID,
		"value":    voteValue,
	}

	reqBody, err := json.Marshal(voteRequest)
	if err != nil {
		fmt.Println("Failed to marshal vote request:", err)
		return
	}

	req, err := http.NewRequest("POST", voteEndpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Println("Failed to create request:", err)
		return
	}

	req.Header.Set("x-api-key", apiKey)
	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		fmt.Println("Failed to send vote request:", err)
	}
}

/////////////////////////////////
