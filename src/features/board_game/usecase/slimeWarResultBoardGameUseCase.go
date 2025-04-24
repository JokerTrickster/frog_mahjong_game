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

	// 방에 접소한 유저 정보를 가져온다.
	roomUserDTOs, err := d.Repository.FindGameRoomUser(ctx, req.RoomID)
	if err != nil {
		return response.ResSlimeWarResult{}, err
	}

	// 맵 정보 가져오기
	maps, err := d.Repository.FindRoomMaps(ctx, req.RoomID)
	if err != nil {
		return response.ResSlimeWarResult{}, err
	}

	res := response.ResSlimeWarResult{}
	userResult := make([]response.SlimeWarResult, 0)
	// 유저 별로 맵 정보 기반으로 점수를 계산한다.
	for _, roomUserDTO := range roomUserDTOs {
		userID := roomUserDTO.UserID
		userScore := 0
		for _, mapDTO := range maps {
			if mapDTO.UserID == userID {
				userScore += 1
			}
		}
		userResult = append(userResult, response.SlimeWarResult{UserID: userID, Score: userScore})
	}
	res.Result = userResult

	return res, nil
}
