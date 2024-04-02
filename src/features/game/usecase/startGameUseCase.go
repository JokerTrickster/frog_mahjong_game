package usecase

import (
	"context"
	"fmt"
	_interface "main/features/game/model/interface"
	"main/features/game/model/request"
	"time"
)

type StartGameUseCase struct {
	Repository     _interface.IStartGameRepository
	ContextTimeout time.Duration
}

func NewStartGameUseCase(repo _interface.IStartGameRepository, timeout time.Duration) _interface.IStartGameUseCase {
	return &StartGameUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *StartGameUseCase) Start(c context.Context, email string, req *request.ReqStart) error {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()
	fmt.Println(ctx)
	// 방장이 게임 시작 요청했는지 체크
	err := d.Repository.CheckOwner(ctx, email, req.RoomID)
	if err != nil {
		return err
	}

	// 방에 있는 유저들이 모두 레디 상태인지 확인
	roomUsers, err := d.Repository.CheckReady(ctx, req.RoomID)
	if err != nil {
		return err
	}
	if allReady := CheckRoomUsersReady(roomUsers); !allReady {
		return fmt.Errorf("not all users are ready")
	}

	// room user 데이터 변경 (대기 -> 플레이)
	err = d.Repository.UpdateRoomUser(ctx, req.RoomID, req.State)
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
