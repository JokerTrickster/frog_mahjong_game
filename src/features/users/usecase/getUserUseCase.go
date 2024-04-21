package usecase

import (
	"context"
	"fmt"
	_interface "main/features/users/model/interface"
	"time"
)

type GetUsersUseCase struct {
	Repository     _interface.IGetUsersRepository
	ContextTimeout time.Duration
}

func NewGetUsersUseCase(repo _interface.IGetUsersRepository, timeout time.Duration) _interface.IGetUsersUseCase {
	return &GetUsersUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *GetUsersUseCase) Get(c context.Context) error {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()
	fmt.Println(ctx)
	return nil
}
