package usecase

import (
	"context"
	"time"

	_interface "main/features/board_game/model/interface"
	"main/features/board_game/model/request"
	"main/features/board_game/model/response"
)

type SequenceResultBoardGameUseCase struct {
	Repository     _interface.ISequenceResultBoardGameRepository
	ContextTimeout time.Duration
}

func NewSequenceResultBoardGameUseCase(repo _interface.ISequenceResultBoardGameRepository, timeout time.Duration) _interface.ISequenceResultBoardGameUseCase {
	return &SequenceResultBoardGameUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *SequenceResultBoardGameUseCase) SequenceResult(c context.Context, req *request.ReqSequenceResult) (response.ResSequenceResult, error) {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()

	gameResultDTOs, err := d.Repository.FindGameResult(ctx, req.RoomID)
	if err != nil {
		return response.ResSequenceResult{}, err
	}

	res := response.ResSequenceResult{}
	userResult := make([]response.SequenceResult, 0)
	for _, gameResultDTO := range gameResultDTOs {
		userResult = append(userResult, response.SequenceResult{UserID: gameResultDTO.UserID, Score: gameResultDTO.Score, Result: gameResultDTO.Result})
	}
	res.Users = userResult

	return res, nil
}
