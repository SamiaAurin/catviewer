package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"catviewer/controllers"
	"github.com/beego/beego/v2/server/web/context" // Correct import for context
)

func TestCastVote(t *testing.T) {
	// Mock the vote request
	voteRequest := map[string]interface{}{
		"vote":     "1", // Upvote
		"image_id": "mock-image-id",
	}

	reqBody, err := json.Marshal(voteRequest)
	if err != nil {
		t.Fatalf("Error marshaling request body: %v", err)
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", "/cat/vote", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	// Use httptest to record the response
	rr := httptest.NewRecorder()

	// Create the Beego web context with the HTTP request and response writer
	ctx := &context.Context{
		ResponseWriter: rr, // ResponseWriter from httptest (correct type)
		Request:        req, // The request created above
	}

	// Create the controller instance and mock the context
	c := &controllers.CatController{}
	c.Ctx = ctx

	// Call the CastVote function (the function you're testing)
	c.CastVote()

	// Check the response status code
	if status := rr.Code; status != http.StatusFound {
		t.Errorf("Expected status %d, but got %d", http.StatusFound, status)
	}
}
