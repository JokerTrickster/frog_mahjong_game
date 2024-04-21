package usecase

import (
	"context"
	_interface "main/features/rooms/model/interface"
	"main/features/rooms/model/response"
	"time"
)

type ListRoomsUseCase struct {
	Repository      _interface.IListRoomsRepository
	ContextTimeList time.Duration
}

func NewListRoomsUseCase(repo _interface.IListRoomsRepository, timeList time.Duration) _interface.IListRoomsUseCase {
	return &ListRoomsUseCase{Repository: repo, ContextTimeList: timeList}
}

func (d *ListRoomsUseCase) List(c context.Context, page, pageSize int) (response.ResListRoom, error) {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeList)
	defer cancel()
	Rooms, err := d.Repository.FindRoomList(ctx, page, pageSize)
	if err != nil {
		return response.ResListRoom{}, err
	}
	total, err := d.Repository.CountRoomList(ctx)
	if err != nil {
		return response.ResListRoom{}, err
	}

	//create res
	res, err := CreateResListRoom(Rooms, total)
	if err != nil {
		return response.ResListRoom{}, err
	}

	return res, nil
}
