package _interface

import (
	"context"
	"main/features/profiles/model/response"
)

type IListProfilesUseCase interface {
	List(c context.Context, userID uint) (response.ResListProfile, error)
}
