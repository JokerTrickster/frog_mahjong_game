package usecase

import (
	"context"
	_interface "main/features/game_users/model/interface"
	"main/features/game_users/model/response"
	"time"
)

type ListProfilesUsersUseCase struct {
	Repository     _interface.IListProfilesUsersRepository
	ContextTimeout time.Duration
}

func NewListProfilesUsersUseCase(repo _interface.IListProfilesUsersRepository, timeout time.Duration) _interface.IListProfilesUsersUseCase {
	return &ListProfilesUsersUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *ListProfilesUsersUseCase) ListProfiles(c context.Context, userID uint) (response.ResListProfileGameUser, error) {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()

	//유저 프로필 정보를 모두 가져온다.
	profileList, err := d.Repository.FindAllProfiles(ctx, userID)
	if err != nil {
		return response.ResListProfileGameUser{}, err
	}
	res := CreateResProfileList(profileList)

	return res, nil
}
