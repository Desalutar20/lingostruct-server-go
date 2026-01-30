package dto

type SignInRequest struct {
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"min=8,max=40"`
}
