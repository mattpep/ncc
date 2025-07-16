package main

import "ncc/api"
import "ncc/moderator"

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/gorilla/mux"
	"net/http"
	"runtime/debug"
)

func requestLogger(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Request: " + r.Method + " " + r.URL.Path)
		f(w, r)
	}
}

func setupRouter() *mux.Router {
	router := mux.NewRouter()
	// legacy endpoints
	router.HandleFunc("/flag/{comment_id}", requestLogger(api.FlagComment)).Methods("POST")
	router.HandleFunc("/comments/{postref}", requestLogger(api.OptionsRequest)).Methods("OPTIONS")

	// v2 endpoints, which have a blogref
	router.HandleFunc("/v2/blog/{blogref}/commentcounts", requestLogger(api.GetBlogCommentCounts)).Methods("GET")
	router.HandleFunc("/v2/blog/{blogref}/commentcounts", requestLogger(api.OptionsRequest)).Methods("OPTIONS")
	router.HandleFunc("/v2/blog/{blogref}/comments/{postref}", requestLogger(api.GetPostComments)).Methods("GET")
	router.HandleFunc("/v2/blog/{blogref}/comments/{postref}", requestLogger(api.AddComment)).Methods("POST")
	router.HandleFunc("/v2/blog/{blogref}/comments/{postref}", requestLogger(api.OptionsRequest)).Methods("OPTIONS")
	router.HandleFunc("/v2/js/blog/{blogref}/counts", requestLogger(CountInsert)).Methods("GET")

	// static endpoints for Javascript inserts
	router.HandleFunc("/js/insert/blog/{blogref}/form", requestLogger(FormInsert)).Methods("GET")

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "endpoint not found",
		})
	})

	return router
}

func runServer(port string) {
	// serve the app
	fmt.Println("ncc - no cookies comment system")
	fmt.Println("Copyright 2023 by Matt Peperell")
	version, err := getVCSVersion()
	if err != nil {
		fmt.Println("Could not determine version: Not built from a git repo?")
	} else {
		fmt.Printf("Git version: %s\n", version)
	}
	router := setupRouter()
	fmt.Printf("Server running at %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func getVCSVersion() (string, error) {
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				return setting.Value, nil
			}
		}
	}
	return "", errors.New("Unknown version. Not built from a git repo?")

}

func main() {
	if len(os.Args) > 1 {
		if os.Args[1] == "modqueue" {
			moderator.ShowTasks()
		} else if os.Args[1] == "server" {
			port, present := os.LookupEnv("PORT")
			if !present {
				port = "8080"
			}
			runServer(port)
		} else if os.Args[1] == "version" {
			version, err := getVCSVersion()
			if err != nil {
				fmt.Println("Not built from a git repo?")
			} else {
				fmt.Printf("Git version: %s\n", version)
			}
		} else {
			fmt.Println("Unknown action")
		}
	} else {
		port, present := os.LookupEnv("PORT")
		if !present {
			port = "8080"
		}
		runServer(port)
	}
}
