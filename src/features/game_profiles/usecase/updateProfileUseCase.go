package usecase

import (
	"context"
	_interface "main/features/game_profiles/model/interface"
	"main/features/game_profiles/model/response"
	"time"
)

type UpdateProfilesUseCase struct {
	Repository     _interface.IUpdateProfilesRepository
	ContextTimeout time.Duration
}

func NewUpdateProfilesUseCase(repo _interface.IUpdateProfilesRepository, timeout time.Duration) _interface.IUpdateProfilesUseCase {
	return &UpdateProfilesUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *UpdateProfilesUseCase) Update(c context.Context, userID int, profileID int) (response.ResUpdateProfile, error) {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()

	err := d.Repository.UpdateOneProfile(ctx, userID, profileID)
	if err != nil {
		return response.ResUpdateProfile{}, err
	}

	return response.ResUpdateProfile{
		ProfileID: profileID,
		UserID:    userID,
	}, nil
}
