package models

type CustomEndpoint struct  {
    ID         int    `json:"id"`
    Username   string `json:"username"`
    URI        string `json:"uri"`
    Content    string `json:"content"`
    StatusCode int    `json:"statusCode"`
}
