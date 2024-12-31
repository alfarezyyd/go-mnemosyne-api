package dto

type UpdateNoteDto struct {
	Title      string `json:"title" validate:"required,min=3,max=100"`
	Content    string `json:"content" validate:"max=255"`
	CategoryId uint64 `json:"category_id" validate:"required,gte=1"`
	Priority   string `json:"priority" validate:"required,oneof=Low Medium High"`
	DueDate    string `json:"due_date" validate:"required,datetime"`
	IsPinned   bool   `json:"is_pinned"`
	IsArchived bool   `json:"is_archived"`
}
