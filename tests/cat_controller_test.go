package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"catviewer/controllers"
	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"

	//"github.com/stretchr/testify/assert"
	"fmt"
	"strings"
    //"io/ioutil"
	"github.com/stretchr/testify/assert"
	//"github.com/stretchr/testify/mock"

	"github.com/jarcoal/httpmock"
	//"github.com/astaxie/beego"
    //"github.com/beego/beego/v2/core/config"

	//"log"
	//"io"
	//"os"
    
	
)

/////////////// TESTS FOR VOTES STARTS ///////////////////

func TestShowVotePage(t *testing.T) {
	r, _ := http.NewRequest("GET", "/cat/vote", nil)
	w := httptest.NewRecorder()
	
	ctx := context.NewContext()
	ctx.Reset(w, r)
	
	controller := &controllers.CatController{}
	controller.Init(ctx, "CatController", "CatController", nil)
	
	controller.ShowVotePage()
	
	if w.Code != http.StatusOK {
		t.Errorf("ShowVotePage returned wrong status code: got %v want %v", w.Code, http.StatusOK)
	}
}

func TestCastVote(t *testing.T) {
	voteData := map[string]string{
		"image_id": "test123",
		"vote":     "1",
	}
	jsonData, _ := json.Marshal(voteData)
	
	r, _ := http.NewRequest("POST", "/cat/vote", bytes.NewBuffer(jsonData))
	w := httptest.NewRecorder()
	
	ctx := context.NewContext()
	ctx.Reset(w, r)
	
	controller := &controllers.CatController{}
	controller.Init(ctx, "CatController", "CatController", nil)
	
	controller.Ctx.Input.SetParam("image_id", "test123")
	controller.Ctx.Input.SetParam("vote", "1")
	
	controller.CastVote()
	
	if w.Code != http.StatusFound {
		t.Errorf("CastVote returned wrong status code: got %v want %v", w.Code, http.StatusFound)
	}
}

func TestCastVote_MissingImageID(t *testing.T) {
	
	voteData := map[string]string{
		"vote": "1",
	}
	jsonData, _ := json.Marshal(voteData)

	// Create the request and recorder
	r, _ := http.NewRequest("POST", "/cat/vote", bytes.NewBuffer(jsonData))
	w := httptest.NewRecorder()

	// Create the context and initialize the controller
	ctx := context.NewContext()
	ctx.Reset(w, r)
	controller := &controllers.CatController{}
	controller.Init(ctx, "CatController", "CatController", nil)

	controller.Ctx.Input.SetParam("vote", "1")
	controller.CastVote()

	// Check if the status code is 400 (Bad Request)
	if w.Code != http.StatusBadRequest {
		t.Errorf("CastVote returned wrong status code: got %v want %v", w.Code, http.StatusBadRequest)
	}

	// Check if the response body contains the expected error message
	expectedResponse := `{"error":"MissingimageID"}`
	actualResponse := w.Body.String()
	actualResponse = strings.ReplaceAll(actualResponse, "\n", "")
	actualResponse = strings.ReplaceAll(actualResponse, " ", "")
	
	if actualResponse != expectedResponse {
		t.Errorf("CastVote returned wrong body: got %v want %v", actualResponse, expectedResponse)
	}
}

func TestCastVoteInvalidValue(t *testing.T) {
	voteData := map[string]string{
		"image_id": "test123",
		"vote":     "0",
	}
	jsonData, _ := json.Marshal(voteData)
	
	r, _ := http.NewRequest("POST", "/cat/vote", bytes.NewBuffer(jsonData))
	w := httptest.NewRecorder()
	
	ctx := context.NewContext()
	ctx.Reset(w, r)
	
	controller := &controllers.CatController{}
	controller.Init(ctx, "CatController", "CatController", nil)
	
	controller.Ctx.Input.SetParam("image_id", "test123")
	controller.Ctx.Input.SetParam("vote", "0")
	
	controller.CastVote()
	
	if w.Code != http.StatusBadRequest {
		t.Errorf("CastVote returned wrong status code: got %v want %v", w.Code, http.StatusFound)
	}
}

func TestShowVotedImages(t *testing.T) {
	// Mock response data
	mockResponse := []map[string]interface{}{
		{
			"id": 123,
			"image_id": "test123",
			"value": 1,
		},
	}
	
	jsonResponse, _ := json.Marshal(mockResponse)
	
	// Create test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
	}))
	defer ts.Close()
	
	r, _ := http.NewRequest("GET", "/cat/voted_pics", nil)
	w := httptest.NewRecorder()
	
	ctx := context.NewContext()
	ctx.Reset(w, r)
	
	controller := &controllers.CatController{}
	controller.Init(ctx, "CatController", "CatController", nil)
	
	controller.ShowVotedImages()
	
	if w.Code != http.StatusOK {
		t.Errorf("ShowVotedImages returned wrong status code: got %v want %v", w.Code, http.StatusOK)
	}
}

