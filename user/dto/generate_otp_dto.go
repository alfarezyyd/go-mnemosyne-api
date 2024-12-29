package dto

type GenerateOtpDto struct {
	Email string `json:"email" validate:"required,email"`
}
