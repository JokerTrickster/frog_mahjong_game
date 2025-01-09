package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"main/features/ws/model/entity"
	"main/features/ws/model/request"
	"main/features/ws/repository"
	"main/utils"
	"main/utils/db/mysql"
	"time"

	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

func CreateChatDTO(req request.ReqWSChat) *mysql.Chats {
	chatDTO := mysql.Chats{
		UserID:  int(req.UserID),
		RoomID:  int(req.RoomID),
		Name:    req.Name,
		Message: req.Message,
	}
	return &chatDTO
}

func Deepcopy(src entity.RoomInfo) entity.RoomInfo {
	var dst entity.RoomInfo
	b, _ := json.Marshal(src)
	json.Unmarshal(b, &dst)
	return dst
}

func FilterOwnCards(roomInfoMsg *entity.RoomInfo, userID uint) *entity.RoomInfo {
	for i := 0; i < len(roomInfoMsg.Users); i++ {
		if roomInfoMsg.Users[i].ID != userID {
			roomInfoMsg.Users[i].Cards = nil
		}
	}
	return roomInfoMsg
}

func CalcPlayTurn(playTurn, playerCount int) int {
	return (playTurn % playerCount) + 1
}

func CreateRoomInfoMSG(ctx context.Context, preloadUsers []entity.RoomUsers, playTurn int, roomInfoError *entity.ErrorInfo) *entity.RoomInfo {
	roomInfoMsg := entity.RoomInfo{}
	allReady := true
	timer := 30
	roomID := 0
	password := ""
	var startTime int64
	//유저 정보 저장
	for _, roomUser := range preloadUsers {
		timer = roomUser.Room.Timer
		password = roomUser.Room.Password
		user := entity.User{
			ID:          uint(roomUser.UserID),
			PlayerState: roomUser.PlayerState,
			Coin:        roomUser.User.Coin,
			Name:        roomUser.User.Name,
			Email:       roomUser.User.Email,
			TurnNumber:  roomUser.TurnNumber,
			ProfileID:   roomUser.User.ProfileID,
		}
		for _, card := range roomUser.Cards {
			if card.State == "owned" {
				user.Cards = append(user.Cards, &entity.Card{
					CardID: uint(card.CardID),
					UserID: uint(card.UserID),
				})
			} else if card.State == "discard" {
				user.DiscardedCards = append(user.DiscardedCards, &entity.Card{
					CardID: uint(card.CardID),
					UserID: uint(card.UserID),
				})
			}
		}
		roomID = roomUser.RoomID
		//시작 시간 추가
		if !roomUser.Room.StartTime.IsZero() {
			// 시작 시간을 epoch time milliseconds로 변환 +3초 추가
			startTime = roomUser.Room.StartTime.UnixNano()/int64(time.Millisecond) + 5000
		}

		if roomUser.Room.OwnerID == roomUser.UserID {
			user.IsOwner = true
		}
		roomInfoMsg.Users = append(roomInfoMsg.Users, &user)
	}

	//게임 정보 저장
	gameInfo := entity.GameInfo{
		PlayTurn:      playTurn,
		AllReady:      allReady,
		IsLoanAllowed: false,
		Timer:         timer,
		IsFull:        false,
		RoomID:        uint(roomID),
		Password:      password,
		StartTime:     startTime,
	}
	// 도라 정보를 가져온다.
	dora, _ := repository.FindOneDoraCard(ctx, roomID)
	if dora != nil {
		gameInfo.Dora = &entity.Card{
			CardID: uint(dora.CardID),
		}
	}
	roomInfoMsg.GameInfo = &gameInfo
	if roomInfoError != nil {
		roomInfoMsg.ErrorInfo = roomInfoError
	}
	return &roomInfoMsg

}
func CreateChatMessage(chatInfoMsg *entity.ChatInfo) (string, error) {
	// 구조체를 JSON 문자열로 변환 (마샬링)
	jsonData, err := json.Marshal(chatInfoMsg)
	if err != nil {
		return "", fmt.Errorf("JSON 마샬링 에러: %s", err)
	}

	return string(jsonData), nil
}
func CreateMessage(roomInfoMsg *entity.RoomInfo) (string, error) {
	// 구조체를 JSON 문자열로 변환 (마샬링)
	jsonData, err := json.Marshal(roomInfoMsg)
	if err != nil {
		return "", fmt.Errorf("JSON 마샬링 에러: %s", err)
	}

	return string(jsonData), nil
}

