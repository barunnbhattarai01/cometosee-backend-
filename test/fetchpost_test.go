package test

import (
	"cometosee/config"
	"cometosee/di"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchPost(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/posts", nil)
	rec := httptest.NewRecorder()

	//di
	config.InitCache()
	postcontroller := di.SetupPostController()
	postcontroller.FetchFeed(rec, req)
	t.Log("fetching posts")
	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("excepted status 200,got %d", res.StatusCode)
	}

	var body map[string]interface{}
	err := json.NewDecoder(res.Body).Decode(&body)
	if err != nil {
		t.Fatalf("failed to decode response : %v", err)
	}

	if body["message"] != "fetch sucessfully" {
		t.Errorf("unexpected message :%v", body["message"])
	}

	if _, ok := body["posts"]; !ok {
		t.Errorf("posts field missing in response")
	}
	t.Logf("post fetch sucessfully")
}
