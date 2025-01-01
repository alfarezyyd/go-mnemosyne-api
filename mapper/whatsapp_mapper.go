package mapper

import (
	"go-mnemosyne-api/model"
	"go-mnemosyne-api/whatsapp/dto"
	"strconv"
	"time"
)

func MapPayloadIntoWhatsAppMessageModel(payloadMessageDto *dto.PayloadMessageDto) []model.WhatsappMessage {
	var allWhatsappMessages []model.WhatsappMessage
	for _, entry := range payloadMessageDto.Entry {
		for _, changes := range entry.Changes {
			var whatsAppMessage model.WhatsappMessage
			for _, message := range changes.Value.Messages {
				if message.ID == "" {
					continue
				}
				whatsAppMessage.ID = message.ID
				whatsAppMessage.Type = message.Type
				unixTimestamp, _ := strconv.ParseInt(message.Timestamp, 10, 64)
				whatsAppMessage.Timestamp = time.Unix(unixTimestamp, 0)
				whatsAppMessage.Text = message.Text.Body
				whatsAppMessage.SenderPhoneNumber = message.From
			}
			for _, contact := range changes.Value.Contacts {
				if contact.Profile.Name == "" {
					continue
				}
				whatsAppMessage.Name = contact.Profile.Name
				whatsAppMessage.WhatsAppId = contact.WhatsAppId
			}
			if whatsAppMessage.Text != "" {
				allWhatsappMessages = append(allWhatsappMessages, whatsAppMessage)
			}
		}
	}
	return allWhatsappMessages
}
