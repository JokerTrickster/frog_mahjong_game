package _interface

import (
	"context"
	"main/features/board_game/model/entity"
	"main/utils/db/mysql"
)

type IFindItSoloPlayBoardGameRepository interface {
	FindRandomImage(ctx context.Context, round int) ([]*mysql.FindItImages, error)
	FindCorrectByImageID(ctx context.Context, imageID uint) ([]*mysql.FindItImageCorrectPositions, error)
}

type IFindItRankBoardGameRepository interface {
	FindTop3UserCorrect(ctx context.Context) ([]*entity.FindItRankEntity, error)
	FindOneUser(ctx context.Context, userID int) (*mysql.GameUsers, error)
}

type IFindItCoinBoardGameRepository interface {
	UpdateUserCoin(ctx context.Context, userID int, coin int) error
}

type IFindItPasswordCheckBoardGameRepository interface {
	FindPasswordCheck(ctx context.Context, password string) (bool, error)
}

type ISlimeWarGetsCardBoardGameRepository interface {
	FindCardList(ctx context.Context) ([]*mysql.SlimeWarCards, error)
}

type ISlimeWarResultBoardGameRepository interface {
	FindGameRoomUser(ctx context.Context, roomID int) ([]*mysql.SlimeWarUsers, error)
	FindRoomMaps(ctx context.Context, roomID int) ([]*mysql.SlimeWarRoomMaps, error)
}

type ISlimeWarRankBoardGameRepository interface {
	FindTop3User(ctx context.Context) ([]*entity.SlimeWarRankEntity, error)
	FindOneUser(ctx context.Context, userID int) (*mysql.GameUsers, error)
}

type ISequenceResultBoardGameRepository interface {
	FindGameRoomUser(ctx context.Context, roomID int) ([]*mysql.SequenceUsers, error)
	FindRoomMaps(ctx context.Context, roomID int) ([]*mysql.SequenceRoomMaps, error)
}
type ISequenceRankBoardGameRepository interface {
	FindTop3User(ctx context.Context) ([]*entity.SequenceRankEntity, error)
	FindOneUser(ctx context.Context, userID int) (*mysql.GameUsers, error)
}

type IGameOverBoardGameRepository interface {
	GameOverInsertGameResult(ctx context.Context, gameResultDTO *mysql.GameResults) error
}
