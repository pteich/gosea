package api

type Post struct {
	Username    string `json:"username"`
	CompanyName string `json:"company_name"`
	Name        string `json:"name"`
	Title       string `json:"title,omitempty"`
	Body        string `json:"body,omitempty"`
}
