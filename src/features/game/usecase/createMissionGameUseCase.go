package usecase

import (
	"context"
	_interface "main/features/game/model/interface"
	"main/features/game/model/request"
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

	// mission DTO 생성
	missionDTO := CreateMissionDTO(req)

	// db에 저장
	err := d.Repository.SaveMission(ctx, missionDTO)
	if err != nil {
		return err
	}
	return nil

}
