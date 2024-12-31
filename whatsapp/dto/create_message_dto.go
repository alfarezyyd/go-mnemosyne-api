package dto

import "go-mnemosyne-api/model/whatsapp"

type PayloadMessageDto struct {
	Object string           `json:"object"`
	Entry  []whatsapp.Entry `json:"entry"`
}
