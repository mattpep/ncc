package db

import "ncc/types"

import "testing"
import "database/sql"

import "log"

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
	_, err = CreateBlog("testsite")

	if err != nil {
		tb.Errorf("got error when connecting to database: %s", err)
	}

	return dbh, func(tb testing.TB) {
		log.Println("DB: teardown test suite")
		_, err := dbh.Query("TRUNCATE blogs CASCADE")
		if err != nil {
			tb.Errorf("Could not rollback database during testing")
		}
	}
}

func setupTest(dbh *sql.DB, tb testing.TB) func(fb testing.TB) {
	// log.Println("set up a test")
	// _, err := dbh.Exec("INSERT INTO blogs (ref) VALUES ('testsite')")
	// if err != nil {
	// 	tb.Errorf("got error when seeding data: %s", err)
	// }

	sql := "INSERT INTO comments (body, display_name, post_ref, blog) VALUES($1, $2, $3, $4) returning id"

	var test_comments = []types.Comment{
		{Body: "comment_body 1", DisplayName: "dispname", PostRef: "post_ref", BlogRef: "testsite"},
		{Body: "comment_body 2", DisplayName: "dispname", PostRef: "post_ref", BlogRef: "testsite"},
		{Body: "comment_body 3", DisplayName: "dispname", PostRef: "post_ref", BlogRef: "testsite"},
		{Body: "comment_body 4", DisplayName: "dispname", PostRef: "otherpost", BlogRef: "testsite"},
	}

	for _, c := range test_comments {
		_, err := dbh.Exec(sql, c.Body, c.DisplayName, c.PostRef, c.BlogRef)
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
		_, err = dbh.Query("TRUNCATE blogs CASCADE")
		if err != nil {
			tb.Errorf("Could not truncate test blogs: %v", err)
		}
	}
}

func TestAddComment(t *testing.T) {
	_, teardown := setupSuite(t)
	defer teardown(t)
	// teardown_test := setupTest(dbh, t)
	// defer teardown_test(t)

	comment := types.Comment{DisplayName: "author", Body: "comment body here", PostRef: "test_post_ref", BlogRef: "testsite"}

	t.Run("Adding a comment", func(t *testing.T) {
		got, _ := AddComment(comment)

		if got == 0 {
			t.Errorf("got 0, wanted non-zero")
		}
	})
}

func TestFlagComment(t *testing.T) {
	_, teardown := setupSuite(t)
	defer teardown(t)

	comment := types.Comment{DisplayName: "author", Body: "comment body here", PostRef: "test_post_ref", BlogRef: "testsite"}
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
		got, _ := GetPostComments("testsite", "post_ref")

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
		got, _ := GetPostComments("testsite", "post_ref")

		if got[0].Id > got[1].Id || got[1].Id > got[2].Id || got[0].Id > got[2].Id {
			t.Errorf("Posts returned in unexpected order: %v", got)
		}
	})
	t.Run("testing another post's comments", func(t *testing.T) {
		got, _ := GetPostComments("testsite", "otherpost")

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
