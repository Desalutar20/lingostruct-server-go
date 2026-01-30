package dto

type SignUpRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8,max=40"`
	FirstName string `json:"firstName" validate:"required,min=3,max=40"`
	LastName  string `json:"lastName" validate:"required,min=3,max=40"`
}
