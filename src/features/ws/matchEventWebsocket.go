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

func MatchEventWebsocket(msg *entity.WSMessage) {
	ctx := context.Background()
	uID := msg.UserID

	//string to struct
	req := request.ReqWSMatchEvent{}
	err := json.Unmarshal([]byte(msg.Message), &req)
	if err != nil {
		log.Fatalf("JSON 언마샬링 에러: %s", err)
	}

	//비즈니스 로직
	roomInfoMsg := entity.RoomInfo{}
	preloadUsers := []entity.RoomUsers{}
	RoomDTO, _ := repository.MatchFindOneWaitingRoom(ctx, uint(req.Count), uint(req.Timer))
	roomID := RoomDTO.ID
	fmt.Println(RoomDTO)
	err = mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		// 방 참여 가능한지 체크
		if RoomDTO.CurrentCount == RoomDTO.MaxCount {
			return fmt.Errorf("방이 꽉 찼습니다.")
		}
		if RoomDTO.State != "wait" {
			return fmt.Errorf("게임 중인 방입니다.")
		}
		// TODO 기존에 방 유저 정보가 있는지 가져온다.
		// 유저 정보가 있으면 삭제하고 방 인원수를 감소시킨다.
		err = repository.MatchFindOneAndDeleteRoomUser(ctx, tx, uID, roomID)
		if err != nil {
			return err
		}
		// 방 유저 정보를 생성한다.
		RoomUserDTO := CreateMatchRoomUserDTO(uID, int(roomID), "wait")
		if err != nil {
			return err
		}

		err = repository.MatchInsertOneRoomUser(ctx, tx, RoomUserDTO)
		if err != nil {
			return err
		}
		// 방 현재 인원을 증가시킨다.
		err = repository.MatchFindOneAndUpdateRoom(ctx, tx, roomID)
		if err != nil {
			return err
		}

		//유저 정보를 업데이트 한다.
		err = repository.MatchFindOneAndUpdateUser(ctx, tx, uID, roomID)
		if err != nil {
			return err
		}

		preloadUsers, err = repository.MatchFindAllRoomUsers(ctx, tx, roomID)
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
		if roomInfoMsg.ErrorInfo.Msg == "방이 꽉 찼습니다." {
			roomInfoMsg.ErrorInfo.Type = _errors.ErrRoomFull
		} else if roomInfoMsg.ErrorInfo.Msg == "비밀번호가 일치하지 않습니다." {
			roomInfoMsg.ErrorInfo.Type = _errors.ErrWrongPassword
		} else if roomInfoMsg.ErrorInfo.Msg == "게임 중인 방입니다." {
			roomInfoMsg.ErrorInfo.Type = _errors.ErrGameInProgress
		}
	}

	// 메시지 생성
	roomInfoMsg = *CreateRoomInfoMSG(ctx, preloadUsers, 1, roomInfoMsg.ErrorInfo)
	roomInfoMsg.GameInfo.AllReady = false
	// 구조체를 JSON 문자열로 변환 (마샬링)
	message, err := CreateMessage(&roomInfoMsg)
	if err != nil {
		fmt.Println(err)
	}
	msg.Message = message
	msg.RoomID = roomID

	//방 유저들에게 메시지 전달
	if clients, ok := entity.WSClients[msg.RoomID]; ok {
		//에러 발생시 이벤트 요청한 유저에게만 메시지를 전달한다.
		if roomInfoMsg.ErrorInfo != nil || err != nil {
			for client := range clients {
				if clients[client].UserID == msg.UserID {
					_ = client.WriteJSON(msg)
					clientData := clients[client]
					clientData.Close()
					clients[client] = clientData
					delete(clients, client)
				}
			}
		} else {
			for client := range clients {
				err := client.WriteJSON(msg)
				if err != nil {
					fmt.Printf("error: %v", err)
					client.Close()
					delete(clients, client)
				}
			}
		}
	}
}