func CalcScore(cards []*mysql.FrogUserCards, score int) *entity.ErrorInfo {
	if score >= 5 {
		return nil
	}
	return CreateErrorMessage(400, "BadRequest", "점수가 부족합니다.")
}

func CreateErrorMessage(errCode int, errType, errMsg string) *entity.ErrorInfo {
	result := &entity.ErrorInfo{
		Code: errCode,
		Type: errType,
		Msg:  errMsg,
	}
	return result
}

// 클라이언트에 메시지 전송
func sendMessageToClients(roomID uint, msg *entity.WSMessage) {
	// 로그 메시지 생성
	utils.LogError(msg.Message)

	// // 메시지 암호화
	// encryptedMessage, err := utils.EncryptAES(msg.Message)
	// if err != nil {
	// 	fmt.Printf("Failed to encrypt message: %v\n", err)
	// 	return
	// }
	// msg.Message = encryptedMessage

	// 방에 있는 모든 클라이언트에 메시지 전송
	if sessionIDs, ok := entity.RoomSessions[roomID]; ok {
		for _, sessionID := range sessionIDs {
			if client, exists := entity.WSClients[sessionID]; exists {
				if err := client.Conn.WriteJSON(msg); err != nil {
					client.Closed = true
				}
			}
		}
	}
}

// 특정 크라이언트에 메시지 전송
func sendMessageToClient(roomID uint, msg *entity.WSMessage) {
	// 로그 메시지 생성
	utils.LogError(msg.Message)

	// // 메시지 암호화
	// encryptedMessage, err := utils.EncryptAES(msg.Message)
	// if err != nil {
	// 	fmt.Printf("Failed to encrypt message: %v\n", err)
	// 	return
	// }
	// msg.Message = encryptedMessage

	// 방에 있는 모든 클라이언트에 메시지 전송
	if sessionIDs, ok := entity.RoomSessions[roomID]; ok {
		for _, sessionID := range sessionIDs {
			if client, exists := entity.WSClients[sessionID]; exists && client.UserID == msg.UserID {
				if err := client.Conn.WriteJSON(msg); err != nil {
					client.Close()
					delete(entity.WSClients, sessionID)
					removeSessionFromRoom(roomID, sessionID)
				}
			}
		}
	}
}

func SendWebSocketCloseMessage(ws *websocket.Conn, closeCode int, message string) error {
	utils.LogError(message)
	closeMessage := websocket.FormatCloseMessage(closeCode, message)
	err := ws.WriteMessage(websocket.CloseMessage, closeMessage)
	return err
}

// 기존 연결 복구
func restoreSession(ws *websocket.Conn, sessionID string, roomID uint, userID uint) {
	// 타이머 취소
	if timer, ok := reconnectTimers.Load(sessionID); ok {
		timer.(*time.Timer).Stop()
		reconnectTimers.Delete(sessionID)
		fmt.Printf("Reconnection successful for session %s in room %d. Timer canceled.\n", sessionID, roomID)
	}
	// 세션 ID 생성
	newSessionID := generateSessionID()

	// 세션 ID 저장
	errInfo := repository.RedisSessionSet(context.TODO(), newSessionID, roomID)
	if errInfo != nil {
		fmt.Println(errInfo)
	}
	// 새로운 세션으로 등록
	registerNewSession(ws, newSessionID, roomID, userID)
}

// 새로운 세션 등록
func registerNewSession(ws *websocket.Conn, sessionID string, roomID uint, userID uint) {
	// 세션 등록
	wsClient := &entity.WSClient{
		SessionID: sessionID,
		RoomID:    roomID,
		UserID:    userID,
		Conn:      ws,
		Closed:    false,
	}
	entity.WSClients[sessionID] = wsClient
	fmt.Println(entity.WSClients[sessionID], sessionID)

	// 방에 세션 추가
	entity.RoomSessions[roomID] = append(entity.RoomSessions[roomID], sessionID)
	fmt.Println("룸 세션 수 ", len(entity.RoomSessions[roomID]))
	// 핑/퐁 핸들링 시작
	go HandlePingPong(wsClient)

	// 메시지 처리 루프 시작
	go readMessages(ws, sessionID, roomID, userID)

}

