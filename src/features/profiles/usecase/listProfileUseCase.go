package usecase

import (
	"context"
	_interface "main/features/profiles/model/interface"
	"main/features/profiles/model/response"
	"time"
)

type ListProfilesUseCase struct {
	Repository     _interface.IListProfilesRepository
	ContextTimeout time.Duration
}

func NewListProfilesUseCase(repo _interface.IListProfilesRepository, timeout time.Duration) _interface.IListProfilesUseCase {
	return &ListProfilesUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *ListProfilesUseCase) List(c context.Context, userID uint) (response.ResListProfile, error) {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()

	//유저 프로필 정보를 모두 가져온다.
	profileList, err := d.Repository.FindAllProfiles(ctx, userID)
	if err != nil {
		return response.ResListProfile{}, err
	}
	res := CreateResProfileList(profileList)

	return res, nil
}
