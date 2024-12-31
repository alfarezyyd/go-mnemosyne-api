package model

import "time"

type Note struct {
	ID         uint64    `gorm:"column:id;primary_key;autoIncrement"`
	UserID     uint64    `gorm:"column:user_id" mapstructure:"user_id"`
	Title      string    `gorm:"column:title"`
	Content    string    `gorm:"column:content"`
	CategoryId uint64    `gorm:"column:category_id"`
	Priority   string    `gorm:"column:priority;default:Low"`
	DueDate    string    `gorm:"column:due_date"`
	IsPinned   bool      `gorm:"column:is_pinned"`
	IsArchived bool      `gorm:"column:is_archived"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at"`
	User       User      `gorm:"foreignKey:user_id;references:id"`
	Category   Category  `gorm:"foreignKey:category_id;references:id"`
}
