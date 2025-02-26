package usecase

import (
	"context"
	_interface "main/features/game/model/interface"
	"main/features/game/model/response"
	"time"
)

type ListGameUseCase struct {
	Repository     _interface.IListGameRepository
	ContextTimeout time.Duration
}

func NewListGameUseCase(repo _interface.IListGameRepository, timeout time.Duration) _interface.IListGameUseCase {
	return &ListGameUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *ListGameUseCase) ListGame(c context.Context) (response.ResListGame, error) {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()

	// 게임 정보를 모두 가져온다. 
	gameList, err := d.Repository.FindGame(ctx)
	if err != nil {
		return response.ResListGame{}, err
	}
	// 응답을 만든다.
	res := CreateResListGame(gameList)
	return res, nil

}
