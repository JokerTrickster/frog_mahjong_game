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
	FindOneUser(ctx context.Context, email, password string) (mysql.Users, error)
}
