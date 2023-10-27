package types

type JsonResponse struct {
	Comments []Comment `json:comments`
	Count    int       `json:count`
}
