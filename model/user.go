package model

import (
	"database/sql"
	"time"
)

type User struct {
	ID                 uint64         `gorm:"column:id;primary_key;auto_increment"`
	Name               string         `gorm:"column:name"`
	Email              string         `gorm:"column:email"`
	Password           string         `gorm:"column:password"`
	Role               string         `gorm:"column:role"`
	EmailVerifiedAt    *time.Time     `gorm:"column:email_verified_at"`
	PhoneNumber        sql.NullString `gorm:"column:phone_number"`
	ProfilePicture     sql.NullString `gorm:"column:profile_picture"`
	IsActive           bool           `gorm:"column:is_active"`
	LanguagePreference string         `gorm:"column:language_preference"`
	CreatedAt          time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt          time.Time      `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
}
