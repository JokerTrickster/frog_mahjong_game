package _interface

import (
	"context"
	"main/features/auth/model/request"
)

type ISignupAuthUseCase interface {
	Signup(c context.Context, req *request.ReqSignup) error
}

type ISigninAuthUseCase interface {
	Signin(c context.Context, req *request.ReqSignin) error
}