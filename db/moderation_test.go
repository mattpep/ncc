package db

import (
	"database/sql"
	"fmt"
	"log"
	"testing"
)

func setupModerationTest(dbh *sql.DB, tb testing.TB) func(fb testing.TB) {
	var id = 0

	sql := "INSERT INTO comments (body, display_name, post_ref, blog) VALUES($1, $2, $3, $4) returning id"
	err := dbh.QueryRow(sql, "comment_body 1", "dispname", "post_ref", "testsite").Scan(&id)
	if err != nil {
		tb.Errorf("got error when seeding data: %s", err)
	}
	log.Println(fmt.Sprintf("Test comment for moderation has ID %d", id))
	var test_data = []string{"comment_body_2", "comment_body_3", "comment_body_4"}
	for _, t := range test_data {
		_, err = dbh.Exec(sql, t, "dispname", "post_ref", "testsite")
		if err != nil {
			tb.Errorf("got error when seeding data: %s", err)
		}
	}

	sql = "INSERT INTO moderation_actions (comment_id, action, actor) VALUES ($1, 'flag', '::1') RETURNING id"
	_, err = dbh.Exec(sql, id)
	if err != nil {
		tb.Errorf("Could not seed moderation actions: %v", err)
	}

	return func(tb testing.TB) {
		_, err := dbh.Query("TRUNCATE blogs CASCADE")
		if err != nil {
			tb.Errorf("Could not truncate test data: %v", err)
		}
	}
}
func TestGetModerationTasks(t *testing.T) {
	dbh, teardown := setupSuite(t)
	defer teardown(t)
	teardown_test := setupModerationTest(dbh, t)
	defer teardown_test(t)

	t.Run("With some tasks", func(t *testing.T) {
		got, _ := GetModerationTasks()

		if len(got) != 1 {
			t.Errorf("got %v tasks, wanted 1 task", got)
		}
	})
}

func TestApproveComment(t *testing.T) {
	dbh, teardown := setupSuite(t)
	defer teardown(t)
	teardown_test := setupModerationTest(dbh, t)
	defer teardown_test(t)

	for_approval, _ := GetModerationTasks()
	err := ApproveComment(for_approval[0].Id)

	if err != nil {
		t.Errorf("Could not approve comment: %v", err)
	}
	got, _ := GetModerationTasks()

	if len(got) != 0 {
		t.Errorf("got %v tasks, wanted 0 task", got)
	}
}
func TestDeleteComment(t *testing.T) {
	dbh, teardown := setupSuite(t)
	defer teardown(t)
	teardown_test := setupModerationTest(dbh, t)
	defer teardown_test(t)

	pre_del_comments, _ := GetPostComments("testsite", "post_ref")
	err := DeleteComment(pre_del_comments[2].Id)

	if err != nil {
		t.Errorf("Could not delete comment: %v", err)
	}
	after_del_comments, _ := GetPostComments("testsite", "post_ref")

	// 4 comments of which 1 is flagged and 1 is deleted = 2 comments visible
	if len(after_del_comments) != 2 {
		t.Errorf("got %v comments, wanted 2", after_del_comments)
	}
}
