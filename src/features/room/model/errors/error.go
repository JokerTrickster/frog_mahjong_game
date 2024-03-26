package _errors

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidAccessToken = errors.New("invalid access token")
	ErrUserAlreadyExisted = errors.New("user already existed")
	ErrBadRequest         = errors.New("bad request")
	ErrRoomNotFound       = errors.New("room not found")
	ErrRoomFull           = errors.New("room full")
	ErrPlayerStateFailed  = errors.New("player state change failed")
	ErrRoomUserNotFound   = errors.New("room user not found")
)
