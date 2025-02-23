package usecase

import (
	"context"
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

	// 해당 방 유저에 클리어 정보를 가져와야 된다.
	userCorrectPositionsDTO, err := d.Repository.FindFindItUserCorrectPositions(ctx, req.RoomID)
	if err != nil {
		return response.ResFindItResult{}, err
	}

	res := CreateResResult(roomSettingDTO, userCorrectPositionsDTO)

	return res, nil
}
