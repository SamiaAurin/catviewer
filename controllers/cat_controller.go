package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
    "log"
	
	

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

// Define a Breed struct to parse the API response
type Breed struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Origin      string `json:"origin"`
	Description string `json:"description"`
	Image       struct {
		URL string `json:"url"`
	} `json:"image"`
	WikipediaURL string `json:"wikipedia_url"`
}


func (c *CatController) FetchBreeds() {
    // Get API key from the configuration
    apiKey, err := web.AppConfig.String("catapi_key")
    if err != nil || apiKey == "" {
        c.Data["json"] = map[string]string{"error": "API key not found"}
        c.ServeJSON()
        return
    }

    // Extract breed ID from query parameters
    breedId := c.GetString("id")
    if breedId == "" {
        // If no breed ID is provided, fetch all breeds
        fetchAllBreeds(apiKey, c)
    } else {
        // If breed ID is provided, fetch specific breed details and images
        fetchBreedWithImages(breedId, apiKey, c)
    }
}

// Fetch all breeds
func fetchAllBreeds(apiKey string, c *CatController) {
    dataChannel := make(chan interface{})
    errorChannel := make(chan error)

    // Fetch all breeds details concurrently
    go fetchAllBreedDetails(apiKey, dataChannel, errorChannel)

    select {
    case breeds := <-dataChannel:
        c.Data["json"] = breeds
		//log.Println("Fetched breeds:", breeds)
    case err := <-errorChannel:
        c.Data["json"] = map[string]string{"error": err.Error()}
        c.ServeJSON()
        return
    }

    c.ServeJSON()
}

// Fetch details for all breeds
func fetchAllBreedDetails(apiKey string, dataChannel chan interface{}, errorChannel chan error) {
    apiURL := "https://api.thecatapi.com/v1/breeds"
    client := &http.Client{Timeout: 10 * time.Second}
    req, _ := http.NewRequest("GET", apiURL, nil)
    req.Header.Add("x-api-key", apiKey)

    resp, err := client.Do(req)
    if err != nil {
        errorChannel <- err
        return
    }
    defer resp.Body.Close()

    body, _ := ioutil.ReadAll(resp.Body)
    var breeds []Breed
    err = json.Unmarshal(body, &breeds)
    if err != nil {
        errorChannel <- err
        return
    }

    dataChannel <- breeds
}

// Fetch specific breed with images
func fetchBreedWithImages(breedId, apiKey string, c *CatController) {
    dataChannel := make(chan interface{})
    imageChannel := make(chan interface{})
    errorChannel := make(chan error)

    // Fetch breed details concurrently
    go fetchBreedDetails(breedId, apiKey, dataChannel, errorChannel)
    // Fetch breed images concurrently
    go fetchBreedImages(breedId, apiKey, imageChannel, errorChannel)

    var breedDetails interface{}
    var breedImages interface{}
    var err error

    select {
    case breedDetails = <-dataChannel:
    case err = <-errorChannel:
        c.Data["json"] = map[string]string{"error": err.Error()}
        c.ServeJSON()
        return
    }

    select {
    case breedImages = <-imageChannel:
    case err = <-errorChannel:
        c.Data["json"] = map[string]string{"error": err.Error()}
        c.ServeJSON()
        return
    }

    c.Data["json"] = map[string]interface{}{
        "BreedDetails": breedDetails,
        "BreedImages":  breedImages,
    }
    c.ServeJSON()
}

// Fetch breed details
func fetchBreedDetails(breedId, apiKey string, dataChannel chan interface{}, errorChannel chan error) {
    apiURL := fmt.Sprintf("https://api.thecatapi.com/v1/breeds/%s", breedId)
    client := &http.Client{Timeout: 10 * time.Second}
    req, _ := http.NewRequest("GET", apiURL, nil)
    req.Header.Add("x-api-key", apiKey)

    resp, err := client.Do(req)
    if err != nil {
        errorChannel <- err
        return
    }
    defer resp.Body.Close()

    body, _ := ioutil.ReadAll(resp.Body)
    var breedDetails Breed
    err = json.Unmarshal(body, &breedDetails)
    if err != nil {
        errorChannel <- err
        return
    }

    dataChannel <- breedDetails
}

// Fetch breed images
func fetchBreedImages(breedId, apiKey string, imageChannel chan interface{}, errorChannel chan error) {
    apiURL := fmt.Sprintf("https://api.thecatapi.com/v1/images/search?limit=10&breed_ids=%s", breedId)
    client := &http.Client{Timeout: 10 * time.Second}
    req, _ := http.NewRequest("GET", apiURL, nil)
    req.Header.Add("x-api-key", apiKey)

    resp, err := client.Do(req)
    if err != nil {
        errorChannel <- err
        return
    }
    defer resp.Body.Close()

    body, _ := ioutil.ReadAll(resp.Body)
    var breedImages []struct {
        URL string `json:"url"`
    }
    err = json.Unmarshal(body, &breedImages)
    if err != nil {
        errorChannel <- err
        return
    }

    imageChannel <- breedImages
}

///////////////////////////////// BREEDS ENDS /////////////////////////////////////

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

// DeleteFavoriteImage handles the deletion of a favorite image
func (c *CatController) DeleteFavoriteImage() {
	// Retrieve the API key from the app.conf file
	apiKey, err := web.AppConfig.String("catapi_key")
	if err != nil {
		log.Printf("Error retrieving API key: %v\n", err)
		c.Data["json"] = map[string]string{"message": "API key not found."}
		c.ServeJSON()
		return
	}

	favoriteId := c.Ctx.Input.Param(":id") // Get the favorite ID from the URL parameter

	// Debug: Log the favorite ID received
	log.Printf("Received request to delete favorite image with ID: %s\n", favoriteId)

	// Create a channel to get the result of the deletion operation
	resultChan := make(chan bool)

	// Use a goroutine to handle the deletion in the background
	go func() {
		defer close(resultChan) // Close the channel once the goroutine is done

		// Make the DELETE request to The Cat API
		apiUrl := fmt.Sprintf("https://api.thecatapi.com/v1/favourites/%s", favoriteId)
		req, err := http.NewRequest("DELETE", apiUrl, nil)
		if err != nil {
			log.Printf("Error creating request: %v\n", err)
			resultChan <- false
			return
		}

		// Log the API URL and method being used
		log.Printf("Sending DELETE request to: %s\n", apiUrl)

		// Set the necessary headers with the API key from the app.conf file
		req.Header.Set("x-api-key", apiKey)

		// Make the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Error sending DELETE request: %v\n", err)
			resultChan <- false
			return
		}
		defer resp.Body.Close()

		// Log the response status code
		log.Printf("Received response with status code: %d\n", resp.StatusCode)

		// Read the response body for debugging purposes
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading response body: %v\n", err)
			resultChan <- false
			return
		}
		log.Printf("Response Body: %s\n", string(body))

		// Check if the request was successful
		if resp.StatusCode == 200 {
			log.Printf("Favorite image with ID %s deleted successfully.\n", favoriteId)
			resultChan <- true
		} else {
			log.Printf("Failed to delete favorite image with ID %s. Status code: %d\n", favoriteId, resp.StatusCode)
			resultChan <- false
		}
	}()

	// Wait for the result from the goroutine
	success := <-resultChan

	// Respond based on the result
	if success {
		c.Data["json"] = map[string]string{"message": "Favorite deleted successfully!"}
	} else {
		c.Data["json"] = map[string]string{"message": "Failed to delete favorite."}
	}

	c.ServeJSON()
}
///////////////////////////////// FAVS ENDS ////////////////////////////////////