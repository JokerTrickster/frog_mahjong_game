package usecase

import (
	"context"
	_interface "main/features/rooms/model/interface"
	"main/features/rooms/model/response"
	"time"
)

type MetaRoomsUseCase struct {
	Repository     _interface.IMetaRoomsRepository
	ContextTimeout time.Duration
}

func NewMetaRoomsUseCase(repo _interface.IMetaRoomsRepository, timeout time.Duration) _interface.IMetaRoomsUseCase {
	return &MetaRoomsUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *MetaRoomsUseCase) Meta(c context.Context) (response.ResMetaRoom, error) {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()
	timeDTO, err := d.Repository.FindAllTimeMeta(ctx)
	if err != nil {
		return response.ResMetaRoom{}, err
	}

	res := CreateResMetaData(timeDTO)

	return res, nil
}
