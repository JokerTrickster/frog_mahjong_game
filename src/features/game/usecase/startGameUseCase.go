package usecase

import (
	"context"
	_interface "main/features/game/model/interface"
	"main/features/game/model/request"
	"main/utils"
	"time"
)

type StartGameUseCase struct {
	Repository     _interface.IStartGameRepository
	ContextTimeout time.Duration
}

func NewStartGameUseCase(repo _interface.IStartGameRepository, timeout time.Duration) _interface.IStartGameUseCase {
	return &StartGameUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *StartGameUseCase) Start(c context.Context, uID uint, req *request.ReqStart) error {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()
	// 방장이 게임 시작 요청했는지 체크
	err := d.Repository.CheckOwner(ctx, uID, req.RoomID)
	if err != nil {
		return err
	}

	// 방에 있는 유저들이 모두 레디 상태인지 확인
	roomUsers, err := d.Repository.CheckReady(ctx, req.RoomID)
	if err != nil {
		return err
	}
	if allReady := CheckRoomUsersReady(roomUsers); !allReady {
		return utils.ErrorMsg(ctx, utils.ErrBadRequest, utils.Trace(), "All users are not ready", utils.ErrFromClient)
	}

	// room user 데이터 변경 (대기 -> 플레이, 플레이 순번 랜덤으로 생성)
	updatedRoomUsers, err := StartUpdateRoomUsers(roomUsers)
	if err != nil {
		return err
	}
	// room user 데이터 변경 (대기 -> 플레이)
	err = d.Repository.UpdateRoomUser(ctx, updatedRoomUsers)
	if err != nil {
		return err
	}

	// room 데이터 상태 변경 (대기 -> 플레이)
	err = d.Repository.UpdateRoom(ctx, req.RoomID, req.State)
	if err != nil {
		return err
	}

	// cards 데이터 생성
	cards := CreateInitCards(req.RoomID)
	err = d.Repository.CreateCards(ctx, req.RoomID, cards)
	if err != nil {
		return err
	}

	return nil
}
