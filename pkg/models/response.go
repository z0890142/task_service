package models

import (
	"google.golang.org/genproto/googleapis/rpc/code"
)

type Response struct {
	Code    code.Code
	Message string
	Data    interface{}
}
