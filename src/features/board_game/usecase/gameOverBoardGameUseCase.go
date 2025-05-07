package usecase

import (
	"context"
	_interface "main/features/board_game/model/interface"
	"main/features/board_game/model/request"
	"main/utils/db/mysql"
	"time"
)

type GameOverBoardGameUseCase struct {
	Repository     _interface.IGameOverBoardGameRepository
	ContextTimeout time.Duration
}

func NewGameOverBoardGameUseCase(repo _interface.IGameOverBoardGameRepository, timeout time.Duration) _interface.IGameOverBoardGameUseCase {
	return &GameOverBoardGameUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *GameOverBoardGameUseCase) GameOver(c context.Context, userID int, req *request.ReqGameOverBoardGame) error {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()
	gameResultDTO := &mysql.GameResults{
		UserID:   userID,
		GameType: req.GameType,
		RoomID:   req.RoomID,
	}
	if req.Winner {
		gameResultDTO.Result = 1
	} else {
		gameResultDTO.Result = 0
	}
	err := d.Repository.GameOverInsertGameResult(ctx, gameResultDTO)
	if err != nil {
		return err
	}

	return nil
}
