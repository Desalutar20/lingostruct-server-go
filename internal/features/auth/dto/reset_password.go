package dto

type ResetPasswordRequest struct {
	Email       string `json:"email" validate:"required,email"`
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required,min=8,max=40"`
}
