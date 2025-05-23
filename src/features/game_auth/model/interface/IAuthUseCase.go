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

type IRequestPasswordAuthUseCase interface {
	RequestPassword(c context.Context, entity entity.RequestPasswordAuthEntity) (string, error)
}
type IRequestSignupAuthUseCase interface {
	RequestSignup(c context.Context, entity entity.RequestSignupAuthEntity) (string, error)
}
type IValidatePasswordAuthUseCase interface {
	ValidatePassword(c context.Context, entity entity.ValidatePasswordAuthEntity) error
}

type IGoogleOauthCallbackAuthUseCase interface {
	GoogleOauthCallback(c context.Context, code string) (response.ResGameGoogleOauthCallback, error)
}

type IFCMTokenAuthUseCase interface {
	FCMToken(c context.Context, userID uint, req *request.ReqGameFCMToken) error
}

type ICheckSigninAuthUseCase interface {
	CheckSignin(c context.Context, req *request.ReqGameCheckSignin) (bool, error)
}
type INameCheckAuthUseCase interface {
	NameCheck(c context.Context, req *request.ReqGameNameCheck) error
}

type IValidateSignupAuthUseCase interface {
	ValidateSignup(c context.Context, entity entity.ValidateSignupAuthEntity) error
}
