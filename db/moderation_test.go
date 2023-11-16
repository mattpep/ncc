package db

import "testing"
import "database/sql"
import "fmt"

func setupModerationTest(dbh *sql.DB, tb testing.TB) func(fb testing.TB) {
	sql := "INSERT INTO comments  (body, display_name, post_ref) VALUES($1, $2, $3) returning id"
	var id = 0
	err := dbh.QueryRow(sql, "comment_body 1", "dispname", "post_ref").Scan(&id)
	if err != nil {
		tb.Errorf("got error when seeding data: %s", err)
	}
	fmt.Printf("inserted a comment with ID %d\n", id)
	var test_data = []string{"comment_body_2", "comment_body_3", "comment_body_4"}
	for _, t := range test_data {
		_, err = dbh.Exec(sql, t, "dispname", "post_ref")
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
		_, err := dbh.Query("TRUNCATE comments CASCADE")
		if err != nil {
			tb.Errorf("Could not truncate test comments: %v", err)
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

	pre_del_comments, _ := GetPostComments("post_ref")
	err := DeleteComment(pre_del_comments[2].Id)

	if err != nil {
		t.Errorf("Could not delete comment: %v", err)
	}
	after_del_comments, _ := GetPostComments("post_ref")

	// 4 comments of which 1 is flagged and 1 is deleted = 2 comments visible
	if len(after_del_comments) != 2 {
		t.Errorf("got %v comments, wanted 2", after_del_comments)
	}
}
