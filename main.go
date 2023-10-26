package main

import "ncc/api"
import "ncc/snippet"

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gorilla/mux"
	"net/http"
)

func runServer(port string) {
	router := mux.NewRouter()

	router.HandleFunc("/comments/{postref}", api.GetPostComments).Methods("GET")
	router.HandleFunc("/commentcount/{postref}", api.GetPostCommentCount).Methods("GET")
	router.HandleFunc("/comments/{postref}", api.AddComment).Methods("POST")
	router.HandleFunc("/js/insert", snippet.ServeWebsiteInsert).Methods("GET")

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "endpoint not found",
		})
	})

	// serve the app
	fmt.Println("ncc - no cookies comment system")
	fmt.Println("Copyright 2023 by Matt Peperell")
	fmt.Printf("Server running at %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func main() {
	port, present := os.LookupEnv("PORT")
	if !present {
		port = "8080"
	}
	runServer(port)
}
