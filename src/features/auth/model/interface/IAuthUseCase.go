package _interface

import (
	"context"
	"main/features/auth/model/request"
)

type ISignupAuthUseCase interface {
	Signup(c context.Context, req *request.ReqSignup) error
}
