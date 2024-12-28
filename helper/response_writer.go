package helper

import (
	"go-mnemosyne-api/web"
)

func WriteSuccess(message string, data interface{}) web.ResponseContract {
	return web.ResponseContract{
		Status:  true,
		Message: message,
		Data:    &data,
	}
}
