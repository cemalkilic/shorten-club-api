package service

import "github.com/cemalkilic/shorten-backend/models"

type GetContentParams struct{
    Slug string `json:"slug" validate:"required,uri"`
}

type GetResponse struct {
    Username   string `json:"username"`
    Slug       string `json:"slug"`
    Content    interface{} `json:"content"`
    Err        error  `json:"err,omitempty"`
}

type AddRecordParams struct {
    Username   string `json:"username" validate:"omitempty,alphanum"`
    Slug       string `json:"slug" validate:"required,numeric"`
    Content    interface{} `json:"content"` // validate:"required"`
}

type AddRecordResponse struct {
    Record   models.Record `json:"record"`
    Err      error  `json:"err,omitempty"`
}

type UpdateRecordParams struct {
    Username   string `json:"username" validate:"omitempty,alphanum"`
    Slug       string `json:"slug" validate:"required,numeric"`
    Content    interface{} `json:"content"` // validate:"required"`
}

type UpdateRecordResponse struct {
    Record   models.Record `json:"record"`
    Err      error  `json:"err,omitempty"`
}