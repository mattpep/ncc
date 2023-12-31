package db

import "ncc/types"

import "testing"
import "database/sql"

// import "log"

func setupSuite(tb testing.TB) (*sql.DB, func(tb testing.TB)) {
	// log.Println("setup test suite")

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
		// log.Println("teardown test suite")
		// _, err := dbh.Query("ROLLBACK")
		// if err != nil {
		// 	tb.Errorf("Could not rollback database during testing")
		// }

	}
}

func setupTest(dbh *sql.DB, tb testing.TB) func(fb testing.TB) {
	// log.Println("set up a test")
	sql := "INSERT INTO comments  (body, display_name, post_ref) VALUES($1, $2, $3) returning id"

	var test_comments = []types.Comment{
		types.Comment{Body: "comment_body 1", DisplayName: "dispname", PostRef: "post_ref"},
		types.Comment{Body: "comment_body 2", DisplayName: "dispname", PostRef: "post_ref"},
		types.Comment{Body: "comment_body 3", DisplayName: "dispname", PostRef: "post_ref"},
		types.Comment{Body: "comment_body 4", DisplayName: "dispname", PostRef: "otherpost"},
	}

	for _, c := range test_comments {
		_, err := dbh.Exec(sql, c.Body, c.DisplayName, c.PostRef)
		if err != nil {
			tb.Errorf("got error when seeding data: %s", err)
		}
	}

	return func(tb testing.TB) {
		// log.Println("tear down a test")
		_, err := dbh.Query("TRUNCATE comments CASCADE")
		if err != nil {
			tb.Errorf("Could not truncate test comments: %v", err)
		}
	}
}

func TestAddComment(t *testing.T) {
	_, teardown := setupSuite(t)
	defer teardown(t)
	// teardown_test := setupTest(dbh, t)
	// defer teardown_test(t)

	comment := types.Comment{DisplayName: "author", Body: "comment body here", PostRef: "test_post_ref"}

	t.Run("Adding a comment", func(t *testing.T) {
		got, _ := AddComment(comment)

		if got == 0 {
			t.Errorf("got %v, wanted non-zero", got)
		}
	})
}

func TestFlagComment(t *testing.T) {
	_, teardown := setupSuite(t)
	defer teardown(t)

	comment := types.Comment{DisplayName: "author", Body: "comment body here", PostRef: "test_post_ref"}
	t.Run("Flagging a comment", func(t *testing.T) {
		id, err := AddComment(comment)
		if err != nil {
			t.Errorf("Could not add comment: %v", err)
		}
		err = FlagComment(id, "0.0.0.0")
		if err != nil {
			t.Errorf("Could not flag comment: %v", err)
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
	t.Run("posts returned in order", func(t *testing.T) {
		sql := "UPDATE comments SET body = $1 WHERE body = $2"
		_, err := dbh.Exec(sql, "updated body 1 (one)", "comment_body 1")
		if err != nil {
			t.Errorf("got error when trying to update a db record in a test: %v", err)
		}
		got, _ := GetPostComments("post_ref")

		if got[0].Id > got[1].Id || got[1].Id > got[2].Id || got[0].Id > got[2].Id {
			t.Errorf("Posts returned in unexpected order: %v", got)
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
