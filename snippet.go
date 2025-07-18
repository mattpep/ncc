package main

import (
	"bytes"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"os"
	"text/template"

	_ "embed"
)

type templateParams struct {
	Endpoint string
	PostRef  string
	BlogRef  string
}

var endpoint string

//go:embed snippets/comment.js
var comment_snippet string

//go:embed snippets/count.js
var count_snippet string

func countInsert(blogRef string) (string, error) {
	buf := &bytes.Buffer{}
	tmpl, err := template.New("tmpl").Parse(count_snippet)
	if err != nil {
		return "", err
	}

	err = tmpl.Execute(buf, templateParams{Endpoint: endpoint, BlogRef: blogRef})
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func formInsert(blogRef string) (string, error) {
	buf := &bytes.Buffer{}
	log.Println(fmt.Sprintf("Snippet: Rendering form insert with blogref: %s", blogRef))
	tmpl, err := template.New("tmpl").Parse(comment_snippet)
	if err != nil {
		return "", err
	}

	err = tmpl.Execute(buf, templateParams{Endpoint: endpoint, BlogRef: blogRef})
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func CountInsert(w http.ResponseWriter, r *http.Request) {
	ext_endpoint, present := os.LookupEnv("EXT_ENDPOINT")
	params := mux.Vars(r)
	blog_ref := params["blogref"]
	log.Println(fmt.Sprintf("Snippet: Making the count insert for %s", blog_ref))
	if present {
		endpoint = ext_endpoint
	} else {
		endpoint = "http://localhost:8080"
	}

	count_snippet, err := countInsert(blog_ref)
	if err != nil {
		log.Println(fmt.Sprintf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"status":"error","message":"Cannot serve snippet"}`)
		return
	}
	ServeWebsiteInsert(w, r, count_snippet)
}

func FormInsert(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	blogref := params["blogref"]

	ext_endpoint, present := os.LookupEnv("EXT_ENDPOINT")
	if present {
		endpoint = ext_endpoint
	} else {
		endpoint = "http://localhost:8080"
	}

	comment_snippet, err := formInsert(blogref)
	if err != nil {
		log.Println(fmt.Sprintf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"status":"error","message":"Cannot serve snippet"}`)
		return
	}
	log.Println(fmt.Sprintf("serving js insert: blogref is >%s<", blogref))
	ServeWebsiteInsert(w, r, comment_snippet)
}

func ServeWebsiteInsert(w http.ResponseWriter, r *http.Request, content string) {
	w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	io.WriteString(w, content)
}
