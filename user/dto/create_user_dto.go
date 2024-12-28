package dto

type CreateUserDto struct {
	Name            string `json:"name" validate:"required,min=3,max=100"`
	Email           string `json:"email" validate:"required,email"`
	PhoneNumber     string `json:"phone_number" validate:"required,number"`
	Password        string `json:"password" validate:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=6,eqfield=Password"`
}
