package _interface

import (
	"context"
	"main/utils/db/mysql"
)

type IGetUsersRepository interface {
	FindOneUser(ctx context.Context, userID int) (mysql.Users, error)
}

type IListUsersRepository interface {
	FindUsers(ctx context.Context) ([]mysql.Users, error)
	CountUsers(ctx context.Context) (int, error)
}
