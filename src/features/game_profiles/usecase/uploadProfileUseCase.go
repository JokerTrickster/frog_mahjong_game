package usecase

import (
	"context"
	"fmt"
	"main/features/game_profiles/model/entity"
	_interface "main/features/game_profiles/model/interface"
	_aws "main/utils/aws"
	"time"
)

type UploadProfilesUseCase struct {
	Repository     _interface.IUploadProfilesRepository
	ContextTimeout time.Duration
}

func NewUploadProfilesUseCase(repo _interface.IUploadProfilesRepository, timeout time.Duration) _interface.IUploadProfilesUseCase {
	return &UploadProfilesUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *UploadProfilesUseCase) Upload(c context.Context, entity entity.ImageUploadProfileEntity) error {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()

	// 랜덤으로 이미지 이름 생성
	filename := fmt.Sprintf("%s.png", _aws.FileNameGenerateRandom())

	//profileDTO 생성
	profileDTO := CreateProfileDTO(entity, filename)
	// 디비에 이미지 파일 이름 저장
	err := d.Repository.InsertOneProfile(ctx, profileDTO)
	if err != nil {
		return err
	}

	// s3 이미지 파일 업로드
	err = _aws.ImageUpload(ctx, entity.Image, filename, _aws.ImgTypeProfile)
	if err != nil {
		return err
	}

	return nil
}
