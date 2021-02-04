package models

type Record struct {
    ID       int64 `json:"id"`
    Username string `json:"username"`
    Slug     string `json:"slug"`
    Content  interface{} `json:"content"`
}
