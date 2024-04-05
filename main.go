package main

import "ncc/api"
import "ncc/moderator"

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gorilla/mux"
	"net/http"
	"runtime/debug"
)

func runServer(port string) {
	router := mux.NewRouter()

	router.HandleFunc("/comments/{postref}", api.GetPostComments).Methods("GET")
	router.HandleFunc("/commentcount/{postref}", api.GetPostCommentCount).Methods("GET")
	router.HandleFunc("/commentcount/{postref}", api.OptionsRequest).Methods("OPTIONS")
	router.HandleFunc("/comments/{postref}", api.AddComment).Methods("POST")
	router.HandleFunc("/flag/{comment_id}", api.FlagComment).Methods("POST")
	router.HandleFunc("/comments/{postref}", api.OptionsRequest).Methods("OPTIONS")
	router.HandleFunc("/js/insert/form", FormInsert).Methods("GET")
	router.HandleFunc("/js/insert/counts", CountInsert).Methods("GET")

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
			if info, ok := debug.ReadBuildInfo(); ok {
				for _, setting := range info.Settings {
					if setting.Key == "vcs.revision" {
						fmt.Printf("Version: %s\n", setting.Value)
					}
				}
			} else {
				fmt.Println("Unknown version. Not built from a git repo?")
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
