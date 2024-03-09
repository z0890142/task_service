package models

import "google.golang.org/genproto/googleapis/rpc/code"

type HttpError struct {
	Code    code.Code `json:"code" validate:"required"`
	Message string    `json:"message" validate:"required"`
}

type ErrorDetails struct {
	Row      interface{} `json:"row" validate:"required"`
	ErrorMsg string      `json:"error_msg" validate:"required"`
}
