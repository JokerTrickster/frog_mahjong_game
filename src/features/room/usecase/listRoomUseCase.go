package usecase

import (
	"context"
	"fmt"
	_interface "main/features/room/model/interface"
	"time"
)

type ListRoomUseCase struct {
	Repository      _interface.IListRoomRepository
	ContextTimeList time.Duration
}

func NewListRoomUseCase(repo _interface.IListRoomRepository, timeList time.Duration) _interface.IListRoomUseCase {
	return &ListRoomUseCase{Repository: repo, ContextTimeList: timeList}
}

func (d *ListRoomUseCase) List(c context.Context, page, pageSize int) error {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeList)
	defer cancel()
	fmt.Println(ctx)

	return nil
}
