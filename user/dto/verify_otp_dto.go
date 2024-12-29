package dto

type VerifyOtpDto struct {
	Email                string `json:"email" validate:"required,email"`
	OneTimePasswordToken string `json:"one_time_password_token" validate:"required,len=4"`
}
