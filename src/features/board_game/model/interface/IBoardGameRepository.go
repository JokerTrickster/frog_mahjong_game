package _interface

import (
	"context"
	"main/utils/db/mysql"
)

type IFindItSoloPlayBoardGameRepository interface {
	FindRandomImage(ctx context.Context, round int) ([]*mysql.FindItImages, error)
	FindCorrectByImageID(ctx context.Context, imageID uint) ([]*mysql.FindItImageCorrectPositions, error)
}
