package usecase

import (
	"context"
	_interface "main/features/rooms/model/interface"
	"main/features/rooms/model/response"
	"time"
)

type UserListRoomsUseCase struct {
	Repository          _interface.IUserListRoomsRepository
	ContextTimeUserList time.Duration
}

func NewUserListRoomsUseCase(repo _interface.IUserListRoomsRepository, timeUserList time.Duration) _interface.IUserListRoomsUseCase {
	return &UserListRoomsUseCase{Repository: repo, ContextTimeUserList: timeUserList}
}

func (d *UserListRoomsUseCase) UserList(c context.Context, RoomID uint) (response.ResUserListRoom, error) {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeUserList)
	defer cancel()

	userList, err := d.Repository.FindRoomUser(ctx, RoomID)
	if err != nil {
		return response.ResUserListRoom{}, err
	}
	// 방장이 누구인지 체d
	Rooms, err := d.Repository.FindOneRoom(ctx, RoomID)
	if err != nil {
		return response.ResUserListRoom{}, err
	}
	res := CreateResUserListRoom(userList, Rooms)
	return res, nil
}
