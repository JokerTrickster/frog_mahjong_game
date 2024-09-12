package usecase

import (
	"context"
	_interface "main/features/rooms/model/interface"
	"main/features/rooms/model/request"
	"time"
)

type JoinPlayRoomsUseCase struct {
	Repository     _interface.IJoinPlayRoomsRepository
	ContextTimeout time.Duration
}

func NewJoinPlayRoomsUseCase(repo _interface.IJoinPlayRoomsRepository, timeout time.Duration) _interface.IJoinPlayRoomsUseCase {
	return &JoinPlayRoomsUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *JoinPlayRoomsUseCase) JoinPlay(c context.Context, req *request.ReqJoinPlay) error {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()

	// 방 참여 가능한지 체크
	err := d.Repository.FindOneRoom(ctx, req)
	if err != nil {
		return err
	}
	return nil

}
