package _interface

import (
	"context"
	"main/utils/db/mysql"
)

type IListProfilesRepository interface {
	FindAllProfiles(ctx context.Context) ([]*mysql.GameProfiles, error)
}

type IUploadProfilesRepository interface {
	InsertOneProfile(ctx context.Context, profile *mysql.GameProfiles) error
}

type IUpdateProfilesRepository interface {
	UpdateOneProfile(ctx context.Context, userID int, profileID int) error
}
