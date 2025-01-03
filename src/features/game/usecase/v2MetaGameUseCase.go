package usecase

import (
	"context"
	_interface "main/features/game/model/interface"
	"main/features/game/model/response"
	"time"
)

type MetaGameUseCase struct {
	Repository     _interface.IMetaGameRepository
	ContextTimeout time.Duration
}

func NewMetaGameUseCase(repo _interface.IMetaGameRepository, timeout time.Duration) _interface.IMetaGameUseCase {
	return &MetaGameUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *MetaGameUseCase) Meta(c context.Context) (response.ResMetaGame, error) {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()

	categoryList, err := d.Repository.FindAllReportCategory(ctx)
	if err != nil {
		return response.ResMetaGame{}, err
	}
	res := CreateResMetaGame(categoryList)

	return res, nil
}
