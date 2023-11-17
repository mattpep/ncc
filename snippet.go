package main

import (
	"bytes"
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
}

var endpoint string

//go:embed snippets/comment.js
var comment_snippet string

//go:embed snippets/count.js
var count_snippet string

func countInsert() (string, error) {
	buf := &bytes.Buffer{}
	tmpl, err := template.New("tmpl").Parse(count_snippet)
	if err != nil {
		return "", err
	}

	err = tmpl.Execute(buf, templateParams{Endpoint: endpoint})
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func formInsert(post_ref string) (string, error) {
	buf := &bytes.Buffer{}
	tmpl, err := template.New("tmpl").Parse(comment_snippet)
	if err != nil {
		return "", err
	}

	err = tmpl.Execute(buf, templateParams{Endpoint: endpoint, PostRef: post_ref})
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func CountInsert(w http.ResponseWriter, r *http.Request) {
	ext_endpoint, present := os.LookupEnv("EXT_ENDPOINT")
	if present {
		endpoint = ext_endpoint
	} else {
		endpoint = "http://localhost:8080"
	}

	count_snippet, err := countInsert()
	if err != nil {
		log.Printf("Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"status":"error","message":"Cannot serve snippet"}`)
		return
	}
	ServeWebsiteInsert(w, r, count_snippet)
}

func FormInsert(w http.ResponseWriter, r *http.Request) {
	postref := r.URL.Query().Get("postref")

	ext_endpoint, present := os.LookupEnv("EXT_ENDPOINT")
	if present {
		endpoint = ext_endpoint
	} else {
		endpoint = "http://localhost:8080"
	}

	// log.Printf("serving js insert: post ref is >%v<", postref)
	comment_snippet, err := formInsert(postref)
	if err != nil {
		log.Printf("Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"status":"error","message":"Cannot serve snippet"}`)
		return
	}
	ServeWebsiteInsert(w, r, comment_snippet)
}

func ServeWebsiteInsert(w http.ResponseWriter, r *http.Request, content string) {
	w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	io.WriteString(w, content)
}
