package _interface

import (
	"context"
	"main/features/users/model/entity"
	"main/utils/db/mysql"
)

type IGetUsersRepository interface {
	FindOneUser(ctx context.Context, userID int) (mysql.Users, error)
}

type IListUsersRepository interface {
	FindUsers(ctx context.Context) ([]mysql.Users, error)
	CountUsers(ctx context.Context) (int, error)
}

type IUpdateUsersRepository interface {
	FindOneAndUpdateUsers(ctx context.Context, entitySQL *entity.UpdateUsersEntitySQL) error
}

type IDeleteUsersRepository interface {
	FindOneAndDeleteUsers(ctx context.Context, userID uint) error
}
