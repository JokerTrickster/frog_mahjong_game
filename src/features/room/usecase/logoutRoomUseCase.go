package usecase

import (
	"context"
	"fmt"
	_interface "main/features/room/model/interface"
	"time"
)

type LogoutRoomUseCase struct {
	Repository        _interface.ILogoutRoomRepository
	ContextTimeLogout time.Duration
}

func NewLogoutRoomUseCase(repo _interface.ILogoutRoomRepository, timeLogout time.Duration) _interface.ILogoutRoomUseCase {
	return &LogoutRoomUseCase{Repository: repo, ContextTimeLogout: timeLogout}
}

func (d *LogoutRoomUseCase) Logout(c context.Context, uID uint) error {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeLogout)
	defer cancel()
	fmt.Println(ctx)
	return nil
}
