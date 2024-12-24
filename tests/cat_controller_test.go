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

/////////////// TESTS FOR BREEDS ENDS /////////////////////

func TestMain(m *testing.M) {
	web.BConfig.RunMode = web.DEV
	m.Run()
}