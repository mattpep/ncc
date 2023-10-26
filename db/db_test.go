package db

import "ncc/types"
import "testing"
import "database/sql"
import "log"

func setupSuite(tb testing.TB) (*sql.DB, func(tb testing.TB)) {
	log.Println("setup test suite")
	dbh, err := SetupDB()
	if err != nil {
		tb.Errorf("got error when connecting to database: %s", err)
	}
	// _, err = dbh.Exec("BEGIN")
	// if err != nil {
	// 	tb.Errorf("Could not open transaction during testing")
	// }

	if err != nil {
		tb.Errorf("got error when connecting to database: %s", err)
	}

	return dbh, func(tb testing.TB) {
		log.Println("teardown test suite")
		// _, err := dbh.Query("ROLLBACK")
		// if err != nil {
		// 	tb.Errorf("Could not rollback database during testing")
		// }

	}
}

func setupTest(dbh *sql.DB, tb testing.TB) func(fb testing.TB) {
	log.Println("set up a test")

	sql := "INSERT INTO comments  (body, display_name, post_ref) VALUES($1, $2, $3) returning id"
	_, err := dbh.Exec(sql, "comment_body", "dispname", "post_ref")
	_, err = dbh.Exec(sql, "comment_body", "dispname", "post_ref")
	_, err = dbh.Exec(sql, "comment_body", "dispname", "post_ref")
	_, err = dbh.Exec(sql, "comment_body", "dispname", "otherpost")
	if err != nil {
		tb.Errorf("got error when seeding data: %s", err)
	}

	return func(tb testing.TB) {
		log.Println("tear down a test")
		_, err := dbh.Query("TRUNCATE comments")
		if err != nil {
			tb.Errorf("Could not truncate test comments")
		}
	}
}

func TestAddComment(t *testing.T) {
	_, teardown := setupSuite(t)
	defer teardown(t)
	// teardown_test := setupTest(dbh, t)
	// defer teardown_test(t)

	comment := types.CommentEntry{DisplayName: "author", Body: "comment body here", PostRef: "test_post_ref"}

	t.Run("Adding a comment", func(t *testing.T) {
		got, _ := AddComment(comment)

		if got == 0 {
			t.Errorf("got %v, wanted non-zero", got)
		}
	})
}

func TestGetPostComments(t *testing.T) {
	dbh, teardown := setupSuite(t)
	defer teardown(t)
	teardown_test := setupTest(dbh, t)
	defer teardown_test(t)

	t.Run("testing one post's comments", func(t *testing.T) {
		got, _ := GetPostComments("post_ref")

		if len(got) != 3 {
			t.Errorf("got %v, wanted 3", len(got))
		}
	})
	t.Run("testing another post's comments", func(t *testing.T) {
		got, _ := GetPostComments("otherpost")

		if len(got) != 1 {
			t.Errorf("got %v, wanted 1", len(got))
		}
	})
}
func TestGetPostCommentCount(t *testing.T) {
	dbh, teardown := setupSuite(t)
	defer teardown(t)
	teardown_test := setupTest(dbh, t)
	defer teardown_test(t)

	t.Run("testing one post's comments", func(t *testing.T) {
		got, _ := GetPostCommentCount("post_ref")

		if got != 3 {
			t.Errorf("got %v, wanted 3", got)
		}
	})
	t.Run("testing another post's comments", func(t *testing.T) {
		got, _ := GetPostCommentCount("otherpost")

		if got != 1 {
			t.Errorf("got %v, wanted 1", got)
		}
	})
}
