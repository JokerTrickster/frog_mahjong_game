package usecase

import (
	"context"
	"time"

	_interface "main/features/board_game/model/interface"
	"main/features/board_game/model/request"
	"main/features/board_game/model/response"
)

type SlimeWarResultBoardGameUseCase struct {
	Repository     _interface.ISlimeWarResultBoardGameRepository
	ContextTimeout time.Duration
}

func NewSlimeWarResultBoardGameUseCase(repo _interface.ISlimeWarResultBoardGameRepository, timeout time.Duration) _interface.ISlimeWarResultBoardGameUseCase {
	return &SlimeWarResultBoardGameUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *SlimeWarResultBoardGameUseCase) SlimeWarResult(c context.Context, req *request.ReqSlimeWarResult) (response.ResSlimeWarResult, error) {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()

	gameResultDTOs, err := d.Repository.FindGameResult(ctx, req.RoomID)
	if err != nil {
		return response.ResSlimeWarResult{}, err
	}

	res := response.ResSlimeWarResult{}
	userResult := make([]response.SlimeWarResult, 0)
	for _, gameResultDTO := range gameResultDTOs {
		userResult = append(userResult, response.SlimeWarResult{UserID: gameResultDTO.UserID, Score: gameResultDTO.Score, Result: gameResultDTO.Result})
	}
	res.Users = userResult
	return res, nil
}
