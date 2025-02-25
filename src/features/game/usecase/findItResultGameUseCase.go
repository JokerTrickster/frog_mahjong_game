package usecase

import (
	"context"
	"fmt"
	_interface "main/features/game/model/interface"
	"main/features/game/model/request"
	"main/features/game/model/response"
	"time"
)

type FindItResultGameUseCase struct {
	Repository     _interface.IFindItResultGameRepository
	ContextTimeout time.Duration
}

func NewFindItResultGameUseCase(repo _interface.IFindItResultGameRepository, timeout time.Duration) _interface.IFindItResultGameUseCase {
	return &FindItResultGameUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *FindItResultGameUseCase) FindItResult(c context.Context, req *request.ReqFindItResult) (response.ResFindItResult, error) {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()

	// 게임 정보 가져와야 되고
	roomSettingDTO, err := d.Repository.FindOneRoomSetting(ctx, req.RoomID)
	if err != nil {
		return response.ResFindItResult{}, err
	}
	// 방에 접소한 유저 정보를 가져온다.
	roomUserDTOs, err := d.Repository.FindGameRoomUser(ctx, req.RoomID)
	if err != nil {
		return response.ResFindItResult{}, err
	}
	// 방에 접속한 유저들의 ID를 가져와야 된다.
	var userIDList []int
	for _, roomUserDTO := range roomUserDTOs {
		userIDList = append(userIDList, roomUserDTO.UserID)
	}
	// 해당 방에 접속한 유저들의 정보를 가져와야 된다.
	userDTOs, err := d.Repository.FindGameUsers(ctx, userIDList)
	if err != nil {
		return response.ResFindItResult{}, err
	}

	// 해당 방 유저에 클리어 정보를 가져와야 된다.
	userCorrectPositionsDTO, err := d.Repository.FindFindItUserCorrectPositions(ctx, req.RoomID)
	if err != nil {
		return response.ResFindItResult{}, err
	}

	res := CreateResResult(roomSettingDTO, userCorrectPositionsDTO, userDTOs)

	fmt.Println(res)
	return res, nil
}
