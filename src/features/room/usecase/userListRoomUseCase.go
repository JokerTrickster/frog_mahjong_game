package usecase

import (
	"context"
	_interface "main/features/room/model/interface"
	"main/features/room/model/response"
	"time"
)

type UserListRoomUseCase struct {
	Repository          _interface.IUserListRoomRepository
	ContextTimeUserList time.Duration
}

func NewUserListRoomUseCase(repo _interface.IUserListRoomRepository, timeUserList time.Duration) _interface.IUserListRoomUseCase {
	return &UserListRoomUseCase{Repository: repo, ContextTimeUserList: timeUserList}
}

func (d *UserListRoomUseCase) UserList(c context.Context, roomID uint) (response.ResUserListRoom, error) {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeUserList)
	defer cancel()

	userList, err := d.Repository.FindRoomUser(ctx, roomID)
	if err != nil {
		return response.ResUserListRoom{}, err
	}
	// 방장이 누구인지 체d
	rooms, err := d.Repository.FindOneRoom(ctx, roomID)
	if err != nil {
		return response.ResUserListRoom{}, err
	}
	res := CreateResUserListRoom(userList, rooms)
	return res, nil
}