/////////////// TESTS FOR VOTES ENDS ///////////////////

/////////////// TESTS FOR BREEDS STARTS ///////////////////
func setupMockConfig() {
    web.AppConfig.Set("catapi_key", "mock_api_key")
}

// Mock breed data
var mockBreed = controllers.Breed{
	ID:          "abys",
	Name:        "Abyssinian",
	Origin:      "Egypt",
	Description: "The Abyssinian is easy to care for, and a joy to have in your home.",
	WikipediaURL: "https://en.wikipedia.org/wiki/Abyssinian_(cat)",
}

// Mock breed list
var mockBreeds = []controllers.Breed{mockBreed}

var fetchAllBreedDetailsURL string
var fetchBreedDetailsURL string
var fetchBreedImagesURL string

// Testing to fetch all the breeds
func TestFetchBreedsAll(t *testing.T) {
    setupMockConfig() // Ensure mock configuration is set

    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path != "/v1/breeds" {
            t.Errorf("Expected to request '/v1/breeds', got: %s", r.URL.Path)
        }
        if r.Header.Get("x-api-key") == "" {
            t.Error("Expected API Key header to be set")
        }

        // Mock response
        json.NewEncoder(w).Encode(mockBreeds)
    }))
    defer ts.Close()

    // Override API endpoint
    oldURL := fetchAllBreedDetailsURL
    fetchAllBreedDetailsURL = ts.URL
    defer func() { fetchAllBreedDetailsURL = oldURL }()

    r, _ := http.NewRequest("GET", "/cat/fetch_breeds", nil)
    w := httptest.NewRecorder()

    ctx := context.NewContext()
    ctx.Reset(w, r)

    controller := &controllers.CatController{}
    controller.Init(ctx, "CatController", "CatController", nil)

    controller.FetchBreeds()

    // Check response
    if w.Code != http.StatusOK {
        t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
    }

    var response []controllers.Breed
    err := json.Unmarshal(w.Body.Bytes(), &response)
    if err != nil {
        t.Errorf("Failed to unmarshal response: %v", err)
    }

    if len(response) == 0 {
        t.Error("Expected non-empty breed list in response")
    }
}
// Testing Breeds With Images
func TestFetchBreedWithImages(t *testing.T) {
    setupMockConfig() // Ensure mock configuration is set

    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/breeds/abys" {
			json.NewEncoder(w).Encode(mockBreed)
		} else if r.URL.Path == "/v1/images/search" {
			mockImages := []struct {
				URL string `json:"url"`
			}{
				{URL: "http://example.com/cat1.jpg"},
				{URL: "http://example.com/cat2.jpg"},
			}
			json.NewEncoder(w).Encode(mockImages)
		} else {
			t.Errorf("Unexpected request to path: %s", r.URL.Path)
		}
	}))
	defer ts.Close()
	
	fetchBreedDetailsURL = ts.URL
	fetchBreedImagesURL = ts.URL

    // Override API endpoints
    oldBreedDetailsURL := fetchBreedDetailsURL
    oldBreedImagesURL := fetchBreedImagesURL
    fetchBreedDetailsURL = ts.URL
    fetchBreedImagesURL = ts.URL
    defer func() {
        fetchBreedDetailsURL = oldBreedDetailsURL
        fetchBreedImagesURL = oldBreedImagesURL
    }()

    r, _ := http.NewRequest("GET", "/cat/fetch_breeds?id=abys", nil)
    w := httptest.NewRecorder()

    ctx := context.NewContext()
    ctx.Reset(w, r)

    controller := &controllers.CatController{}
    controller.Init(ctx, "CatController", "CatController", nil)
    controller.Ctx.Input.SetParam("id", "abys")

    controller.FetchBreeds()

    // Check response
    if w.Code != http.StatusOK {
        t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
    }

    var response map[string]interface{}
    err := json.Unmarshal(w.Body.Bytes(), &response)
    if err != nil {
        t.Errorf("Failed to unmarshal response: %v", err)
    }

    if response["BreedDetails"] == nil {
        t.Error("Expected BreedDetails in response")
    }
    if response["BreedImages"] == nil {
        t.Error("Expected BreedImages in response")
    }
}


