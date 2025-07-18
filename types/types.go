package types

type Comment struct {
	Id          int    `json:"id"`
	DisplayName string `json:"display_name"`
	Body        string `json:"body"`
	PostRef     string `json:"post_ref,omitempty"`
	BlogRef     string `json:"blog"`
	DateTime    string `json:"date_time"`
}

type CommentCountInfo struct {
	PostRef string `json:"post_ref"`
	BlogRef string `json:"blog"`
	Count   int    `json:"count"`
}

type PostCommentCount struct {
	PostRef      string `json:"postref"`
	CommentCount int    `json:"count"`
}

type BlogCommentCounts struct {
	Status    string             `json:"status"`
	CountInfo []PostCommentCount `json:"counts"`
}

type JsonResponse struct {
	Status   string    `json:"status"`
	Comments []Comment `json:"comments"`
	Count    int       `json:"count"`
}

type ModTask struct {
	Id          int
	DisplayName string
	Body        string
	PostRef     string
	DateTime    string
	Actor       string
}
