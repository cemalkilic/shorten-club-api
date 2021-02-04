package models

import "time"

type User struct {
    Username  string    `json:"username" validate:"required"`
    Password  string    `json:"password" validate:"required"`
    CreatedAt time.Time `json:"createdAt"`
}
