package _interface

import (
	"context"
	"main/features/profiles/model/entity"
	"main/features/profiles/model/response"
)

type IListProfilesUseCase interface {
	List(c context.Context) (response.ResListProfile, error)
}

type IUploadProfilesUseCase interface {
	Upload(c context.Context, e entity.ImageUploadProfileEntity) error
}
