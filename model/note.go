package model

import "time"

type Note struct {
	ID         uint64    `gorm:"column:id;primary_key;autoIncrement"`
	UserID     uint64    `gorm:"column:user_id" mapstructure:"user_id"`
	Title      string    `gorm:"column:title" mapstructure:"title"`
	Content    string    `gorm:"column:content" mapstructure:"content"`
	CategoryId uint64    `gorm:"column:category_id" mapstructure:"category_id"`
	Priority   uint64    `gorm:"column:priority" mapstructure:"priority"`
	DueDate    string    `gorm:"column:due_date" mapstructure:"due_date"`
	IsPinned   bool      `gorm:"column:is_pinned" mapstructure:"is_pinned"`
	IsArchived bool      `gorm:"column:is_archived" mapstructure:"is_archived"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at"`
	User       User      `gorm:"foreignKey:user_id;references:id"`
	Category   Category  `gorm:"foreignKey:category_id;references:id"`
}
