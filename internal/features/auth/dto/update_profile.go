package dto

import "mime/multipart"

type UpdateProfileRequest struct {
	FirstName string `json:",omitempty" validate:"min=3,max=40"`
	LastName  string `json:",omitempty" validate:"min=3,max=40"`
	Image     multipart.File
}
