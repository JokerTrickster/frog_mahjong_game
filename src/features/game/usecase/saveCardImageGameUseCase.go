package usecase

import (
	"context"
	"main/features/game/model/entity"
	_interface "main/features/game/model/interface"
	_aws "main/utils/aws"
	"time"
)

type SaveCardImageGameUseCase struct {
	Repository     _interface.ISaveCardImageGameRepository
	ContextTimeout time.Duration
}

func NewSaveCardImageGameUseCase(repo _interface.ISaveCardImageGameRepository, timeout time.Duration) _interface.ISaveCardImageGameUseCase {
	return &SaveCardImageGameUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *SaveCardImageGameUseCase) SaveCardImage(c context.Context, e entity.SaveCardImageGameEntity) error {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()
	// s3 이미지 파일 업로드
	err := _aws.ImageUpload(ctx, e.Image, e.Image.Filename, _aws.ImgTypeBirdCard)
	if err != nil {
		return err
	}
	return nil

}
