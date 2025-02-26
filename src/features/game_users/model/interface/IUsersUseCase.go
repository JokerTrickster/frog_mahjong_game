package _interface

import (
	"context"
	"main/features/game_users/model/request"
	"main/features/game_users/model/response"
)

type IGetUsersUseCase interface {
	Get(c context.Context, userID int) (response.ResGetGameUser, error)
}

type IListUsersUseCase interface {
	List(c context.Context) (response.ResListGameUser, error)
}

type IUpdateUsersUseCase interface {
	Update(c context.Context, userID uint, req *request.ReqUpdateGameUsers) error
}

type IDeleteUsersUseCase interface {
	Delete(c context.Context, userID uint) error
}
type IListProfilesUsersUseCase interface {
	ListProfiles(c context.Context, userID uint) (response.ResListProfileGameUser, error)
}
type IFullCoinUsersUseCase interface {
	FullCoin(c context.Context) error
}
type IOneCoinUsersUseCase interface {
	OneCoin(c context.Context) error
}

type IAlertUsersUseCase interface {
	Alert(c context.Context, userID uint, e *request.ReqAlertGameUsers) error
}

type IPushUsersUseCase interface {
	Push(c context.Context, req *request.ReqPushGameUsers) error
}
