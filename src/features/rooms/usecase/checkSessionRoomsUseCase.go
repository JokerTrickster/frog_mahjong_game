package usecase

import (
	"context"
	_interface "main/features/rooms/model/interface"
	"main/features/rooms/model/request"
	"time"
)

type CheckSessionRoomsUseCase struct {
	Repository     _interface.ICheckSessionRoomsRepository
	ContextTimeout time.Duration
}

func NewCheckSessionRoomsUseCase(repo _interface.ICheckSessionRoomsRepository, timeout time.Duration) _interface.ICheckSessionRoomsUseCase {
	return &CheckSessionRoomsUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *CheckSessionRoomsUseCase) CheckSession(c context.Context, req *request.ReqCheckSession) (bool, error) {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()

	existed, err := d.Repository.RedisCheckSession(ctx, req)
	if err != nil {
		return false, err
	}

	return existed, nil
}
