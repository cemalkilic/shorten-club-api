package database

import (
    "github.com/cemalkilic/shorten-backend/models"
)

type DataStore interface {
    Insert(endpoint models.Record) error
    Select(username string, slug string) (models.Record, error)
    SelectByID(id int) (models.Record, error)
    SelectBySlug(slug string) (models.Record, error)
    UpdateBySlug(record models.Record) error
    SelectAllByUser(username string) ([]models.Record, error)
    Delete(id int) error
}

type UserStore interface {
    InsertUser(user models.User) error
    SelectByUsername(username string) (models.User, error)
}
