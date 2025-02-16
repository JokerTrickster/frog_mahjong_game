package _interface

import (
	"context"
	"main/features/game_auth/model/entity"
	"main/features/game_auth/model/request"
	"main/features/game_auth/model/response"
)

type ISignupAuthUseCase interface {
	Signup(c context.Context, req *request.ReqGameSignup) error
}

type ISigninAuthUseCase interface {
	Signin(c context.Context, req *request.ReqGameSignin) (response.ResGameSignin, error)
}

type ILogoutAuthUseCase interface {
	Logout(c context.Context, uID uint) error
}

type IReissueAuthUseCase interface {
	Reissue(c context.Context, req *request.ReqGameReissue) (response.ResGameReissue, error)
}

type IGoogleOauthAuthUseCase interface {
	GoogleOauth(c context.Context) (string, error)
}

type IGoogleOauthCallbackAuthUseCase interface {
	GoogleOauthCallback(c context.Context, code string) (response.GameGoogleOauthCallbackRes, error)
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
	V02GoogleOauthCallback(c context.Context, code string) (response.ResGameV02GoogleOauthCallback, error)
}

type IFCMTokenAuthUseCase interface {
	FCMToken(c context.Context, userID uint, req *request.ReqGameFCMToken) error
}

type ICheckSigninAuthUseCase interface {
	CheckSignin(c context.Context, req *request.ReqGameCheckSignin) (bool, error)
}
