package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"main/features/ws/model/entity"
	_errors "main/features/ws/model/errors"
	"main/features/ws/model/request"
	"main/features/ws/repository"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func FailedLoanEventWebsocket(msg *entity.WSMessage) {
	ctx := context.Background()
	uID := msg.UserID
	roomID := msg.RoomID

	//string to struct
	req := request.ReqWSFailedLoan{}
	err := json.Unmarshal([]byte(msg.Message), &req)
	if err != nil {
		log.Printf("JSON 언마샬링 에러: %s", err)
	}
	loanEntity := entity.WSLoanEntity{
		RoomID:       roomID,
		CardID:       req.CardID,
		TargetUserID: req.TargetUserID,
		UserID:       uID,
	}

	// 비즈니스 로직
	roomInfoMsg := entity.RoomInfo{}
	doraDTO := &mysql.FrogUserCards{}
	preloadUsers := []entity.RoomUsers{}
	err = mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		// 소유하고 있는 카드인지 체크
		err := repository.FailedLoanCheckCard(ctx, tx, &loanEntity)
		if err != nil {
			return err
		}

		// 카드 정보를 롤백한다.
		err = repository.FailedLoanRollbackCard(ctx, tx, &loanEntity)
		if err != nil {
			return err
		}

		// // 패널티를 부여한다. (코인 차감)
		// penaltyCoin := (len(entity.WSClients[msg.RoomID]) - 1) * 2
		// err = repository.FailedLoanPenalty(ctx, tx, &loanEntity, penaltyCoin)
		// if err != nil {
		// 	return err
		// }

		// // 모든 플레이어에게 코인 추가
		// err = repository.FailedLoanAddCoin(ctx, tx, &loanEntity)
		// if err != nil {
		// 	return err
		// }

		//dora 카드 가져오기
		doraDTO, err = repository.LoanCardFindOneDora(ctx, tx, roomID)
		if err != nil {
			return err
		}

		preloadUsers, err = repository.FailedLoanFindAllRoomUsers(ctx, tx, roomID)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		roomInfoMsg.ErrorInfo = &entity.ErrorInfo{
			Code: 500,
			Msg:  err.Error(),
			Type: _errors.ErrInternalServer,
		}
	}

	//유저들에게 메시지 전송한다.
	if clients, ok := entity.WSClients[msg.RoomID]; ok {
		// 메시지 생성
		roomInfoMsg = *CreateRoomInfoMSG(ctx, preloadUsers, req.PlayTurn, roomInfoMsg.ErrorInfo)

		// 론 가능 여부를 true로 변경
		roomInfoMsg.GameInfo.IsLoanAllowed = true

		// 론 실패한 유저ID 저장
		roomInfoMsg.GameInfo.FailedLoanUserID = uID

		//도라 카드 정보 저장
		doraCardInfo := entity.Card{}
		doraCardInfo.CardID = uint(doraDTO.CardID)
		roomInfoMsg.GameInfo.Dora = &doraCardInfo

		//에러 발생시 이벤트 요청한 유저에게만 메시지를 전달한다.
		if roomInfoMsg.ErrorInfo != nil || err != nil {
			for client := range clients {
				if clients[client].UserID == msg.UserID {
					// 구조체를 JSON 문자열로 변환 (마샬링)
					message, err := CreateMessage(&roomInfoMsg)
					if err != nil {
						fmt.Println(err)
					}
					msg.Message = message
					err = client.WriteJSON(msg)
					if err != nil {
						log.Printf("error: %v", err)
						client.Close()
						delete(clients, client)
					}
				}
			}
		} else {
			for client := range clients {
				filterRoomInfoMsg := Deepcopy(roomInfoMsg)

				// 구조체를 JSON 문자열로 변환 (마샬링)
				message, err := CreateMessage(&filterRoomInfoMsg)
				if err != nil {
					fmt.Println(err)
				}
				msg.Message = message
				err = client.WriteJSON(msg)
				if err != nil {
					log.Printf("error: %v", err)
					client.Close()
					delete(clients, client)
				}
			}
		}
	}
}
