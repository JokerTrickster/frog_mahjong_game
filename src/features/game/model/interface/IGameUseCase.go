package _interface

import (
	"context"
	"main/features/game/model/entity"
	"main/features/game/model/request"
	"main/features/game/model/response"
)

type IStartGameUseCase interface {
	Start(c context.Context, uID uint, req *request.ReqStart) error
}

type IDoraGameUseCase interface {
	Dora(c context.Context, userID int, req *request.ReqDora) error
}
type IOwnershipGameUseCase interface {
	Ownership(c context.Context, req *request.ReqOwnership) error
}

type IDiscardGameUseCase interface {
	Discard(c context.Context, userID int, req *request.ReqDiscard) error
}

type INextTurnGameUseCase interface {
	NextTurn(c context.Context, req *request.ReqNextTurn) error
}

type ILoanGameUseCase interface {
	Loan(c context.Context, userID uint, req *request.ReqLoan) error
}

type IScoreCalculateGameUseCase interface {
	ScoreCalculate(c context.Context, userID uint, req *request.ReqScoreCalculate) (int, []string, error)
}

type IWinRequestGameUseCase interface {
	WinRequest(c context.Context, req *request.ReqWinRequest) (bool, error)
}

type IResultGameUseCase interface {
	Result(c context.Context, userID uint, req *request.ReqResult) (response.ResResult, error)
}

type IReportGameUseCase interface {
	Report(c context.Context, userID uint, req *request.ReqReport) error
}

type IMetaGameUseCase interface {
	Meta(c context.Context) (response.ResMetaGame, error)
}
type IDeckCardGameUseCase interface {
	DeckCard(c context.Context, userID, roomID int) (response.ResDeckCardGame, error)
}

type IListMissionGameUseCase interface {
	ListMission(c context.Context) (response.ResListMissionGame, error)
}

type ICreateMissionGameUseCase interface {
	CreateMission(c context.Context, req *request.ReqCreateMission) error
}
type IListCardGameUseCase interface {
	ListCard(c context.Context) (response.ResListCardGame, error)
}
type IV2DeckCardGameUseCase interface {
	V2DeckCard(c context.Context, userID, roomID int) (response.ResV2DeckCardGame, error)
}
type IV2ReportGameUseCase interface {
	V2Report(c context.Context, userID uint, req *request.ReqV2Report) error
}
type IV2ResultGameUseCase interface {
	V2Result(c context.Context, req *request.ReqV2Result) (response.ResV2Result, error)
}

type ISaveCardInfoGameUseCase interface {
	SaveCardInfo(c context.Context, req *request.ReqSaveCardInfo) error
}

type ISaveCardImageGameUseCase interface {
	SaveCardImage(c context.Context, e entity.SaveCardImageGameEntity) error
}

type IUpdateCardGameUseCase interface {
	UpdateCard(c context.Context, req *request.ReqUpdateCard) error
}

type IReportImageUploadGameUseCase interface {
	ReportImageUpload(c context.Context, req *request.ReqReportImageUploadGame) error
}
