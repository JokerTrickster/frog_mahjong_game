package _interface

import (
	"context"
	"main/utils/db/mysql"
)

type IListProfilesRepository interface {
	FindAllProfiles(ctx context.Context, userID uint) ([]*mysql.UserProfiles, error)
}
