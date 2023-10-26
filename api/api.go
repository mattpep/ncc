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

func AddComment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	// check the parent post exists (404? 400?)
	// check post is not locked (and/or that comments are allowed on this post) - 403 if locked
	// create comment
	log.Printf("new comment: %v", params)
	comment := types.CommentEntry{DisplayName: params["name"], Body: params["body"], PostRef: params["post_ref"]}
	lastInsertID, err := db.AddComment(comment)
	if err != nil || lastInsertID == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "{\"status\":\"error\",\"message\":\"Database error when storing comment\"}")

		fmt.Println(fmt.Sprintf("Database error when writing: %v", err))
	} else {
		fmt.Println(fmt.Sprintf("Created record %d", lastInsertID))
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
	json.NewEncoder(w).Encode(response)
}

func GetPostComments(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	post_ref := params["postref"]
	dbcomments, err := db.GetPostComments(post_ref)
	fmt.Printf("response is %v", dbcomments)
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
	var response = types.JsonResponse{Comments: json_comments}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(response)
}
