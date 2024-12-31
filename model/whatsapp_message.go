package model

import "time"

type WhatsappMessage struct {
	ID                string    `gorm:"column:id;primary_key"`
	Name              string    `gorm:"column:name"`
	WhatsAppId        string    `gorm:"column:whatsapp_id"`
	SenderPhoneNumber string    `gorm:"column:sender_phone_number"`
	Timestamp         time.Time `gorm:"column:timestamp"`
	Type              string    `gorm:"column:type"`
	Text              string    `gorm:"column:text"`
}
