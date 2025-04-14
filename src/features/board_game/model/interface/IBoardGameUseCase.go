package _interface

import (
	"context"
	"main/features/board_game/model/request"
	"main/features/board_game/model/response"
)

type IFindItSoloPlayBoardGameUseCase interface {
	FindItSoloPlay(c context.Context, userID int, req *request.ReqFindItSoloPlayBoardGame) (response.ResFindItSoloPlayBoardGame, error)
}

type IFindItRankBoardGameUseCase interface {
	FindItRank(c context.Context) (response.ResFindItRankBoardGame, error)
}

type IFindItCoinBoardGameUseCase interface {
	FindItCoin(c context.Context, userID int, req *request.ReqFindItCoinBoardGame) error
}
type IFindItPasswordCheckBoardGameUseCase interface {
	FindItPasswordCheck(c context.Context, req *request.ReqFindItPasswordCheckBoardGame) (bool, error)
}
