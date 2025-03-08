package usecase

import (
	"context"
	_interface "main/features/board_game/model/interface"
	"main/features/board_game/model/request"
	"main/features/board_game/model/response"
	"time"
)

type FindItSoloPlayBoardGameUseCase struct {
	Repository     _interface.IFindItSoloPlayBoardGameRepository
	ContextTimeout time.Duration
}

func NewFindItSoloPlayBoardGameUseCase(repo _interface.IFindItSoloPlayBoardGameRepository, timeout time.Duration) _interface.IFindItSoloPlayBoardGameUseCase {
	return &FindItSoloPlayBoardGameUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *FindItSoloPlayBoardGameUseCase) FindItSoloPlay(c context.Context, userID int, req *request.ReqFindItSoloPlayBoardGame) (response.ResFindItSoloPlayBoardGame,error) {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()

	// 랜덤으로 이미지를 가져온다.
	ImageDTOList, err := d.Repository.FindRandomImage(ctx, req.Round)
	if err != nil {
		return response.ResFindItSoloPlayBoardGame{}, err
	}
	res := response.ResFindItSoloPlayBoardGame{}
	gameInfoList := []response.SoloPlayGameInfo{}
	// 해당 이미지에 정답 좌표들을 가져온다.
	round := 1 
	for _, image := range ImageDTOList {
		correctDTOList, err := d.Repository.FindCorrectByImageID(ctx, image.ID)
		if err != nil {
			return response.ResFindItSoloPlayBoardGame{}, err
		}
		gameInfo := CreateSoloPlayGameInfo(image, correctDTOList,round )
		gameInfoList = append(gameInfoList, gameInfo)
		round++
	}
	res.GameInfoList = gameInfoList

	return res,nil
}
