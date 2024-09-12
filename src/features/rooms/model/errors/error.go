package _errors

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExisted = errors.New("user already existed")
	ErrBadRequest         = errors.New("bad request")
	ErrRoomNotFound       = errors.New("room not found")
	ErrRoomFull           = errors.New("방이 꽉 찼습니다.")
	ErrWrongPassword      = errors.New("비밀번호가 일치하지 않습니다.")
	ErrPlayerStateFailed  = errors.New("player state change failed")
	ErrRoomUserNotFound   = errors.New("room user not found")
	ErrServerError        = errors.New("server error")
)
