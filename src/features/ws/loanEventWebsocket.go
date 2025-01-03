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

func LoanEventWebsocket(msg *entity.WSMessage) {
	ctx := context.Background()
	uID := msg.UserID
	roomID := msg.RoomID

	//string to struct
	req := request.ReqWSLoan{}
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
		// loan 가능한지 체크 (마지막으로 버려진 카드인지 체크)
		err := repository.LoanCheckLoan(ctx, tx, &loanEntity)
		if err != nil {
			return err
		}

		// loan 하기 (상대방이 버린 카드를 가져온다)
		err = repository.LoanCardLoan(ctx, tx, &loanEntity)
		if err != nil {
			return err
		}

		// 룸 유저 카드 수와 상태값을 변경한다.
		err = repository.LoanUpdateRoomUserCardCount(ctx, tx, &loanEntity)
		if err != nil {
			return err
		}
		//dora 카드 가져오기
		doraDTO, err = repository.LoanCardFindOneDora(ctx, tx, roomID)
		if err != nil {
			return err
		}
		preloadUsers, err = repository.LoanFindAllRoomUsers(ctx, tx, roomID)
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

		//론한 유저에 대한 정보를 게임정보에 저장한다.
		LoanInfo := entity.LoanInfo{
			CardID:       int(req.CardID),
			UserID:       uID,
			TargetUserID: req.TargetUserID,
		}
		roomInfoMsg.GameInfo.LoanInfo = &LoanInfo

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
						fmt.Printf("error: %v", err)
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
					fmt.Printf("error: %v", err)
					client.Close()
					delete(clients, client)
				}
			}
		}
	}
}