// 메시지 읽기 및 처리
func readMessages(ws *websocket.Conn, sessionID string, roomID uint, userID uint) {
	client := entity.WSClients[sessionID]
	defer func() {
		// 연결 종료 시 세션 정리
		client.Closed = true
		fmt.Println("Session", sessionID, "closed. Read loop stopped.")
	}()

	for {
		if client.Closed {
			log.Printf("Session %s is closed. Stopping read loop.", sessionID)
			return
		}

		var msg entity.WSMessage
		err := ws.ReadJSON(&msg)
		fmt.Println(msg)
		if err != nil {
			if closeErr, ok := err.(*websocket.CloseError); ok {
				if closeErr.Code == websocket.CloseNormalClosure {
					log.Printf("Session %s closed normally (Code 1000).", sessionID)
					break
				}
				log.Printf("Session %s closed with error: %v", sessionID, closeErr)
			} else {
				log.Printf("Error reading message for session %s: %v", sessionID, err)
			}
			break
		}

		// 메시지 브로드캐스트
		msg.RoomID = roomID
		msg.UserID = userID
		msg.SessionID = sessionID
		entity.WSBroadcast <- msg
		select {
		case entity.WSBroadcast <- msg:
			log.Printf("Message successfully sent to WSBroadcast: %+v", msg)
		default:
			log.Printf("WSBroadcast channel is full. Dropping message: %+v", msg)
		}
	}
}

// SendErrorMessage processes errors and sends them to the corresponding client.
func SendErrorMessage(msg *entity.WSMessage, errMsg *entity.ErrorInfo) {
	roomInfoMsg := &entity.RoomInfo{}
	roomInfoMsg.ErrorInfo = errMsg
	utils.LogError(errMsg.Msg)

	// Retrieve all sessionIDs for the room
	if sessionIDs, ok := entity.RoomSessions[msg.RoomID]; ok {
		for _, sessionID := range sessionIDs {
			// Find the client associated with the sessionID
			if client, exists := entity.WSClients[sessionID]; exists && client.UserID == msg.UserID {
				// Create an error message
				message, err := CreateMessage(roomInfoMsg)
				if err != nil {
					fmt.Println("Error creating error message:", err)
					continue
				}
				msg.Message = message

				// encrypt the message
				// encryptedMessage, err := utils.EncryptAES(message)
				// if err != nil {
				// 	fmt.Println("Error encrypting message:", err)
				// 	continue
				// }

				// // Set the encrypted message
				// msg.Message = encryptedMessage

				// Attempt to send the error message
				err = client.Conn.WriteJSON(msg)
				if err != nil {
					// Mark the client as closed (instead of immediate removal)
					client.Closed = true
					// Optionally retry sending the message (if needed)
					// Retry logic can be implemented here

					// Remove the client only after retries or severe errors
					closeAndRemoveClient(client, sessionID, msg.RoomID)
				}
			}
		}
	}

	// If the room has no active sessions, delete it
	if len(entity.RoomSessions[msg.RoomID]) == 0 {
		delete(entity.RoomSessions, msg.RoomID)
	}
}

func cleanGameInfo(ctx context.Context, userID uint) *entity.ErrorInfo {
	var errInfo *entity.ErrorInfo
	err := mysql.GormMysqlDB.Transaction(func(tx *gorm.DB) error {
		// frog_user_cards 제거
		errInfo := repository.DeleteAllFrogUserCards(ctx, tx, userID)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}
		// frog_room_users 제거
		errInfo = repository.DeleteAllFrogRoomUsers(ctx, tx, userID)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}
		// rooms 제거
		errInfo = repository.DeleteAllRooms(ctx, tx, userID)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}

		return nil
	})

	if err != nil {
		return errInfo
	}

	return nil
}
