package db

import (
	"errors"
	"fmt"
	"log"
)
import "ncc/types"

func GetModerationTasks() ([]types.ModTask, error) {
	// gets a list of comments for which the most recent action is a flagging
	db, err := SetupDB()
	defer db.Close()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not set up db: %v", err))
	}
	flagged_comments_sql := `SELECT comments.id, comments.display_name, comments.body, rma.date_time, comments.post_ref, rma.actor, comments.blog
				FROM comments
				INNER JOIN (
					SELECT id, date_time, comment_id, action, actor, ROW_NUMBER() OVER (PARTITION BY comment_id ORDER BY DATE_TIME DESC) recency
					FROM moderation_actions) rma
				ON rma.comment_id = comments.id
				WHERE recency = 1 AND action = 'flag'`
	rows, err := db.Query(flagged_comments_sql)

	if err != nil {
		log.Println(fmt.Sprintf("could not read moderation tasks from db: %v", err))
		return nil, errors.New("Error reading moderation tasks from database")
	}

	var tasks []types.ModTask

	for rows.Next() {
		var id int
		var body string
		var display_name string
		var date_time string
		var post_ref string
		var blog_ref string
		var actor string

		err = rows.Scan(&id, &display_name, &body, &date_time, &post_ref, &actor, &blog_ref)

		if err != nil {
			log.Println(fmt.Sprintf("could not scan row: %v", err))
			return nil, errors.New(fmt.Sprintf("Error parsing response from db: %v", err))
		}
		tasks = append(tasks, types.ModTask{Id: id, Body: body, DisplayName: display_name, PostRef: post_ref, DateTime: date_time, Actor: actor, BlogRef: blog_ref})
	}
	return tasks, nil
}

func ApproveComment(comment_id int) error {
	var action int
	db, err := SetupDB()
	defer db.Close()
	if err != nil {
		return errors.New("Could not connect to database")
	}
	err = db.QueryRow("INSERT INTO moderation_actions (comment_id, action, date_time, actor) VALUES($1, 'approve', NOW(), '::1') returning id", comment_id).Scan(&action)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not approve comment: %v", err))
	}
	return nil
}
func DeleteComment(comment_id int) error {
	db, err := SetupDB()
	defer db.Close()
	if err != nil {
		return errors.New("Could not connect to database")
	}
	_, err = db.Exec("DELETE FROM comments WHERE id = $1", comment_id)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not delete comment: %v", err))
	}
	return nil
}
