package _interface

import (
	"context"
	"main/features/auth/model/entity"
	"main/features/auth/model/request"
	"main/features/auth/model/response"
)

type ISignupAuthUseCase interface {
	Signup(c context.Context, req *request.ReqSignup) error
}

type ISigninAuthUseCase interface {
	Signin(c context.Context, req *request.ReqSignin) (response.ResSignin, error)
}

type ILogoutAuthUseCase interface {
	Logout(c context.Context, uID uint) error
}

type IReissueAuthUseCase interface {
	Reissue(c context.Context, req *request.ReqReissue) (response.ResReissue, error)
}

type IGoogleOauthAuthUseCase interface {
	GoogleOauth(c context.Context) (string, error)
}

type IGoogleOauthCallbackAuthUseCase interface {
	GoogleOauthCallback(c context.Context, code string) (response.GoogleOauthCallbackRes, error)
}

type IRequestPasswordAuthUseCase interface {
	RequestPassword(c context.Context, entity entity.RequestPasswordAuthEntity) (string, error)
}
type IRequestSignupAuthUseCase interface {
	RequestSignup(c context.Context, entity entity.RequestSignupAuthEntity) (string, error)
}
type IValidatePasswordAuthUseCase interface {
	ValidatePassword(c context.Context, entity entity.ValidatePasswordAuthEntity) error
}

type IV02GoogleOauthCallbackAuthUseCase interface {
	V02GoogleOauthCallback(c context.Context, code string) (response.ResV02GoogleOauthCallback, error)
}

type IFCMTokenAuthUseCase interface {
	FCMToken(c context.Context, userID uint, req *request.ReqFCMToken) error
}
