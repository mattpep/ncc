package api

import "ncc/db"
import "ncc/types"

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
)

func OptionsRequest(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling OPTIONS request")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.WriteHeader(http.StatusNoContent)
	io.WriteString(w, "")
}

func AddComment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	// check the parent post exists (404? 400?)
	// check post is not locked (and/or that comments are allowed on this post) - 403 if locked
	var form_fields map[string]interface{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"status":"error","message":"Error parsing POST request"}`)
		return
	}

	err = json.Unmarshal([]byte(body), &form_fields)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"status":"error","message":"Error parsing JSON"}`)
		fmt.Println(fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	comment := types.CommentEntry{DisplayName: (form_fields["display_name"]).(string), Body: (form_fields["body"]).(string), PostRef: params["postref"]}

	lastInsertID, err := db.AddComment(comment)
	if err != nil || lastInsertID == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		io.WriteString(w, `{"status":"error","message":"Database error when storing comment"}`)

		fmt.Println(fmt.Sprintf("Database error when writing: %v", err))
	} else {
		fmt.Println(fmt.Sprintf("Created record %d", lastInsertID))
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		// we need to set the status header _after_ the other headers
		// even though it gets emitted first in the output. If we don't
		// then these extra CORS headers won't get emitted.
		w.WriteHeader(http.StatusNoContent)
		io.WriteString(w, "")
	}
}

func GetPostCommentCount(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	post_ref := params["postref"]
	count, _ := db.GetPostCommentCount(post_ref)

	type countresponse struct {
		Status string `json:"status"`
		Count  int    `json:"count"`
	}
	response := countresponse{
		Status: "ok",
		Count:  count,
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	json.NewEncoder(w).Encode(response)
}

func GetPostComments(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	post_ref := params["postref"]
	dbcomments, err := db.GetPostComments(post_ref)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error saving comment to database: %v", err)
		json.NewEncoder(w).Encode("{\"status\":\"error\",\"message\":\"Error saving to database\"}")
		return
	}
	// comments := []types.CommentEntry{}
	json_comments := []types.Comment{}

	for _, dbcomment := range dbcomments {
		json_comment := types.Comment{Id: dbcomment.Id, DisplayName: dbcomment.DisplayName, Body: dbcomment.Body}
		json_comments = append(json_comments, json_comment)
	}
	var response = types.JsonResponse{Comments: json_comments, Count: len(json_comments)}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	json.NewEncoder(w).Encode(response)
}
