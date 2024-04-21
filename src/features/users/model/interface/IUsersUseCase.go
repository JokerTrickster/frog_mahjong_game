package _interface

import (
	"context"
	"main/features/users/model/response"
)

type IGetUsersUseCase interface {
	Get(c context.Context, userID int) (response.ResGetUser, error)
}
