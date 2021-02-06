package service

import "github.com/cemalkilic/shorten-backend/models"

type GetContentParams struct{
    Slug string `json:"slug" validate:"required,uri"`
}

type GetResponse struct {
    Record   models.Record `json:"record"`
    Err        error  `json:"err,omitempty"`
}

type AddRecordParams struct {
    Username   string `json:"username" validate:"omitempty,alphanum"`
    Slug       string `json:"slug" validate:"required,numeric"`
    Type       string `json:"type" validate:"required"`
    Content    interface{} `json:"content"` // validate:"required"`
}

type AddRecordResponse struct {
    Record   models.Record `json:"record"`
    Err      error  `json:"err,omitempty"`
}

type UpdateRecordParams struct {
    Slug       string `json:"slug" validate:"required,numeric"`
    Type       string `json:"type"`
    Content    interface{} `json:"content"` // validate:"required"`
}

type UpdateRecordResponse struct {
    Record   models.Record `json:"record"`
    Err      error  `json:"err,omitempty"`
}