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
    // Set a mock API key
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

/*
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
*/
/////////////// TESTS FOR BREEDS ENDS /////////////////////

func TestMain(m *testing.M) {
	web.BConfig.RunMode = web.DEV
	m.Run()
}