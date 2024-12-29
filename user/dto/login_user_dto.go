package dto

type LoginUserDto struct {
	UserIdentifier string `json:"user_identifier" validate:"required"`
	Password       string `json:"password" validate:"required"`
}
