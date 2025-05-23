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

type ISlimeWarGetsCardBoardGameUseCase interface {
	SlimeWarGetsCard(c context.Context) (response.ResSlimeWarGetsCardBoardGame, error)
}

type ISlimeWarResultBoardGameUseCase interface {
	SlimeWarResult(c context.Context, req *request.ReqSlimeWarResult) (response.ResSlimeWarResult, error)
}

type ISlimeWarRankBoardGameUseCase interface {
	SlimeWarRank(c context.Context) (response.ResSlimeWarRankBoardGame, error)
}

type ISequenceResultBoardGameUseCase interface {
	SequenceResult(c context.Context, req *request.ReqSequenceResult) (response.ResSequenceResult, error)
}

type ISequenceRankBoardGameUseCase interface {
	SequenceRank(c context.Context) (response.ResSequenceRank, error)
}

type IGameOverBoardGameUseCase interface {
	GameOver(c context.Context, userID int, req *request.ReqGameOverBoardGame) error
}

type IFrogCardListBoardGameUseCase interface {
	FrogCardList(c context.Context) (response.ResFrogCardListBoardGame, error)
}
