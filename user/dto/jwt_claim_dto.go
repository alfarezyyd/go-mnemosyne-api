package dto

type JwtClaimDto struct {
	Email       *string `json:"email" mapstructure:"email"`
	PhoneNumber *string `json:"phone_number" mapstructure:"phone_number"`
}
