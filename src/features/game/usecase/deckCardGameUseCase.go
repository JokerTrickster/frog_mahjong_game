package usecase

import (
	"context"
	_interface "main/features/game/model/interface"
	"main/features/game/model/request"
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

func (d *DeckCardGameUseCase) DeckCard(c context.Context, userID int, req *request.ReqDeckCard) (response.ResDeckCardGame, error) {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()

	// 유저가 방안에 있는지 체크
	err := d.Repository.CheckRoomUser(ctx, userID, req.RoomID)
	if err != nil {
		return response.ResDeckCardGame{}, err
	}
	// 카드 정보를 생성한다.

	//응답한다.

	return response.ResDeckCardGame{}, nil

}
