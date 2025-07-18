package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"runtime/debug"
	"strings"
	"testing"
)
import "ncc/db"
import "ncc/types"

func setupTest(tb testing.TB) func(tb testing.TB) {
	_, err := db.CreateBlog("testsite")
	if err != nil {
		tb.Logf("creating testsite entry failed. Stacktrace follows")
		debug.PrintStack()
		tb.Errorf("API setup - Could not create test blog: %s", err)
	}

	return func(tb testing.TB) {
		dbh, err := db.SetupDB()
		defer dbh.Close()
		if err != nil {
			tb.Errorf("got error when connecting to database: %s", err)
		}

		_, err = dbh.Query("TRUNCATE blogs CASCADE")
		if err != nil {
			tb.Errorf("Could not rollback database during testing")
		}
	}
}

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

func TestGetPostComments(t *testing.T) {
	teardown := setupTest(t)
	defer teardown(t)

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
	teardown := setupTest(t)
	defer teardown(t)

	t.Run("A comment which exists", func(t *testing.T) {
		comment := types.Comment{DisplayName: "author", Body: "comment body here", PostRef: "test_post_ref", BlogRef: "testsite"}

		id, err := db.AddComment(comment)
		if err != nil {
			t.Logf(fmt.Sprintf("Got error when creating a comment: %v", err))
			t.Errorf("Could not add a comment (in order to flag it)")
		}
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
	teardown := setupTest(t)
	defer teardown(t)
	comment := types.Comment{DisplayName: "author", Body: "comment body here", PostRef: "test_post_ref", BlogRef: "testsite"}
	data, err := json.Marshal(comment)
	t.Logf(fmt.Sprintf("API/TestAddComment: Going to post: %s", data))
	post_url := fmt.Sprintf("/v2/blog/%s/comments/%s", comment.BlogRef, comment.PostRef)
	t.Logf(fmt.Sprintf("API/TestAddComment: Post destination: %s", post_url))
	req, err := http.NewRequest(http.MethodPost, post_url, strings.NewReader(string(data[:])))
	req.Header.Add("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/v2/blog/{blogref}/comments/{postref}", AddComment)
	router.ServeHTTP(w, req)
	t.Logf("Posted")

	if w.Code != http.StatusNoContent {
		t.Errorf("API/TestaAddComment: Got unexpected HTTP status: %v", w.Code)
	}

	response, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	if len(response) != 0 {
		t.Errorf("API/TestAddComment: Expected empty response but got %v", string(response[:]))
	}
}
