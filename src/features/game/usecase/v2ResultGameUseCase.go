package usecase

import (
	"context"
	_interface "main/features/game/model/interface"
	"main/features/game/model/request"
	"main/features/game/model/response"
	"time"
)

type V2ResultGameUseCase struct {
	Repository     _interface.IV2ResultGameRepository
	ContextTimeout time.Duration
}

func NewV2ResultGameUseCase(repo _interface.IV2ResultGameRepository, timeout time.Duration) _interface.IV2ResultGameUseCase {
	return &V2ResultGameUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *V2ResultGameUseCase) V2Result(c context.Context, req *request.ReqV2Result) (response.ResV2Result, error) {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()

	// 승리한 유저에 미션 성공 카드들을 모두 가져온다.
	// user_missions 에서 미션 정보를 모두 가져온다.
	userMissions, err := d.Repository.GetUserMissions(ctx, req)
	if err != nil {
		return response.ResV2Result{}, err
	}
	res := response.ResV2Result{
		Winner: req.UserID,
	}
	for _, userMission := range userMissions {
		// 미션마다 달성한 카드 정보를 모두 가져온다,
		userMissionCards, err := d.Repository.GetUserMissionCards(ctx, userMission.ID)
		if err != nil {
			return response.ResV2Result{}, err
		}
		cards := []int{}
		for _, userMissionCard := range userMissionCards {
			cards = append(cards, userMissionCard.CardID)
		}

		mission := response.ResultMission{
			MissionID: uint(userMission.MissionID),
			Cards:     cards,
		}

		res.Missions = append(res.Missions, mission)
	}

	return res, nil
}
