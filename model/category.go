package model

import "time"

type Category struct {
	ID          uint64    `gorm:"primary_key;autoIncrement"`
	UserID      uint64    `gorm:"column:user_id" mapstructure:"user_id"`
	User        *User     `gorm:"foreignKey:user_id;references:id"`
	Name        string    `gorm:"column:name" mapstructure:"name"`
	Description string    `gorm:"column:description" mapstructure:"description"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	Note        []Note    `gorm:"foreignKey:category_id;references:id"`
}