// TestFetchBreedsError tests error handling when API key is missing
func TestFetchBreedsError(t *testing.T) {
	r, _ := http.NewRequest("GET", "/cat/fetch_breeds", nil)
	w := httptest.NewRecorder()
	
	ctx := context.NewContext()
	ctx.Reset(w, r)
	
	controller := &controllers.CatController{}
	controller.Init(ctx, "CatController", "CatController", nil)
	
	// Ensure API key is not set
	web.AppConfig.Set("catapi_key", "")
	
	controller.FetchBreeds()
	
	// Check error response
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}
	
	if response["error"] != "API key not found" {
		t.Errorf("Expected error message 'API key not found', got %s", response["error"])
	}
}

/////////////// TESTS FOR BREEDS ENDS /////////////////////

/////////////// TESTS FOR FAVS STARTS /////////////////////

// Favorite Img Post Test
var favEndpoint string

func TestFavoriteImage(t *testing.T) {
    tests := []struct {
        name           string
        imageID        string
        mockAPIStatus  int
        mockAPIResponse string
        expectedStatus int
        expectError    bool
    }{
        {
            name:           "Successful favorite",
            imageID:        "test123",
            mockAPIStatus:  http.StatusCreated,
            mockAPIResponse: `{"id": "fav_123"}`,
            expectedStatus: http.StatusFound,
            expectError:    false,
        },
        {
            name:           "Empty image ID",
            imageID:        "",
            mockAPIStatus:  http.StatusOK,
            mockAPIResponse: "",
            expectedStatus: http.StatusBadRequest,
            expectError:    true,
        },
        {
            name:           "API error response",
            imageID:        "test123",
            mockAPIStatus:  http.StatusUnauthorized,
            mockAPIResponse: `{"message": "Invalid API key"}`,
            expectedStatus: http.StatusInternalServerError,
            expectError:    true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Create mock server
            ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                // Verify request method and headers
                if r.Method != "POST" {
                    t.Errorf("Expected POST request, got %s", r.Method)
                }
                if r.Header.Get("Content-Type") != "application/json" {
                    t.Errorf("Expected Content-Type: application/json, got %s", r.Header.Get("Content-Type"))
                }
                if r.Header.Get("x-api-key") == "" {
                    t.Error("Expected x-api-key header to be set")
                }

                // Send mock response
                w.WriteHeader(tt.mockAPIStatus)
                w.Write([]byte(tt.mockAPIResponse))
            }))
            defer ts.Close()

            // Override the API endpoint for testing
            originalEndpoint := favEndpoint
            controllers.SetFavEndpoint(ts.URL)  // Use the SetFavEndpoint function
            defer controllers.SetFavEndpoint(originalEndpoint)  // Restore the original endpoint after the test

            // Set mock API key in config
            web.AppConfig.Set("catapi_key", "mock_api_key")

            // Create test request
            r := httptest.NewRequest("POST", "/cat/favorite", nil)
            w := httptest.NewRecorder()

            // Setup form data
            r.ParseForm()
            r.Form.Set("image_id", tt.imageID)

            // Setup controller
            ctx := context.NewContext()
            ctx.Reset(w, r)
            controller := &controllers.CatController{}
            controller.Init(ctx, "CatController", "CatController", nil)

            // Execute
            controller.FavoriteImage()

            if w.Code != tt.expectedStatus {
                t.Errorf("Expected status %v, got %v", tt.expectedStatus, w.Code)
            }

            // For error cases, verify error message
            if tt.expectError {
                var response map[string]string
                err := json.Unmarshal(w.Body.Bytes(), &response)
                if err != nil {
                    t.Fatalf("Failed to unmarshal response: %v", err)
                }
                if response["error"] == "" {
                    t.Error("Expected error message in response")
                }
            }
        })
    }
}

// Show the FavoriteImages Test
var fetchFavoriteImagesURL string

func TestShowFavoriteImages(t *testing.T) {
	// Mock response data
	mockResponse := []map[string]interface{}{
		{
			"id": 1,
			"image_id": "test_image_id",
		},
	}

	jsonResponse, _ := json.Marshal(mockResponse)

	// Create test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
	}))
	defer ts.Close()

	// Override API URL
	oldURL := fetchFavoriteImagesURL
	fetchFavoriteImagesURL = ts.URL
	defer func() { fetchFavoriteImagesURL = oldURL }()

	r, _ := http.NewRequest("GET", "/cat/favorites", nil)
	w := httptest.NewRecorder()

	ctx := context.NewContext()
	ctx.Reset(w, r)

	controller := &controllers.CatController{}
	controller.Init(ctx, "CatController", "CatController", nil)

	controller.ShowFavoriteImages()

	// Check response status code
	if w.Code != http.StatusOK {
		t.Errorf("ShowFavoriteImages returned wrong status code: got %v want %v", w.Code, http.StatusOK)
	}
}

