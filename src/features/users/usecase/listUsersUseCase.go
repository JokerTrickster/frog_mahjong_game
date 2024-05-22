package usecase

import (
	"context"
	_interface "main/features/users/model/interface"
	"main/features/users/model/response"
	"time"
)

type ListUsersUseCase struct {
	Repository     _interface.IListUsersRepository
	ContextTimeout time.Duration
}

func NewListUsersUseCase(repo _interface.IListUsersRepository, timeout time.Duration) _interface.IListUsersUseCase {
	return &ListUsersUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *ListUsersUseCase) List(c context.Context) (response.ResListUser, error) {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()
	users, err := d.Repository.FindUsers(ctx)
	if err != nil {
		return response.ResListUser{}, err
	}
	total, err := d.Repository.CountUsers(ctx)
	if err != nil {
		return response.ResListUser{}, err
	}
	res := CreateResListUser(users, total)

	return res, nil
}
