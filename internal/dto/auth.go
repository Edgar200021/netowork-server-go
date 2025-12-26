package dto

import "github.com/Edgar200021/netowork-server-go/.gen/netowork/public/model"

type SignInRequest struct {
	Email    string `json:"email" validate:"required,email,lte=100"`
	Password string `json:"password" validate:"required,gte=8,lte=40"`
}

type SignInResponse struct {
	UserResponse
	SessionId string
}

type SignUpRequest struct {
	FirstName string         `json:"firstName" validate:"required,gte=3,lte=40"`
	LastName  string         `json:"lastName" validate:"required,gte=4,lte=40"`
	Email     string         `json:"email" validate:"required,email,lte=100"`
	Password  string         `json:"password" validate:"required,gte=8,lte=40"`
	Role      model.UserRole `json:"role" validate:"required,oneof=client freelancer"`
}

type VerifyAccountRequest struct {
	Token string `json:"token" validate:"required"`
	Email string `json:"email" validate:"required,email,lte=100"`
}

type VerifyAccountResponse struct {
	UserResponse
	SessionId string
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email,lte=100"`
}

type ResetPasswordRequest struct {
	Email    string `json:"email" validate:"required,email,lte=100"`
	Password string `json:"password" validate:"required,gte=8,lte=40"`
	Token    string `json:"token" validate:"required"`
}
