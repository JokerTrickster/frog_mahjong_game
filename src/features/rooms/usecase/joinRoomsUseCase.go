package usecase

import (
	"context"
	_errors "main/features/rooms/model/errors"
	_interface "main/features/rooms/model/interface"
	"main/features/rooms/model/request"
	"main/features/rooms/model/response"
	"main/utils"
	"main/utils/db/mysql"
	"time"

	"gorm.io/gorm"
)

type JoinRoomsUseCase struct {
	Repository     _interface.IJoinRoomsRepository
	ContextTimeout time.Duration
}

func NewJoinRoomsUseCase(repo _interface.IJoinRoomsRepository, timeout time.Duration) _interface.IJoinRoomsUseCase {
	return &JoinRoomsUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *JoinRoomsUseCase) Join(c context.Context, uID uint, email string, req *request.ReqJoin) (response.ResJoinRoom, error) {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()

	res := response.ResJoinRoom{}
	err := mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		// 방 참여 가능한지 체크
		RoomDTO, err := d.Repository.FindOneRoom(ctx,tx, req)
		if err != nil {
			return err
		}
		if RoomDTO.CurrentCount == RoomDTO.MaxCount {
			return utils.ErrorMsg(ctx, utils.ErrRoomFull, utils.Trace(), _errors.ErrRoomFull.Error(), utils.ErrFromClient)
		}
		// 방 유저 정보를 생성한다.
		RoomUserDTO, err := CreateRoomUserDTO(uID, int(req.RoomID), "wait")
		if err != nil {
			return err
		}
		err = d.Repository.InsertOneRoomUser(ctx,tx, RoomUserDTO)
		if err != nil {
			return err
		}
		// 방 현재 인원을 증가시킨다.
		err = d.Repository.FindOneAndUpdateRoom(ctx,tx, req.RoomID)
		if err != nil {
			return err
		}

		//유저 정보를 업데이트 한다.
		err = d.Repository.FindOneAndUpdateUser(ctx,tx, uID, req.RoomID)
		if err != nil {
			return err
		}
		res = response.ResJoinRoom{
			RoomID: int(req.RoomID),
		}
		return nil
	})
	return res, err
}
