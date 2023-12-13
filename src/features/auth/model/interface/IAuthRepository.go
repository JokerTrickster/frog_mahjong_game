package _interface

import (
	"context"
	"main/utils/db/mysql"
)

type ISignupAuthRepository interface {
	FindOneUserAuth(ctx context.Context, name string) error
	InsertOneUserDTO(ctx context.Context, userDTO mysql.GormUserDTO) (string, error)
	InsertOneUserAuthDTO(ctx context.Context, userAuthDTO mysql.GormUserAuthDTO) error
}
