package api

type Post struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Title    string `json:"title,omitempty"`
	Body     string `json:"body,omitempty"`
}
