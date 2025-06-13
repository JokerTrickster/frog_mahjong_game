package _interface

import (
	"context"
	"main/features/game_profiles/model/entity"
	"main/features/game_profiles/model/response"
)

type IListProfilesUseCase interface {
	List(c context.Context) (response.ResListGameProfile, error)
}

type IUploadProfilesUseCase interface {
	Upload(c context.Context, e entity.ImageUploadProfileEntity) error
}

type IUpdateProfilesUseCase interface {
	Update(c context.Context, userID int, profileID int) (response.ResUpdateProfile, error)
}
