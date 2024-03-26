package usecase

import (
	"context"
	"fmt"
	_interface "main/features/room/model/interface"
	"main/features/room/model/request"
	"time"
)

type ReadyRoomUseCase struct {
	Repository       _interface.IReadyRoomRepository
	ContextTimeReady time.Duration
}

func NewReadyRoomUseCase(repo _interface.IReadyRoomRepository, timeReady time.Duration) _interface.IReadyRoomUseCase {
	return &ReadyRoomUseCase{Repository: repo, ContextTimeReady: timeReady}
}

func (d *ReadyRoomUseCase) Ready(c context.Context, uID uint, req *request.ReqReady) error {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeReady)
	defer cancel()
	fmt.Println(ctx)

	// room user에 player state 를 변경한다.
	err := d.Repository.FindOneAndUpdateRoomUser(ctx, uID, req)
	if err != nil {
		return err
	}

	return nil
}
