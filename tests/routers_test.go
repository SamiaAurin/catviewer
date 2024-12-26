package tests

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/stretchr/testify/assert"
    "catviewer/routers"
    "github.com/beego/beego/v2/server/web"
)

func init() {
    // Ensure Beego looks for views in the correct directory for testing
    web.SetViewsPath("./views")
    web.BConfig.WebConfig.EnableDocs = false
}

func TestCatRoutes(t *testing.T) {
    // Initialize Beego's routers
    routers.Init()

    // Create a new request for the ShowVotePage route (GET /cat/vote)
    req, err := http.NewRequest("GET", "/cat/vote", nil)
    if err != nil {
        t.Fatal(err)
    }

    // Record the response
    rr := httptest.NewRecorder()

    // Serve the HTTP request through Beego's test server
    web.BeeApp.Handlers.ServeHTTP(rr, req)

    // Check if the response status code is OK (200)
    assert.Equal(t, http.StatusOK, rr.Code)

    // Add your other test cases for different routes here (POST, DELETE, etc.)
}
