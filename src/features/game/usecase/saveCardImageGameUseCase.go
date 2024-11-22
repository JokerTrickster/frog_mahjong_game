package usecase

import (
	"context"
	"fmt"
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
	cardName := e.Image.Filename
	//file 이름 랜덤으로 변경 + 확장자
	newFileName := fmt.Sprintf("%s.png", _aws.FileNameGenerateRandom())
	//디비에 image 정보 수정
	err := d.Repository.FindOneUpdateCardImage(ctx, cardName, newFileName)
	if err != nil {
		return err
	}
	// s3 이미지 파일 업로드
	err = _aws.ImageUpload(ctx, e.Image, newFileName, _aws.ImgTypeBirdCard)
	if err != nil {
		return err
	}
	return nil

}
