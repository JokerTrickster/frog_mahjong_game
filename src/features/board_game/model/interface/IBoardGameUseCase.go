package _interface

import (
	"context"
	"main/features/board_game/model/request"
	"main/features/board_game/model/response"
)

type IFindItSoloPlayBoardGameUseCase interface {
	FindItSoloPlay(c context.Context, userID int, req *request.ReqFindItSoloPlayBoardGame) (response.ResFindItSoloPlayBoardGame, error)
}
