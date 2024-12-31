package dto

type CreateNoteDto struct {
	Title      string `json:"title" validate:"required,min=3,max=100"`
	Content    string `json:"content" validate:"max=255"`
	CategoryId string `json:"category_id" validate:"required,gte=1"`
	Priority   string `json:"priority" validate:"required,oneof=Low Medium High"`
	DueDate    string `json:"due_date" validate:"required,datetime"`
	IsPinned   bool   `json:"is_pinned" validate:"required"`
	IsArchived bool   `json:"is_archived" validate:"required"`
}
