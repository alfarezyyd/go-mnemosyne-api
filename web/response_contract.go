package web

import "go-mnemosyne-api/exception"

type ResponseContract struct {
	Status  bool                   `json:"status"`
	Message string                 `json:"message"`
	Data    interface{}            `json:"data,omitempty"`
	Error   *exception.ClientError `json:"error,omitempty"`
}
