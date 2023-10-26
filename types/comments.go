package types

type Comment struct {
	Id          int    `json:"id"`
	DisplayName string `json:"display_name"`
	Body        string `json:"body"`
}
type CommentEntry struct {
	Id          int    `json:"id"`
	DisplayName string `json:"display_name"`
	Body        string `json:"body"`
	PostRef     string `json:"post_ref"`
}
