package autherrors

import "errors"

var (
	ErrUserNotFound                      = errors.New("user not found")
	ErrUserAlreadyExists                 = errors.New("user already exists")
	ErrVerificationTokenInvalidOrExpired = errors.New(
		"verification link is invalid or expired",
	)
	ErrUserBanned                         = errors.New("user is banned")
	ErrUserNotVerified                    = errors.New("user not verified")
	ErrUserAlreadyVerified                = errors.New("user already verified")
	ErrInvalidCredentials                 = errors.New("invalid credentials")
	ErrUnauthorized                       = errors.New("unauthorized")
	ErrResetPasswordTokenInvalidOrExpired = errors.New("reset password link is invalid or expireв")
)
