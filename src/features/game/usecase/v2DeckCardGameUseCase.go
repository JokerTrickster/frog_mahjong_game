package usecase

import (
	"context"
	_interface "main/features/game/model/interface"
	"main/features/game/model/response"
	"time"
)

type V2DeckCardGameUseCase struct {
	Repository     _interface.IV2DeckCardGameRepository
	ContextTimeout time.Duration
}

func NewV2DeckCardGameUseCase(repo _interface.IV2DeckCardGameRepository, timeout time.Duration) _interface.IV2DeckCardGameUseCase {
	return &V2DeckCardGameUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *V2DeckCardGameUseCase) V2DeckCard(c context.Context, userID, roomID int) (response.ResV2DeckCardGame, error) {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()

	// 유저가 방안에 있는지 체크
	err := d.Repository.CheckRoomUser(ctx, userID, roomID)
	if err != nil {
		return response.ResV2DeckCardGame{}, err
	}
	
	res := CreateV2RandomCardIDList()
	

	return res, nil

}
