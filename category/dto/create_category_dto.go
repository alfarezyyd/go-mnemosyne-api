package dto

type CreateCategoryDto struct {
	Name        string `json:"name" validate:"required,min=1,max=100"`
	Description string `json:"description" validate:"max=255"`
}
