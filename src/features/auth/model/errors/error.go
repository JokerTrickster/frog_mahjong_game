package _errors

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidAccessToken = errors.New("invalid access token")
	ErrUserAlreadyExisted = errors.New("user already existed")
	ErrUserGoogleExisted  = errors.New("user google already existed")
	ErrCodeNotFound       = errors.New("code not found")
	ErrPasswordNotMatch   = errors.New("password not match")
)
