package types

type Comment struct {
	Id          int    `json:"id"`
	DisplayName string `json:"display_name"`
	Body        string `json:"body"`
	PostRef     string `json:"post_ref,omitempty"`
	DateTime    string `json:"date_time"`
}
