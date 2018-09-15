package transport

import (
	"net/http"
)

// CreateReq contains user registration request
type CreateReq struct {
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=8"`
	PasswordConfirm string `json:"password_confirm" validate:"required"`
}

// Bind binds CreateReq
func (x *CreateReq) Bind(r *http.Request) error {
	return nil
}
