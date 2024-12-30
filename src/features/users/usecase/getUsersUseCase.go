package usecase

import (
	"context"
	_interface "main/features/users/model/interface"
	"main/features/users/model/response"
	"time"
)

type GetUsersUseCase struct {
	Repository     _interface.IGetUsersRepository
	ContextTimeout time.Duration
}

func NewGetUsersUseCase(repo _interface.IGetUsersRepository, timeout time.Duration) _interface.IGetUsersUseCase {
	return &GetUsersUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *GetUsersUseCase) Get(c context.Context, userID int) (response.ResGetUser, error) {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()
	userDTO, err := d.Repository.FindOneUser(ctx, userID)
	if err != nil {
		return response.ResGetUser{}, err
	}
	disconnected, err := d.Repository.CheckDisconnect(ctx, userID)
	if err != nil{
		return response.ResGetUser{}, err
	}	

	// create ResGetUser
	res := CreateResGetUser(userDTO,disconnected)

	return res, nil
}
