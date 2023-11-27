package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)
import "ncc/db"
import "ncc/types"

func TestOptionsHandler(t *testing.T) {
	expected := ""
	req := httptest.NewRequest(http.MethodOptions, "/", nil)
	w := httptest.NewRecorder()
	OptionsRequest(w, req)
	res := w.Result()
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)

	if err != nil {
		t.Errorf("Error: %v", err)
	}

	if string(data) != expected {
		t.Errorf("Expected empty response but got %v", string(data))
	}

}

func TestGetPostCommentCount(t *testing.T) {
	var response map[string]interface{}
	req := httptest.NewRequest(http.MethodGet, "/commentcount/example-post", nil)
	w := httptest.NewRecorder()
	GetPostCommentCount(w, req)
	res := w.Result()
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	err = json.Unmarshal([]byte(data), &response)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	if response["status"] != "ok" {
		t.Errorf("Expected ok response but got %v", string(data))
	}
}

func TestGetPostComments(t *testing.T) {
	var response map[string]interface{}
	req := httptest.NewRequest(http.MethodGet, "/comments/example-post", nil)
	w := httptest.NewRecorder()
	GetPostComments(w, req)
	res := w.Result()
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	err = json.Unmarshal([]byte(data), &response)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	if response["status"] != "ok" {
		t.Errorf("Expected ok response but got %v", string(data))
	}
}

func TestFlagComment(t *testing.T) {
	t.Run("A comment which exists", func(t *testing.T) {
		comment := types.Comment{DisplayName: "author", Body: "comment body here", PostRef: "test_post_ref"}
		id, _ := db.AddComment(comment)
		_ = db.FlagComment(id, "0.0.0.0")

		flag_url := fmt.Sprintf("/flag/%d", id)

		req := httptest.NewRequest(http.MethodPost, flag_url, nil)
		req = mux.SetURLVars(req, map[string]string{"comment_id": fmt.Sprintf("%d", id)})
		w := httptest.NewRecorder()
		FlagComment(w, req)
		res := w.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusNoContent {
			t.Errorf("Got unexpected HTTP status")
		}

		response, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Errorf("Error: %v", err)
		}

		if len(response) != 0 {
			t.Errorf("Expected empty response but got %v", string(response[:]))
		}
	})
}

func TestAddComment(t *testing.T) {
	comment := types.Comment{DisplayName: "author", Body: "comment body here", PostRef: "test_post_ref"}
	data, err := json.Marshal(comment)
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/comments/%s", comment.PostRef), strings.NewReader(string(data[:])))
	req.Header.Add("Content-Type", "application/json")

	w := httptest.NewRecorder()
	AddComment(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		t.Errorf("Got unexpected HTTP status: %v", res.StatusCode)
	}

	response, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	if len(response) != 0 {
		t.Errorf("Expected empty response but got %v", string(response[:]))
	}
}
