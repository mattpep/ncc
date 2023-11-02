package snippet

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

//go:embed snippet.js
var snippet string

func websiteInsert(post_ref string) (string, error) {
	buf := &bytes.Buffer{}
	tmpl, err := template.New("tmpl").Parse(snippet)
	if err != nil {
		return "", err
	}

	err = tmpl.Execute(buf, templateParams{Endpoint: endpoint, PostRef: post_ref})
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func ServeWebsiteInsert(w http.ResponseWriter, r *http.Request) {
	postref := r.URL.Query().Get("postref")

	ext_endpoint, present := os.LookupEnv("EXT_ENDPOINT")
	if present {
		endpoint = ext_endpoint
	} else {
		endpoint = "http://localhost:8080"
	}

	// log.Printf("serving js insert: post ref is >%v<", postref)
	snippet, err := websiteInsert(postref)
	if err != nil {
		log.Printf("Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "{\"status\":\"error\",\"message\":\"Cannot serve snippet\"}")
		return
	}
	w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	io.WriteString(w, snippet)
}
