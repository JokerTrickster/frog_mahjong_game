package usecase

import (
	"context"
	"fmt"
	_interface "main/features/rooms/model/interface"
	"time"
)

type V02JoinRoomsUseCase struct {
	Repository     _interface.IV02JoinRoomsRepository
	ContextTimeout time.Duration
}

func NewV02JoinRoomsUseCase(repo _interface.IV02JoinRoomsRepository, timeout time.Duration) _interface.IV02JoinRoomsUseCase {
	return &V02JoinRoomsUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *V02JoinRoomsUseCase) V02Join(c context.Context) error {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()
	fmt.Println(ctx)

	return nil
}
