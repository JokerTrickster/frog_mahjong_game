package _interface

import (
	"context"
	"main/features/auth/model/entity"
	"main/utils/db/mysql"
)

type ISignupAuthRepository interface {
	UserCheckByEmail(ctx context.Context, email string) error
	InsertOneUser(ctx context.Context, user *mysql.Users) error
	VerifyAuthCode(ctx context.Context, email, code string) error
	FindAllBasicProfile(ctx context.Context) ([]*mysql.Profiles, error)
	InsertOneUserProfile(ctx context.Context, userProfileDTOList []*mysql.UserProfiles) error
}

type ISigninAuthRepository interface {
	FindOneAndUpdateUser(ctx context.Context, email, password string) (mysql.Users, error)
	SaveToken(ctx context.Context, uID uint, accessToken, refreshToken string, refreshTknExpiredAt int64) error
	DeleteToken(ctx context.Context, uID uint) error
}

type ILogoutAuthRepository interface {
	FindOneAndUpdateUser(ctx context.Context, uID uint) error
	DeleteToken(ctx context.Context, uID uint) error
}

type IReissueAuthRepository interface {
	SaveToken(ctx context.Context, token mysql.Tokens) error
	DeleteToken(ctx context.Context, uID uint) error
}

type IGoogleOauthAuthRepository interface {
}

type IGoogleOauthCallbackAuthRepository interface {
	FindOneAndUpdateUser(ctx context.Context, googleOauthCallbackSQLQuery *entity.GoogleOauthCallbackSQLQuery) (*mysql.Users, error)
	SaveToken(ctx context.Context, uID uint, accessToken, refreshToken string, refreshTknExpiredAt int64) error
	DeleteToken(ctx context.Context, uID uint) error
	CreateUser(ctx context.Context, user *mysql.Users) (*mysql.Users, error)
}

type IRequestPasswordAuthRepository interface {
	FindOneUserByEmail(ctx context.Context, email string) error
	InsertAuthCode(ctx context.Context, userAuthDTO mysql.UserAuths) error
}

type IRequestSignupAuthRepository interface {
	FindOneUserByEmail(ctx context.Context, email string) error
	InsertAuthCode(ctx context.Context, userAuthDTO mysql.UserAuths) error
	DeleteAuthCodeByEmail(ctx context.Context, email string) error
}
type IValidatePasswordAuthRepository interface {
	CheckAuthCode(ctx context.Context, email, code string) error
	UpdatePassword(ctx context.Context, user mysql.Users) error
	DeleteAuthCode(ctx context.Context, email string) error
}

type IV02GoogleOauthCallbackAuthRepository interface {
	FindOneAndUpdateUser(ctx context.Context, googleOauthCallbackSQLQuery *entity.V02GoogleOauthCallbackSQLQuery) (*mysql.Users, error)
	SaveToken(ctx context.Context, uID uint, accessToken, refreshToken string, refreshTknExpiredAt int64) error
	DeleteToken(ctx context.Context, uID uint) error
	CreateUser(ctx context.Context, user *mysql.Users) (*mysql.Users, error)
	FindAllBasicProfile(ctx context.Context) ([]*mysql.Profiles, error)
	InsertOneUserProfile(ctx context.Context, userProfileDTOList []*mysql.UserProfiles) error
}
