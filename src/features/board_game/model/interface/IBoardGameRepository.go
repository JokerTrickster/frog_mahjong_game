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
