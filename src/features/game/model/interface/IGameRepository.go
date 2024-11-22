package _interface

import (
	"context"
	"main/features/game/model/entity"
	"main/features/game/model/request"
	"main/utils/db/mysql"
)

type IStartGameRepository interface {
	CheckOwner(c context.Context, uID uint, roomID uint) error
	CheckReady(c context.Context, roomID uint) ([]mysql.RoomUsers, error)
	UpdateRoomUser(c context.Context, updateRoomUsers []mysql.RoomUsers) error
	UpdateRoom(c context.Context, roomID uint, state string) error
	CreateCards(c context.Context, roomID uint, cards []mysql.Cards) error
}

type IDoraGameRepository interface {
	CheckFirstPlayer(c context.Context, userID int, roomID int) error
	UpdateDoraCard(c context.Context, req *request.ReqDora) error
}
type IOwnershipGameRepository interface {
	UpdateCardState(c context.Context, req *request.ReqOwnership) error
	UpdateRoomUserCardCount(c context.Context, req *request.ReqOwnership) error
}

type IDiscardGameRepository interface {
	PlayerCheckTurn(c context.Context, req *request.ReqDiscard) (mysql.RoomUsers, error)
	UpdateCardStateDiscard(c context.Context, req *request.ReqDiscard) error
	UpdateRoomUser(c context.Context, updateRoomUser mysql.RoomUsers) error
}

type INextTurnGameRepository interface {
	UpdatePlayerNextTurn(c context.Context, req *request.ReqNextTurn) error
}

type ILoanGameRepository interface {
	CheckLoan(c context.Context, req *request.ReqLoan) error
	Loan(c context.Context, req *request.ReqLoan) error
	UpdateRoomUserCardCount(c context.Context, userID uint, roomID uint) error
}

type IScoreCalculateGameRepository interface {
	FindOwnedCards(c context.Context, entitySQL *entity.ScoreCalculateEntitySQL) ([]mysql.Cards, error)
	GetDoraCard(c context.Context, req *request.ReqScoreCalculate) (mysql.Cards, error)
}

type IWinRequestGameRepository interface {
	GetRoomUser(c context.Context, userID uint, roomID uint) (mysql.RoomUsers, error)
}
type IResultGameRepository interface {
	FindCards(c context.Context, entitySQL *entity.ResultEntitySQL) ([]mysql.Cards, error)
	GetDoraCard(c context.Context, req *request.ReqResult) (mysql.Cards, error)
}

type IReportGameRepository interface {
	SaveReport(c context.Context, reportDTO *mysql.Reports) error
}
type IMetaGameRepository interface {
	FindAllReportCategory(c context.Context) ([]mysql.Categories, error)
}
type IDeckCardGameRepository interface {
	CheckRoomUser(c context.Context, userID int, roomID int) error
}

type IListMissionGameRepository interface {
	FindAllMission(c context.Context) ([]*mysql.Missions, error)
}

type ICreateMissionGameRepository interface {
	SaveMission(c context.Context, missionDTO *mysql.Missions) error
}
type IListCardGameRepository interface {
	FindAllBirdCard(c context.Context) ([]*mysql.BirdCards, error)
	CountAllBirdCard(c context.Context) (int, error)
}

// v2
type IV2DeckCardGameRepository interface {
	CheckRoomUser(c context.Context, userID int, roomID int) error
}

type IV2ReportGameRepository interface {
	SaveReport(c context.Context, reportDTO *mysql.Reports) error
}
type IV2ResultGameRepository interface {
	GetUserMissions(c context.Context, req *request.ReqV2Result) ([]*mysql.UserMissions, error)
	GetUserMissionCards(c context.Context, missionID uint) ([]*mysql.UserMissionCards, error)
}

type ISaveCardInfoGameRepository interface {
	SaveCardInfo(c context.Context, birdCardsDTO []mysql.BirdCards) error
}

type ISaveCardImageGameRepository interface {
}

type IUpdateCardGameRepository interface {
	UpdateCard(c context.Context, updates mysql.BirdCards) error
}
