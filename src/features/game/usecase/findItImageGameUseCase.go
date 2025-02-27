package usecase

import (
	"context"
	"main/features/game/model/entity"
	_interface "main/features/game/model/interface"
	_aws "main/utils/aws"
	"time"
)

type FindItImageGameUseCase struct {
	Repository     _interface.IFindItImageGameRepository
	ContextTimeout time.Duration
}

func NewFindItImageGameUseCase(repo _interface.IFindItImageGameRepository, timeout time.Duration) _interface.IFindItImageGameUseCase {
	return &FindItImageGameUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *FindItImageGameUseCase) FindItImage(c context.Context, e *entity.FindItImageGameEntity) error {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()
	findItName := e.Image.Filename

	// s3 이미지 파일 업로드
	err := _aws.ImageUpload(ctx, e.Image, findItName, _aws.ImgTypeFindIt)
	if err != nil {
		return err
	}
	return nil

}
