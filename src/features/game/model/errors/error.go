package _errors

import "errors"

var (
	ErrInvalidAccessToken = errors.New("invalid access token")
	ErrBadRequest         = errors.New("bad request")
	ErrNotAllUsersReady   = errors.New("not all users are ready")
	ErrNotOwner           = errors.New("owner did not request start")
	ErrNotFirstPlayer     = errors.New("not first player")
	ErrNotLoanCard        = errors.New("not loan card")
	ErrNotEnoughCard      = errors.New("not enough card")
	ErrNotEnoughCondition = errors.New("not enough condition to score calculate")
	ErrBadRequestCard     = errors.New("bad request card")
	ErrInvalidGoogleCode  = errors.New("invalid google code")
)
