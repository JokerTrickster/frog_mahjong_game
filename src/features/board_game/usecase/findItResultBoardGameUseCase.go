package usecase

import (
	"context"
	"time"

	_interface "main/features/board_game/model/interface"
	"main/features/board_game/model/request"
	"main/features/board_game/model/response"
)

type FindItResultBoardGameUseCase struct {
	Repository     _interface.IFindItResultBoardGameRepository
	ContextTimeout time.Duration
}

func NewFindItResultBoardGameUseCase(repo _interface.IFindItResultBoardGameRepository, timeout time.Duration) _interface.IFindItResultBoardGameUseCase {
	return &FindItResultBoardGameUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *FindItResultBoardGameUseCase) FindItResult(c context.Context, req *request.ReqFindItResult) (response.ResFindItResult, error) {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()

	userDTOs, err := d.Repository.FindGameRoomUser(ctx, req.RoomID)
	if err != nil {
		return response.ResFindItResult{}, err
	}

	res := response.ResFindItResult{}
	userResult := make([]response.FindItResult, 0)

	for _, userDTO := range userDTOs {
		findItResultDTOs, err := d.Repository.FindFindItResult(ctx, int(userDTO.UserID), req.RoomID)
		if err != nil {
			return response.ResFindItResult{}, err
		}
		userResult = append(userResult, response.FindItResult{UserID: userDTO.UserID, Score: len(findItResultDTOs), Result: 0})
	}
	res.Users = userResult

	return res, nil
}
