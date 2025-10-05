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

func CreateBlog(blogref string) (int, error) {
	db, err := SetupDB()
	defer db.Close()
	var id int
	if err != nil {
		return 0, errors.New(fmt.Sprintf("Could not connect to database when trying to create blog %s: %s", blogref, err))
	}
	err = db.QueryRow("INSERT INTO blogs (ref) VALUES($1) returning id", blogref).Scan(&id)
	if err != nil {
		return 0, errors.New(fmt.Sprintf("DB: Could not create blog %s: %v", blogref, err))
	}
	return id, nil
}

func GetBlogs() ([]string, error) {
	db, err := SetupDB()
	defer db.Close()
	if err != nil {
		return []string{}, errors.New("Could not connect to database when obtaining blog list")
	}
	rows, err := db.Query(`SELECT ref FROM blogs`)
	defer rows.Close()
	if err != nil {
		return []string{}, err
	}

	var blogs = []string{}

	for rows.Next() {
		var blogname string

		err = rows.Scan(&blogname)
		if err != nil {
			return []string{}, err
		}

		blogs = append(blogs, blogname)
	}
	return blogs, nil
}

func FlagComment(comment_id int, ipaddr string) error {
	db, err := SetupDB()
	defer db.Close()
	var action int
	if err != nil {
		return errors.New(fmt.Sprintf("Could not connect to database when trying to flag comment %d", comment_id))
	}
	err = db.QueryRow("INSERT INTO moderation_actions (comment_id, action, date_time, actor) VALUES($1, 'flag', NOW(), $2) returning id", comment_id, ipaddr).Scan(&action)
	if err != nil {
		return errors.New(fmt.Sprintf("DB: Could not flag comment: %v", err))
	}
	return nil
}

func GetPostComments(blog_ref string, post_ref string) ([]types.Comment, error) {
	// TODO: check also for needs_moderation and approved flags
	db, err := SetupDB()
	defer db.Close()
	if err != nil {
		return nil, errors.New("Could not connect to database when getting a post's comments")
	}
	rows, err := db.Query(`SELECT comments.id, comments.display_name, comments.body, comments.created_at
                               FROM comments
                               LEFT OUTER JOIN (
                                        SELECT id, date_time, comment_id, action, ROW_NUMBER() OVER (PARTITION BY comment_id ORDER BY DATE_TIME DESC) recency
                                        FROM moderation_actions) rma
                               ON rma.comment_id = comments.id
                               WHERE ((recency = 1 AND action='approve') OR recency IS NULL)
				AND post_ref = $1 AND blog = $2
			       ORDER BY created_at ASC`, post_ref, blog_ref)

	if err != nil {
		log.Println(fmt.Sprintf("could not read from db: %v", err))
		return nil, errors.New("Error reading from database when getting a post's comments")
	}

	var comments []types.Comment

	for rows.Next() {
		var id int
		var body string
		var display_name string
		var date_time string
		// var post_ref string

		err = rows.Scan(&id, &display_name, &body, &date_time)

		if err != nil {
			log.Println(fmt.Sprintf("could not scan row: %v", err))
			return nil, errors.New(fmt.Sprintf("Error parsing response from db: %v", err))
		}
		comments = append(comments, types.Comment{Id: id, Body: body, DisplayName: display_name, PostRef: post_ref, DateTime: date_time})
	}
	return comments, nil
}

func GetBlogCommentCounts(blog_ref string) ([]types.PostCommentCount, error) {
	db, err := SetupDB()
	defer db.Close()
	if err != nil {
		return []types.PostCommentCount{}, errors.New("Could not connect to database")
	}
	rows, err := db.Query(` SELECT post_ref, COUNT(*) AS c
				FROM comments
				LEFT OUTER JOIN (
					SELECT id, date_time, comment_id, action,
					       ROW_NUMBER() OVER (PARTITION BY comment_id
					ORDER BY DATE_TIME DESC) recency
					FROM moderation_actions) rma
				ON rma.comment_id = comments.id
				WHERE ((recency = 1 AND action='approve') OR recency IS NULL)
				AND blog = $1 GROUP BY post_ref`, blog_ref)
	defer rows.Close()
	if err != nil {
		return []types.PostCommentCount{}, err
	}

	var counts = []types.PostCommentCount{}

	for rows.Next() {
		var post_ref string
		var count int

		err = rows.Scan(&post_ref, &count)
		if err != nil {
			return []types.PostCommentCount{}, err
		}

		counts = append(counts, types.PostCommentCount{PostRef: post_ref, CommentCount: count})
	}
	return counts, nil
}

func GetPostCommentCount(post_ref string) (int, error) {
	db, err := SetupDB()
	defer db.Close()
	if err != nil {
		return 0, errors.New("Could not connect to database")
	}
	result := db.QueryRow(`SELECT COUNT(*) AS c
                               FROM comments
                               LEFT OUTER JOIN (
                                        SELECT id, date_time, comment_id, action, ROW_NUMBER() OVER (PARTITION BY comment_id ORDER BY DATE_TIME DESC) recency
                                        FROM moderation_actions) rma
                               ON rma.comment_id = comments.id
                               WHERE ((recency = 1 AND action='approve') OR recency IS NULL)
				AND post_ref = $1`, post_ref)

	var count int
	err = result.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func AddComment(comment types.Comment) (int, error) {
	db, err := SetupDB()
	defer db.Close()
	if err != nil {
		return 0, errors.New(fmt.Sprintf("DB/AddComment: Could not set up db: %v", err))
	}
	var lastInsertID int
	log.Println(fmt.Sprintf("DB/AddComment: Inserting a comment: %v", comment))

	err = db.QueryRow("INSERT INTO comments (body, display_name, post_ref, blog) VALUES($1, $2, $3, $4) returning id;", comment.Body, comment.DisplayName, comment.PostRef, comment.BlogRef).Scan(&lastInsertID)
	if err != nil || lastInsertID == 0 {
		return 0, errors.New(fmt.Sprintf("Could not write to db: %v", err))
	}
	log.Println(fmt.Sprintf("DB/AddComment: Comment was given id %d", lastInsertID))
	return lastInsertID, nil
}