// DeleteFavoriteImage Test
func TestDeleteFavoriteImage(t *testing.T) {
	// Mock favorite ID
	favoriteID := "1"

	// Create a new request
	r, _ := http.NewRequest("DELETE", "/cat/favorites/"+favoriteID, nil)
	w := httptest.NewRecorder()

	ctx := context.NewContext()
	ctx.Reset(w, r)

	controller := &controllers.CatController{}
	controller.Init(ctx, "CatController", "CatController", nil)

	controller.DeleteFavoriteImage()

	// Check response status code
	if w.Code != http.StatusOK {
		t.Errorf("DeleteFavoriteImage returned wrong status code: got %v want %v", w.Code, http.StatusOK)
	}
}

// TestFavoriteImageAPIFailure tests the API request failure scenarios
func TestFavoriteImageAPIFailure(t *testing.T) {
    // Set mock API key
    web.AppConfig.Set("catapi_key", "mock_api_key")

    // Test case for network failure
    t.Run("Network failure", func(t *testing.T) {
        // Use an invalid URL to simulate network failure
        favEndpoint = "http://invalid-url"

        r := httptest.NewRequest("POST", "/cat/favorite", nil)
        w := httptest.NewRecorder()

        r.ParseForm()
        r.Form.Set("image_id", "test123")

        ctx := context.NewContext()
        ctx.Reset(w, r)
        controller := &controllers.CatController{}
        controller.Init(ctx, "CatController", "CatController", nil)

        controller.FavoriteImage()

        if w.Code != http.StatusInternalServerError {
            t.Errorf("Expected status 500, got %v", w.Code)
        }
    })
}

/////////////// TESTS FOR FAVS ENDS //////////////////////

/////////////// TESTS FetchImage functions STARTS //////////////////////

func TestFetchRandomImage(t *testing.T) {
	// Initialize httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Mock the API URL
	apiUrl := "https://api.thecatapi.com/v1/images/search"
	httpmock.RegisterResponder("GET", apiUrl, httpmock.NewStringResponder(200, `[{"id":"abc123", "url":"https://example.com/cat.jpg"}]`))

	// Call the FetchRandomImage function
	imageURL, imageID := controllers.FetchRandomImage()

	// Assertions
	assert.NoError(t, nil) // This is a simple check, but can be expanded based on error handling improvements.
	assert.Equal(t, "https://example.com/cat.jpg", imageURL)
	assert.Equal(t, "abc123", imageID)
}

func TestFetchRandomImage_ErrorInResponseFetch(t *testing.T) {
	// Initialize httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Mock the API URL
	apiUrl := "https://api.thecatapi.com/v1/images/search"
	httpmock.RegisterResponder("GET", apiUrl, func(req *http.Request) (*http.Response, error) {
		// Simulate API fetch failure
		return nil, fmt.Errorf("failed to fetch from API")
	})

	// Call the FetchRandomImage function
	imageURL, imageID := controllers.FetchRandomImage()

	// Assertions
	assert.Empty(t, imageURL)
	assert.Empty(t, imageID)
}

func TestFetchRandomImage_ErrorInReadingBody(t *testing.T) {
	// Initialize httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Mock the API URL with a response
	apiUrl := "https://api.thecatapi.com/v1/images/search"
	httpmock.RegisterResponder("GET", apiUrl, httpmock.NewStringResponder(200, "invalid json"))

	// Call the FetchRandomImage function
	imageURL, imageID := controllers.FetchRandomImage()

	// Assertions
	assert.Empty(t, imageURL)
	assert.Empty(t, imageID)
}

func TestFetchRandomImage_ErrorInJSONParsing(t *testing.T) {
	// Initialize httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Mock the API URL with an invalid JSON response
	apiUrl := "https://api.thecatapi.com/v1/images/search"
	httpmock.RegisterResponder("GET", apiUrl, httpmock.NewStringResponder(200, `{"id":"abc123"}`)) // Missing `url` field

	// Call the FetchRandomImage function
	imageURL, imageID := controllers.FetchRandomImage()

	// Assertions
	assert.Empty(t, imageURL)
	assert.Empty(t, imageID)
}

/////////////// TESTS FetchImage functions ENDS //////////////////////

////////////////////////////////////////////////////////



//////////////////////////////////////////////////////////



func TestMain(m *testing.M) {
	web.BConfig.RunMode = web.DEV
	m.Run()
}