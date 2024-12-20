package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
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

// SaveFavorite handles saving the image to favorites.
func (c *CatController) SaveFavorite() {
	imageID := c.GetString("image_id")
	if imageID == "" {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = map[string]string{"error": "Missing image ID"}
		c.ServeJSON()
		return
	}

	// Call TheCatAPI to save the favorite
	go saveFavoriteToAPI(imageID)

	// Redirect to the same page after saving
	c.Redirect("/cat/vote", http.StatusFound)

	
}

// saveFavoriteToAPI sends the favorite request to TheCatAPI, handling redirects.
func saveFavoriteToAPI(imageID string) {
	apiKey, err := web.AppConfig.String("catapi_key")
	if err != nil {
		fmt.Println("Failed to read API key from config:", err)
		return
	}

	favoriteEndpoint := "https://api.thecatapi.com/v1/favourites"
	favoriteRequest := map[string]interface{}{
		"image_id": imageID,
	}

	reqBody, err := json.Marshal(favoriteRequest)
	if err != nil {
		fmt.Println("Failed to marshal favorite request:", err)
		return
	}

	req, err := http.NewRequest("POST", favoriteEndpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Println("Failed to create request:", err)
		return
	}

	req.Header.Set("x-api-key", apiKey)
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Failed to send favorite request:", err)
		return
	}
	defer resp.Body.Close()

	// Check if we got a 302 redirect response
	if resp.StatusCode == http.StatusFound { // 302
		location := resp.Header.Get("Location")
		fmt.Printf("Redirected to: %s\n", location)

		// Follow the redirect with a new POST request
		newReq, err := http.NewRequest("POST", location, bytes.NewBuffer(reqBody))
		if err != nil {
			fmt.Println("Failed to create new request:", err)
			return
		}
		newReq.Header.Set("x-api-key", apiKey)

		newResp, err := client.Do(newReq)
		if err != nil {
			fmt.Println("Failed to follow redirect:", err)
			return
		}
		defer newResp.Body.Close()

		// Check if the final response is successful
		if newResp.StatusCode == http.StatusOK || newResp.StatusCode == http.StatusCreated {
			fmt.Println("Favorite saved successfully after redirect!")
		} else {
			fmt.Println("Failed to save favorite after redirect. Status:", newResp.StatusCode)
		}
	} else if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		// Handle success on the first request
		fmt.Println("Favorite saved successfully!")
	} else {
		fmt.Printf("Unexpected response status: %d\n", resp.StatusCode)
	}
}


// GetFavorites retrieves the user's favorites from the Cat API
func (c *CatController) GetFavorites() {
    apiKey, err := web.AppConfig.String("catapi_key")
    if err != nil {
        log.Printf("Failed to read API key from config: %s", err.Error())
        c.Ctx.Output.SetStatus(http.StatusInternalServerError)
        c.Data["json"] = map[string]string{"error": "Failed to read API key from config"}
        c.ServeJSON()
        return
    }

    // Create a channel to receive API response
    resultChan := make(chan []byte)
    errorChan := make(chan error)

    // Use a Go routine to fetch favorites
    go func() {
        apiUrl := "https://api.thecatapi.com/v1/favourites"
        req, err := http.NewRequest("GET", apiUrl, nil)
        if err != nil {
            errorChan <- fmt.Errorf("error creating API request: %s", err.Error())
            return
        }

        req.Header.Set("x-api-key", apiKey)

        client := &http.Client{}
        resp, err := client.Do(req)
        if err != nil {
            errorChan <- fmt.Errorf("failed to fetch favorites from API: %s", err.Error())
            return
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusOK {
            errorChan <- fmt.Errorf("unexpected API response status: %d", resp.StatusCode)
            return
        }

        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            errorChan <- fmt.Errorf("error reading API response body: %s", err.Error())
            return
        }

        resultChan <- body
    }()

    // Wait for the response or error
    select {
    case body := <-resultChan:
        // Log the raw response for debugging
        log.Printf("Raw response: %s", string(body))

        // Parse the JSON response
        var favorites []struct {
            ID     int    `json:"id"`
            Image  struct {
                URL string `json:"url"`
            } `json:"image"`
        }

        err := json.Unmarshal(body, &favorites)
        if err != nil {
            log.Printf("Error decoding response: %s", err.Error())
            c.Ctx.Output.SetStatus(http.StatusInternalServerError)
            c.Data["json"] = map[string]string{"error": "Error decoding response"}
            c.ServeJSON()
            return
        }

        // Return JSON response
        c.Data["json"] = favorites
        c.ServeJSON()

    case err := <-errorChan:
        log.Printf("Error fetching favorites: %s", err.Error())
        c.Ctx.Output.SetStatus(http.StatusInternalServerError)
        c.Data["json"] = map[string]string{"error": err.Error()}
        c.ServeJSON()
    }
}


