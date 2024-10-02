package _interface

import (
	"context"
	"main/features/users/model/request"
	"main/features/users/model/response"
)

type IGetUsersUseCase interface {
	Get(c context.Context, userID int) (response.ResGetUser, error)
}

type IListUsersUseCase interface {
	List(c context.Context) (response.ResListUser, error)
}

type IUpdateUsersUseCase interface {
	Update(c context.Context, userID uint, req *request.ReqUpdateUsers) error
}

type IDeleteUsersUseCase interface {
	Delete(c context.Context, userID uint) error
}
type IListProfilesUsersUseCase interface {
	ListProfiles(c context.Context, userID uint) (response.ResListProfileUser, error)
}
