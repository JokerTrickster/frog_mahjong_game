package usecase

import (
	"context"
	_interface "main/features/room/model/interface"
	"main/features/room/model/response"
	"time"
)

type ListRoomUseCase struct {
	Repository      _interface.IListRoomRepository
	ContextTimeList time.Duration
}

func NewListRoomUseCase(repo _interface.IListRoomRepository, timeList time.Duration) _interface.IListRoomUseCase {
	return &ListRoomUseCase{Repository: repo, ContextTimeList: timeList}
}

func (d *ListRoomUseCase) List(c context.Context, page, pageSize int) (response.ResListRoom, error) {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeList)
	defer cancel()
	rooms, err := d.Repository.FindRoomList(ctx, page, pageSize)
	if err != nil {
		return response.ResListRoom{}, err
	}
	total, err := d.Repository.CountRoomList(ctx)
	if err != nil {
		return response.ResListRoom{}, err
	}

	//create res
	res, err := CreateResListRoom(rooms, total)
	if err != nil {
		return response.ResListRoom{}, err
	}

	return res, nil
}
