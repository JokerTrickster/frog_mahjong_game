package _interface

import (
	"context"
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
