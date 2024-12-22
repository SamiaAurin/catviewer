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

///////////////////////////////// VOTE STARTS ////////////////////////////////////
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
	// Return a JSON response
	//c.Data["json"] = map[string]string{"image_id": imageID, "vote_value": voteValue}
	//c.ServeJSON()
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

	// Channel to handle response and error
	responseChannel := make(chan string)
	errorChannel := make(chan error)

	go func() {
		reqBody, err := json.Marshal(voteRequest)
		if err != nil {
			errorChannel <- fmt.Errorf("failed to marshal vote request: %v", err)
			return
		}

		req, err := http.NewRequest("POST", voteEndpoint, bytes.NewBuffer(reqBody))
		if err != nil {
			errorChannel <- fmt.Errorf("failed to create request: %v", err)
			return
		}

		req.Header.Set("x-api-key", apiKey)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			errorChannel <- fmt.Errorf("failed to send vote request for image_id %s: %v", imageID, err)
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			errorChannel <- fmt.Errorf("failed to read API response for image_id %s: %v", imageID, err)
			return
		}

		if resp.StatusCode == http.StatusCreated {
			responseChannel <- fmt.Sprintf("Vote successfully posted for image_id %s with value %s. Response: %s", imageID, voteValue, string(body))
		} else {
			errorChannel <- fmt.Errorf("failed to post vote for image_id %s. Status: %d, Response: %s", imageID, resp.StatusCode, string(body))
		}
	}()

	// Wait for either a response or an error
	select {
	case response := <-responseChannel:
		fmt.Println(response)
	case err := <-errorChannel:
		fmt.Println("Error:", err)
	}
}


// ShowVotedImages handles fetching and displaying voted images.
func (c *CatController) ShowVotedImages() {
	apiKey, err := web.AppConfig.String("catapi_key")
	if err != nil {
		c.Data["json"] = map[string]string{"error": "Failed to read API key"}
		c.ServeJSON()
		return
	}

	// Channel for results and errors
	results := make(chan []map[string]interface{}, 1)
	errors := make(chan error, 1)

	// API URL to fetch voted images
	apiUrl := "https://api.thecatapi.com/v1/votes"

	// Goroutine to fetch data
	go func(apiUrl string, apiKey string) {
		req, err := http.NewRequest("GET", apiUrl, nil)
		if err != nil {
			errors <- err
			return
		}
		req.Header.Set("x-api-key", apiKey)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			errors <- err
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			errors <- err
			return
		}

		var votes []map[string]interface{}
		err = json.Unmarshal(body, &votes)
		if err != nil {
			errors <- err
			return
		}

		// Send the result to the results channel
		results <- votes
	}(apiUrl, apiKey)

	// Wait for either a result or an error
	select {
	case res := <-results:
		// Respond with the fetched votes
		c.Data["json"] = res
		c.ServeJSON()
	case err := <-errors:
		// Handle the error
		c.Data["json"] = map[string]string{"error": err.Error()}
		c.ServeJSON()
	}
}

///////////////////////////////// VOTE ENDS ////////////////////////////////////

///////////////////////////////// BREEDS STARTS ////////////////////////////////////
///////////////////////////////// BREEDS ENDS ////////////////////////////////////

///////////////////////////////// FAVS STARTS ////////////////////////////////////
// FavoriteImage handles the favoriting action for a cat image
func (c *CatController) FavoriteImage() {
    // Get image ID from the form
    imageID := c.GetString("image_id")
    if imageID == "" {
        c.Ctx.Output.SetStatus(http.StatusBadRequest)
        c.Data["json"] = map[string]string{"error": "Missing image ID"}
        c.ServeJSON()
        return
    }

    // Create a channel to receive the result of the API call
    ch := make(chan error)

    // Call the function to favorite the image in a goroutine
    go favoriteImageToAPI(imageID, ch)

    // Wait for the result (blocking until the goroutine completes)
    err := <-ch
    if err != nil {
        // Handle the error if the API request failed
        c.Ctx.Output.SetStatus(http.StatusInternalServerError)
        c.Data["json"] = map[string]string{"error": err.Error()}
        c.ServeJSON()
        return
    }

    // Redirect to the same page after successful operation
    c.Redirect("/cat/vote", http.StatusFound)
}

// Fav pics fetch API
func favoriteImageToAPI(imageID string, ch chan error) {
    apiKey, err := web.AppConfig.String("catapi_key")
    if err != nil {
        ch <- fmt.Errorf("failed to read API key from config: %v", err)
        return
    }

    favEndpoint := "https://api.thecatapi.com/v1/favourites"
    favRequest := map[string]interface{}{
        "image_id": imageID,
    }

    reqBody, err := json.Marshal(favRequest)
    if err != nil {
        ch <- fmt.Errorf("failed to marshal favorite request: %v", err)
        return
    }

    req, err := http.NewRequest("POST", favEndpoint, bytes.NewBuffer(reqBody))
    if err != nil {
        ch <- fmt.Errorf("failed to create request: %v", err)
        return
    }

    req.Header.Set("x-api-key", apiKey)
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        ch <- fmt.Errorf("failed to send favorite request for image_id %s: %v", imageID, err)
        return
    }
    defer resp.Body.Close()

    // Log response for debugging
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        ch <- fmt.Errorf("failed to read API response for image_id %s: %v", imageID, err)
        return
    }

    if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
        // Status 200 or 201 means success
        fmt.Printf("Image successfully favorited for image_id %s. Response: %s\n", imageID, string(body))
        ch <- nil // Indicate success
    } else {
        // Handle unexpected status codes
        ch <- fmt.Errorf("failed to favorite image_id %s. Status: %d, Response: %s", imageID, resp.StatusCode, string(body))
    }
}

func (c *CatController) ShowFavoriteImages() {
    apiKey, err := web.AppConfig.String("catapi_key")
    if err != nil {
        c.Data["json"] = map[string]string{"error": "Failed to read API key"}
        c.ServeJSON()
        return
    }

    // Channel for results and errors
    results := make(chan []map[string]interface{}, 1)
    errors := make(chan error, 1)

    // API URL to fetch favorite images
    apiUrl := "https://api.thecatapi.com/v1/favourites"

    // Goroutine to fetch data
    go func(apiUrl string, apiKey string) {
        req, err := http.NewRequest("GET", apiUrl, nil)
        if err != nil {
            errors <- err
            return
        }
        req.Header.Set("x-api-key", apiKey)

        client := &http.Client{}
        resp, err := client.Do(req)
        if err != nil {
            errors <- err
            return
        }
        defer resp.Body.Close()

        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            errors <- err
            return
        }

        var favs []map[string]interface{}
        err = json.Unmarshal(body, &favs)
        if err != nil {
            errors <- err
            return
        }

        // Send the result to the results channel
        results <- favs
    }(apiUrl, apiKey)

    // Wait for either a result or an error
    select {
    case res := <-results:
        // Respond with the fetched favorites
        c.Data["json"] = res
        c.ServeJSON()
    case err := <-errors:
        // Handle the error
        c.Data["json"] = map[string]string{"error": err.Error()}
        c.ServeJSON()
    }
}

///////////////////////////////// FAVS ENDS ////////////////////////////////////