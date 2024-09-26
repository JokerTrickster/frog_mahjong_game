package usecase

import (
	"context"
	_interface "main/features/game/model/interface"
	"main/features/game/model/response"
	"time"
)

type DeckCardGameUseCase struct {
	Repository     _interface.IDeckCardGameRepository
	ContextTimeout time.Duration
}

func NewDeckCardGameUseCase(repo _interface.IDeckCardGameRepository, timeout time.Duration) _interface.IDeckCardGameUseCase {
	return &DeckCardGameUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *DeckCardGameUseCase) DeckCard(c context.Context, userID, roomID int) (response.ResDeckCardGame, error) {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()

	// 유저가 방안에 있는지 체크
	err := d.Repository.CheckRoomUser(ctx, userID, roomID)
	if err != nil {
		return response.ResDeckCardGame{}, err
	}
	
	// 랜덤으로 카드 생성 (1부터 44까지 랜덤으로 생성해서 배열에 저장)
	res := CreateRandomCardIDList()
	

	return res, nil

}
