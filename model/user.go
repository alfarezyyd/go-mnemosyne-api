package model

import (
	"database/sql"
	"time"
)

type User struct {
	ID                   uint64                 `gorm:"column:id;primary_key;autoIncrement"`
	Name                 string                 `gorm:"column:name" mapstructure:"name"`
	Email                string                 `gorm:"column:email" mapstructure:"email"`
	Password             string                 `gorm:"column:password" mapstructure:"password"`
	Role                 string                 `gorm:"column:role;default:'User'"`
	EmailVerifiedAt      *time.Time             `gorm:"column:email_verified_at"`
	PhoneNumber          sql.NullString         `gorm:"column:phone_number" mapstructure:"phone_number"`
	ProfilePicture       sql.NullString         `gorm:"column:profile_picture"`
	IsActive             bool                   `gorm:"column:is_active"`
	LanguagePreference   string                 `gorm:"column:language_preference"`
	CreatedAt            time.Time              `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt            time.Time              `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	OneTimePasswordToken []OneTimePasswordToken `gorm:"foreignKey:user_id;references:id"`
	Categories           []Category             `gorm:"foreignKey:user_id;references:id"`
	Note                 []Note                 `gorm:"foreignKey:category_id;references:id"`
}
