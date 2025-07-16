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
	"strconv"
	"strings"
)

func OptionsRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.WriteHeader(http.StatusNoContent)
	io.WriteString(w, "")
}

// Get the IP address of the server's connected user.
func getUserIP(req *http.Request) string {
	var ipaddr string
	if len(req.Header.Get("CF-Connecting-IP")) > 1 {
		ipaddr = req.Header.Get("CF-Connecting-IP")
	} else if len(req.Header.Get("X-Forwarded-For")) > 1 {
		ipaddr = req.Header.Get("X-Forwarded-For")
	} else if len(req.Header.Get("X-Real-IP")) > 1 {
		ipaddr = req.Header.Get("X-Real-IP")
	} else {
		parts := strings.Split(req.RemoteAddr, ":")
		parts = parts[:len(parts)-1]
		ipaddr = strings.Join(parts[:], ":")
		if ipaddr[0] == '[' {
			ipaddr = ipaddr[1 : len(ipaddr)-1]
		}
	}
	return ipaddr
}

func FlagComment(w http.ResponseWriter, r *http.Request) {
	// params := mux.Vars(r)
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"status":"error","message":"Error parsing POST request"}`)
		log.Println(fmt.Sprintf("Could not parse POST request: %v", err))
		return
	}

	params := mux.Vars(r)
	// comment_id, err := strconv.Atoi(r.PostForm.Get("comment_id"))
	comment_id, err := strconv.Atoi(params["comment_id"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"status":"error","message":"Error getting commentid when flagging"}`)
		log.Println(fmt.Sprintf("Error getting commentid when flagging: %v", err))
		return
	}
	// check the parent post exists (404? 400?)
	// check post is not locked (and/or that comments are allowed on this post) - 403 if locked

	err = db.FlagComment(comment_id, getUserIP(r))

	if err != nil {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"status":"error","message":"Database error when flagging comment"}`)
		log.Println(fmt.Sprintf("Database error when flagging comment %d: %v", comment_id, err))
	} else {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		// we need to set the status header _after_ the other headers
		// even though it gets emitted first in the output. If we don't
		// then these extra CORS headers won't get emitted.
		w.WriteHeader(http.StatusNoContent)
		io.WriteString(w, "")
	}
}

// Takes a JSON submission, not x-urlencoded
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
		log.Println(fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	comment := types.Comment{DisplayName: (form_fields["display_name"]).(string), Body: (form_fields["body"]).(string), PostRef: params["postref"]}

	lastInsertID, err := db.AddComment(comment)
	if err != nil || lastInsertID == 0 {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"status":"error","message":"Database error when storing comment"}`)

		log.Println(fmt.Sprintf("Database error when writing: %v", err))
	} else {
		log.Println(fmt.Sprintf("Created record %d", lastInsertID))
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
	// fmt.Printf("GetPostCommentCount: returning a count of %d\n", response.Count)
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
		log.Println(fmt.Sprintf("Error reading comments from database: %v", err))
		io.WriteString(w, `{"status":"error","message":"Error reading comments from database"}`)
		return
	}
	json_comments := []types.Comment{}

	for _, dbcomment := range dbcomments {
		json_comments = append(json_comments, dbcomment)
	}
	var response = types.JsonResponse{Status: "ok", Comments: json_comments, Count: len(json_comments)}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	json.NewEncoder(w).Encode(response)
}
