package usecase

import (
	"context"
	_interface "main/features/rooms/model/interface"
	"main/features/rooms/model/request"
	"time"
)

type ReadyRoomsUseCase struct {
	Repository       _interface.IReadyRoomsRepository
	ContextTimeReady time.Duration
}

func NewReadyRoomsUseCase(repo _interface.IReadyRoomsRepository, timeReady time.Duration) _interface.IReadyRoomsUseCase {
	return &ReadyRoomsUseCase{Repository: repo, ContextTimeReady: timeReady}
}

func (d *ReadyRoomsUseCase) Ready(c context.Context, uID uint, req *request.ReqReady) error {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeReady)
	defer cancel()

	// Rooms user에 player state 를 변경한다.
	err := d.Repository.FindOneAndUpdateRoomUser(ctx, uID, req)
	if err != nil {
		return err
	}

	return nil
}
