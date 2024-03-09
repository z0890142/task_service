package models

import (
	"google.golang.org/genproto/googleapis/rpc/code"
)

type Response struct {
	Code    code.Code
	Message string
	Data    []Task
}

type ListTaskResp struct {
	*Response
	Data []Task
}
