package _errors

import "errors"

var (
	ErrInvalidAccessToken = errors.New("invalid access token")
	ErrBadRequest         = errors.New("bad request")
	ErrNotAllUsersReady   = errors.New("not all users are ready")
	ErrNotOwner           = errors.New("owner did not request start")
	ErrNotFirstPlayer     = errors.New("not first player")
	ErrNotLoanCard        = errors.New("not loan card")
)
