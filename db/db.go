package db

import "ncc/types"

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
)

// DB set up
func SetupDB() (*sql.DB, error) {
	db_url, present := os.LookupEnv("DATABASE_URL")
	if !present {
		return nil, errors.New("DATABASE_URL not set in the environment")
	}
	db, err := sql.Open("postgres", db_url)

	if err != nil {
		return nil, err
	}

	// log.Println("db opened")
	return db, nil
}

func GetPostComments(post_ref string) ([]types.Comment, error) {
	// TODO: check also for needs_moderation and approved flags
	db, err := SetupDB()
	if err != nil {
		return nil, errors.New("Could not connect to database")
	}
	rows, err := db.Query("SELECT id, body, display_name FROM comments WHERE post_ref = $1", post_ref)

	if err != nil {
		log.Println(fmt.Sprintf("could not read from db: %v", err))
		return nil, errors.New("Error reading from database")
	}

	var comments []types.Comment

	for rows.Next() {
		var id int
		var body string
		var display_name string
		// var post_ref string

		err = rows.Scan(&id, &body, &display_name)

		if err != nil {
			log.Println(fmt.Sprintf("could not scan row: %v", err))
			return nil, errors.New(fmt.Sprintf("Error parsing response from db: %v", err))
		}
		comments = append(comments, types.Comment{Id: id, Body: body, DisplayName: display_name, PostRef: post_ref})
	}
	return comments, nil
}

func GetPostCommentCount(post_ref string) (int, error) {
	// TODO: check also for needs_moderation and approved flags
	db, err := SetupDB()
	if err != nil {
		return 0, errors.New("Could not connect to database")
	}
	result := db.QueryRow("SELECT count(*) as c FROM comments WHERE post_ref = $1", post_ref)

	if err != nil {
		log.Println(fmt.Sprintf("could not read from db: %v", err))
		return 0, errors.New("Error reading from database")
	}

	var count int
	err = result.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func AddComment(comment types.Comment) (int, error) {
	db, err := SetupDB()
	if err != nil {
		return 0, errors.New(fmt.Sprintf("Could not set up db: %v", err))
	}
	var lastInsertID int
	err = db.QueryRow("INSERT INTO comments  (body, display_name, post_ref) VALUES($1, $2, $3) returning id;", comment.Body, comment.DisplayName, comment.PostRef).Scan(&lastInsertID)
	if err != nil || lastInsertID == 0 {
		return 0, errors.New(fmt.Sprintf("Could not write to db: %v", err))
	}
	return lastInsertID, nil
}
