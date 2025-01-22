package usecase

import (
	"context"
	_interface "main/features/game/model/interface"
	"main/features/game/model/response"
	"time"
)

type V2DrawResultGameUseCase struct {
	Repository     _interface.IV2DrawResultGameRepository
	ContextTimeout time.Duration
}

func NewV2DrawResultGameUseCase(repo _interface.IV2DrawResultGameRepository, timeout time.Duration) _interface.IV2DrawResultGameUseCase {
	return &V2DrawResultGameUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *V2DrawResultGameUseCase) V2DrawResult(c context.Context, roomID int) (response.ResV2DrawResult, error) {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()

	roomUsers, err := d.Repository.FindAllRoomUsers(ctx, roomID)
	if err != nil {
		return response.ResV2DrawResult{}, err
	}
	res := response.ResV2DrawResult{}
	for _, roomUser := range roomUsers {
		userMission, err := d.Repository.FindAllUserMission(ctx, roomUser.UserID, roomID)
		if err != nil {
			return response.ResV2DrawResult{}, err
		}
		drawResult := response.DrawResult{
			UserID: roomUser.UserID,
		}
		missionList := make([]int, 0, len(userMission))
		for _, mission := range userMission {
			missionList = append(missionList, mission.MissionID)
		}
		drawResult.SuccessMissions = missionList
		res.Users = append(res.Users, drawResult)
	}

	return res, nil
}
