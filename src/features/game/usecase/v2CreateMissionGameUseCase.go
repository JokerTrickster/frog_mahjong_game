package usecase

import (
	"context"
	"fmt"
	_interface "main/features/game/model/interface"
	"main/features/game/model/request"
	"main/utils/aws"
	"time"
)

type CreateMissionGameUseCase struct {
	Repository     _interface.ICreateMissionGameRepository
	ContextTimeout time.Duration
}

func NewCreateMissionGameUseCase(repo _interface.ICreateMissionGameRepository, timeout time.Duration) _interface.ICreateMissionGameUseCase {
	return &CreateMissionGameUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *CreateMissionGameUseCase) CreateMission(c context.Context, req *request.ReqCreateMission) error {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()

	//s3 이미지 업로드
	fileName := fmt.Sprintf("%s.png", aws.FileNameGenerateRandom())
	err := aws.ImageUpload(ctx, req.Image, fileName, aws.ImgTypeMission)
	if err != nil {
		return err
	}

	// mission DTO 생성
	missionDTO := CreateMissionDTO(req,fileName)

	// db에 저장
	err = d.Repository.SaveMission(ctx, missionDTO)
	if err != nil {
		return err
	}
	return nil

}
