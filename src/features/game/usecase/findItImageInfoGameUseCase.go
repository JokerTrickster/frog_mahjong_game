package usecase

import (
	"context"
	_interface "main/features/game/model/interface"
	"main/features/game/model/request"
	"time"
)

type FindItImageInfoGameUseCase struct {
	Repository     _interface.IFindItImageInfoGameRepository
	ContextTimeout time.Duration
}

func NewFindItImageInfoGameUseCase(repo _interface.IFindItImageInfoGameRepository, timeout time.Duration) _interface.IFindItImageInfoGameUseCase {
	return &FindItImageInfoGameUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *FindItImageInfoGameUseCase) FindItImageInfo(c context.Context, req *request.ReqFindItImageInfo) error {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()

	//이미지 정보를 저장한다.
	for _, imageInfo := range req.ImageInfoList {
		imageDTO := CreateImageDTO(imageInfo)
		imageID, err := d.Repository.SaveImageInfo(ctx, imageDTO)
		if err != nil {
			return err
		}
		//이미지 좌표를 저장한다.
		imageCorrectDTOs := CreateImageCorrectDTO(imageID, imageInfo)
		err = d.Repository.SaveImageCorrectInfo(ctx, imageCorrectDTOs)
		if err != nil {
			return err
		}
	}

	return nil

}
