package mapper

import (
	"fmt"
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
				whatsAppMessage.SenderPhoneNumber = message.From
				switch message.Type {
				case "text":
					whatsAppMessage.Text = &(message.Text.Body)
					break
				case "image":
					whatsAppMessage.MimeType = &(message.Media.MimeType)
					whatsAppMessage.SHA256 = &(message.Media.SHA256)
					whatsAppMessage.MediaId = &(message.Media.ID)
					break
				}
			}
			for _, contact := range changes.Value.Contacts {
				if contact.Profile.Name == "" {
					continue
				}
				whatsAppMessage.Name = contact.Profile.Name
				whatsAppMessage.WhatsAppId = contact.WhatsAppId
			}
			fmt.Println(*whatsAppMessage.MediaId)
			if (whatsAppMessage.Type == "text" && *(whatsAppMessage.Text) != "") || (whatsAppMessage.Type == "image" && *(whatsAppMessage.MediaId) != "") {
				allWhatsappMessages = append(allWhatsappMessages, whatsAppMessage)
			}
		}
	}
	return allWhatsappMessages
}
