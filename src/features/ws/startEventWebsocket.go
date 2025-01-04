package ws

import (
	"context"
	"fmt"
	"main/features/ws/model/entity"
	"main/features/ws/repository"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func StartEventWebsocket(msg *entity.WSMessage) {
	//유저 상태를 변경한다. (대기실로 이동)
	ctx := context.Background()
	uID := msg.UserID
	roomID := msg.RoomID

	// 비즈니스 로직
	roomInfoMsg := entity.RoomInfo{}
	preloadUsers := []entity.RoomUsers{}
	err := mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {

		// 방장이 게임 시작 요청했는지 체크
		newErr := repository.StartCheckOwner(ctx, tx, uID, roomID)
		if newErr != nil {
			SendErrorMessage(msg, newErr)
			return fmt.Errorf("%s", newErr.Msg)
		}

		roomUsers, newErr := repository.StartFindRoomUsers(ctx, tx, roomID)
		if newErr != nil {
			SendErrorMessage(msg, newErr)
			return fmt.Errorf("%s", newErr.Msg)
		}

		// room user 데이터 변경 (플레이 순번 랜덤으로 생성)
		updatedRoomUsers, newErr := StartUpdateRoomUsers(roomUsers)
		if newErr != nil {
			SendErrorMessage(msg, newErr)
			return fmt.Errorf("%s", newErr.Msg)
		}

		// room user 데이터 변경 (플레이 순번 랜덤으로 생성)
		newErr = repository.StartUpdateRoomUser(ctx, tx, updatedRoomUsers)
		if newErr != nil {
			SendErrorMessage(msg, newErr)
			return fmt.Errorf("%s", newErr.Msg)
		}

		// room 데이터 상태 변경 (대기 -> 플레이)
		newErr = repository.StartUpdateRoom(ctx, tx, roomID, "play")
		if newErr != nil {
			SendErrorMessage(msg, newErr)
			return fmt.Errorf("%s", newErr.Msg)
		}

		// 카드 정보 가져온다.
		cards, newErr := repository.StartFindCards(ctx, tx)
		if newErr != nil {
			SendErrorMessage(msg, newErr)
			return fmt.Errorf("%s", newErr.Msg)
		}

		// cards 데이터 생성
		userCards := CreateInitCards(roomID, cards)
		newErr = repository.StartCreateCards(ctx, tx, userCards)
		if newErr != nil {
			SendErrorMessage(msg, newErr)
			return fmt.Errorf("%s", newErr.Msg)
		}

		preloadUsers, newErr = repository.PreloadFindGameInfo(ctx, tx, roomID)
		if newErr != nil {
			SendErrorMessage(msg, newErr)
			return fmt.Errorf("%s", newErr.Msg)
		}

		return nil
	})
	if err != nil {
		return
	}

	// 메시지 생성
	roomInfoMsg = *CreateRoomInfoMSG(ctx, preloadUsers, 1, roomInfoMsg.ErrorInfo)

	// 구조체를 JSON 문자열로 변환 (마샬링)
	message, err := CreateMessage(&roomInfoMsg)
	if err != nil {
		fmt.Println(err)
	}
	msg.Message = message
	sendMessageToClients(roomID, msg)

}
