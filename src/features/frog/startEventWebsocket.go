package ws

import (
	"context"
	"fmt"
	"main/features/frog/model/entity"
	"main/features/frog/repository"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func StartEventWebsocket(msg *entity.WSMessage) *entity.ErrorInfo {
	//유저 상태를 변경한다. (대기실로 이동)
	ctx := context.Background()
	uID := msg.UserID
	roomID := msg.RoomID

	// 비즈니스 로직
	roomInfoMsg := entity.RoomInfo{}
	preloadUsers := []entity.RoomUsers{}
	var errInfo *entity.ErrorInfo
	err := mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {

		// 방장이 게임 시작 요청했는지 체크
		errInfo := repository.StartCheckOwner(ctx, tx, uID, roomID)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}

		roomUsers, errInfo := repository.StartFindRoomUsers(ctx, tx, roomID)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}

		// room user 데이터 변경 (플레이 순번 랜덤으로 생성)
		updatedRoomUsers, errInfo := StartUpdateRoomUsers(roomUsers)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}

		// room user 데이터 변경 (플레이 순번 랜덤으로 생성)
		errInfo = repository.StartUpdateRoomUser(ctx, tx, updatedRoomUsers)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}

		// room 데이터 상태 변경 (대기 -> 플레이)
		errInfo = repository.StartUpdateRoom(ctx, tx, roomID, "play")
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}

		// 카드 정보 가져온다.
		cards, errInfo := repository.StartFindCards(ctx, tx)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}

		// cards 데이터 생성
		userCards := CreateInitCards(roomID, cards)
		errInfo = repository.StartCreateCards(ctx, tx, userCards)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}

		preloadUsers, errInfo = repository.PreloadFindGameInfo(ctx, tx, roomID)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}

		return nil
	})
	if err != nil {
		return errInfo
	}

	// 메시지 생성
	roomInfoMsg = *CreateRoomInfoMSG(ctx, preloadUsers, 0, roomInfoMsg.ErrorInfo)

	// 구조체를 JSON 문자열로 변환 (마샬링)
	message, err := CreateMessage(&roomInfoMsg)
	if err != nil {
		fmt.Println(err)
	}
	msg.Message = message
	sendMessageToClients(roomID, msg)
	return nil
}
