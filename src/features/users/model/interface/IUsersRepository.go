package _interface

import (
	"context"
	"main/features/users/model/entity"
	"main/features/users/model/request"
	"main/utils/db/mysql"
)

type IGetUsersRepository interface {
	FindOneUser(ctx context.Context, userID int) (mysql.Users, error)
	CheckDisconnect(ctx context.Context, userID int) (int64, error)
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
	DeleteToken(ctx context.Context, userID uint) error
}

type IListProfilesUsersRepository interface {
	FindAllProfiles(ctx context.Context, userID uint) ([]*mysql.UserProfiles, error)
}

type IFullCoinUsersRepository interface {
	FullCoin(ctx context.Context) error
}

type IOneCoinUsersRepository interface {
	OneCoin(ctx context.Context) error
}

type IAlertUsersRepository interface {
	FindOneAndUpdateUsers(ctx context.Context, userID uint, req *request.ReqAlertUsers) error
}

type IPushUsersRepository interface {
	FindUsersForNotifications(ctx context.Context) ([]mysql.Users, error)
	FindOnePushToken(ctx context.Context, userID uint) (string, error)
}
