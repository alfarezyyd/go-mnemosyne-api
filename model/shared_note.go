package model

import "time"

type SharedNote struct {
	ID         uint64    `gorm:"column:ID;primary_key;autoIncrement"`
	NoteId     uint64    `gorm:"column:note_id;"`
	UserId     uint64    `gorm:"column:user_id"`
	Permission string    `gorm:"column:permission"`
	SharedAt   time.Time `gorm:"column:shared_at"`
	ExpiresAt  time.Time `gorm:"column:expires_at"`
	Status     string    `gorm:"column:status"`
}
