package usecase

import (
	"context"
	_interface "main/features/game/model/interface"
	"main/features/game/model/response"
	"time"
)

type ListMissionGameUseCase struct {
	Repository     _interface.IListMissionGameRepository
	ContextTimeout time.Duration
}

func NewListMissionGameUseCase(repo _interface.IListMissionGameRepository, timeout time.Duration) _interface.IListMissionGameUseCase {
	return &ListMissionGameUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *ListMissionGameUseCase) ListMission(c context.Context) (response.ResListMissionGame, error) {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()

	//미션 리스트를 모두 가져온다
	missionDTOList, err := d.Repository.FindAllMission(ctx)
	if err != nil {
		return response.ResListMissionGame{}, err
	}

	res := CreateResListMission(missionDTOList)
	return res, nil

}
