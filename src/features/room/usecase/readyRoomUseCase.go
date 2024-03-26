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
	return nil
}
