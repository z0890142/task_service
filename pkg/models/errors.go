package models

import "google.golang.org/genproto/googleapis/rpc/code"

type HttpError struct {
	Code    code.Code      `json:"code" validate:"required"`
	Message string         `json:"message" validate:"required"`
	Detail  []ErrorDetails `json:"error" validate:"required"`
}

type ErrorDetails struct {
	Row      interface{} `json:"row" validate:"required"`
	ErrorMsg string      `json:"error_msg" validate:"required"`
}

func NewHttpError(code code.Code, message string, Details []ErrorDetails) HttpError {
	return HttpError{
		Code:    code,
		Message: message,
		Detail:  Details,
	}
}
