package _interface

import (
	"context"
	"main/features/auth/model/entity"
	"main/utils/db/mysql"
)

type ISignupAuthRepository interface {
	UserCheckByEmail(ctx context.Context, email string) error
	InsertOneUser(ctx context.Context, user mysql.Users) error
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
