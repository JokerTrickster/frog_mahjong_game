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
