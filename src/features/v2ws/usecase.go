package v2ws

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"main/features/v2ws/model/entity"
	"main/features/v2ws/model/request"
	"main/features/v2ws/repository"
	"main/utils"
	"main/utils/db/mysql"
	"strings"
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

func DiscardCreateRoomInfoMSG(ctx context.Context, preloadUsers []entity.RoomUsers, playTurn int, roomInfoError *entity.ErrorInfo, selectCardID int) *entity.RoomInfo {
	roomInfoMsg := entity.RoomInfo{}
	allReady := true
	timer := 30
	roomID := 0
	password := ""
	missionIDs := make([]int, 0)
	pickedCount := 0
	//유저 정보 저장
	for _, roomUser := range preloadUsers {
		timer = roomUser.Room.Timer
		password = roomUser.Room.Password
		user := entity.User{
			ID:                  uint(roomUser.UserID),
			Coin:                roomUser.User.Coin,
			Name:                roomUser.User.Name,
			Email:               roomUser.User.Email,
			ProfileID:           roomUser.User.ProfileID,
			PlayerState:         "picking",
			MissionSuccessCount: len(roomUser.UserMissions),
		}
		for _, item := range roomUser.UserItems {
			userItem := entity.Item{
				ItemID:        uint(item.ItemID),
				RemainingUses: item.RemainingUses,
			}
			user.Items = append(user.Items, &userItem)
		}
		ownedCount := 0
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
			} else if card.State == "picked" {
				user.PickedCards = append(user.PickedCards, &entity.Card{
					CardID: uint(card.CardID),
					UserID: uint(card.UserID),
				})
				user.DiscardedCards = append(user.DiscardedCards, &entity.Card{
					CardID: uint(card.CardID),
					UserID: uint(card.UserID),
				})
				ownedCount++
			}

		}
		if ownedCount == 1 {
			user.PlayerState = "done"
			pickedCount++
		}
		roomID = roomUser.RoomID

		if roomUser.Room.OwnerID == roomUser.UserID {
			user.IsOwner = true
		}
		if len(missionIDs) == 0 {
			for _, mission := range roomUser.RoomMission {
				missionIDs = append(missionIDs, mission.MissionID)
			}
		}

		roomInfoMsg.Users = append(roomInfoMsg.Users, &user)
	}

	//게임 정보 저장
	gameInfo := entity.GameInfo{
		PlayTurn:   playTurn,
		AllReady:   allReady,
		Timer:      timer,
		IsFull:     false,
		RoomID:     uint(roomID),
		Password:   password,
		MissionIDs: missionIDs,
		AllPicked:  false,
	}
	if pickedCount == len(preloadUsers) {
		gameInfo.AllPicked = true
	}
	openCards, err := repository.FindAllOpenCards(ctx, roomID)
	if err != nil {
		fmt.Println(err)
	}
	gameInfo.OpenCards = openCards
	roomInfoMsg.GameInfo = &gameInfo
	if roomInfoError != nil {
		roomInfoMsg.ErrorInfo = roomInfoError
	}
	return &roomInfoMsg

}

func CreateRoomInfoMSG(ctx context.Context, preloadUsers []entity.RoomUsers, playTurn int, roomInfoError *entity.ErrorInfo, selectCardID int) *entity.RoomInfo {
	roomInfoMsg := entity.RoomInfo{}
	allReady := true
	timer := 30
	roomID := 0
	password := ""
	pickedCount := 0
	missionIDs := make([]int, 0)
	var startTime int64

	//유저 정보 저장
	for _, roomUser := range preloadUsers {
		timer = roomUser.Room.Timer
		password = roomUser.Room.Password
		user := entity.User{
			ID:                  uint(roomUser.UserID),
			Coin:                roomUser.User.Coin,
			Name:                roomUser.User.Name,
			Email:               roomUser.User.Email,
			ProfileID:           roomUser.User.ProfileID,
			PlayerState:         "picking",
			MissionSuccessCount: len(roomUser.UserMissions),
		}
		//유저 상태 변경 (비정상적인 경우)
		if roomUser.RoomUsers.PlayerState == "disconnected" {
			user.PlayerState = "disconnected"
		}

		// 아이템 정보 저장
		for _, item := range roomUser.UserItems {
			userItem := entity.Item{
				ItemID:        uint(item.ItemID),
				RemainingUses: item.RemainingUses,
			}
			user.Items = append(user.Items, &userItem)
		}

		ownedCount := 0
		for _, card := range roomUser.Cards {
			if card.State == "owned" {
				user.Cards = append(user.Cards, &entity.Card{
					CardID: uint(card.CardID),
					UserID: uint(card.UserID),
				})
				ownedCount++
			} else if card.State == "discard" {
				user.DiscardedCards = append(user.DiscardedCards, &entity.Card{
					CardID: uint(card.CardID),
					UserID: uint(card.UserID),
				})
			} else if card.State == "picked" {
				user.Cards = append(user.Cards, &entity.Card{
					CardID: uint(card.CardID),
					UserID: uint(card.UserID),
				})
				user.PickedCards = append(user.PickedCards, &entity.Card{
					CardID: uint(card.CardID),
					UserID: uint(card.UserID),
				})
				ownedCount++
			}

		}
		if ownedCount == 4 {
			user.PlayerState = "done"
			pickedCount++
		}
		roomID = roomUser.RoomID
		// 방장 여부 추가
		if roomUser.Room.OwnerID == roomUser.UserID {
			user.IsOwner = true
		}
		//시작 시간 추가
		if !roomUser.Room.StartTime.IsZero() {
			// 시작 시간을 epoch time milliseconds로 변환 +3초 추가
			startTime = roomUser.Room.StartTime.UnixNano()/int64(time.Millisecond) + 5000
		}

		// 미션 정보 저장
		if len(missionIDs) == 0 {
			for _, mission := range roomUser.RoomMission {
				missionIDs = append(missionIDs, mission.MissionID)
			}
		}

		roomInfoMsg.Users = append(roomInfoMsg.Users, &user)
	}

	//게임 정보 저장
	gameInfo := entity.GameInfo{
		PlayTurn:   playTurn,
		AllReady:   allReady,
		Timer:      timer,
		IsFull:     false,
		RoomID:     uint(roomID),
		Password:   password,
		Winner:     0,
		MissionIDs: missionIDs,
		AllPicked:  false,
		StartTime:  startTime,
	}
	if pickedCount == len(preloadUsers) {
		gameInfo.AllPicked = true
	}

	openCards, err := repository.FindAllOpenCards(ctx, roomID)
	if err != nil {
		fmt.Println(err)
	}
	gameInfo.OpenCards = openCards

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

func CalcScore(cards []*mysql.Cards, score int) error {
	if score >= 5 {
		return nil
	}
	return fmt.Errorf("점수가 부족합니다.")
}

const (
	CONSECUTIVE_PAIRS = iota + 1 // 연속된 숫자 2쌍
	IDENTICAL_PAIRS              //같은 숫자 2쌍
)

// 숫자 문자열과 대응하는 숫자를 맵으로 정의
var numberMap = map[string]int{
	"one":   1,
	"two":   2,
	"three": 3,
	"four":  4,
	"five":  5,
	"six":   6,
	"seven": 7,
	"eight": 8,
	"nine":  9,
}

func CalcMission(missionID int, cards entity.WSRequestWinEntity) (bool, error) {
	switch missionID {
	case CONSECUTIVE_PAIRS:
		return CalcConsecutivePairs(cards), nil
	case IDENTICAL_PAIRS:
		return CalcIdenticalPairs(cards), nil
	}
	return false, nil
}
func CalcConsecutivePairs(cards entity.WSRequestWinEntity) bool {
	result := true
	// 연속된 숫자 2쌍인지 체크
	for i := 0; i < len(cards.Cards); i += 3 {
		card1Int := convertToNumber(cards.Cards[i].Name)
		card2Int := convertToNumber(cards.Cards[i+1].Name)
		card3Int := convertToNumber(cards.Cards[i+2].Name)

		if card1Int+1 == card2Int && card2Int+1 == card3Int {
			continue
		} else {
			result = false
		}
	}
	return result
}
func convertToNumber(word string) int {
	word = strings.ToLower(word) // 대소문자 구분을 없애기 위해 소문자로 변환
	number, _ := numberMap[word]
	return number
}

func CalcIdenticalPairs(cards entity.WSRequestWinEntity) bool {
	// 같은 숫자 2쌍인지 체크
	result := true
	for i := 0; i < len(cards.Cards); i += 3 {
		if cards.Cards[i].Name == cards.Cards[i+1].Name && cards.Cards[i+2].Name == cards.Cards[i+3].Name {
			continue
		} else {
			result = false
		}
	}
	return result
}
func CreateUserMissionDTO(missionEntity entity.V2WSMissionEntity) *mysql.UserMissions {
	userMissionDTO := mysql.UserMissions{
		UserID:    int(missionEntity.UserID),
		MissionID: int(missionEntity.MissionID),
		RoomID:    int(missionEntity.RoomID),
	}
	return &userMissionDTO
}

func CreateUserMissionCardDTO(missionEntity entity.V2WSMissionEntity, userMissionID int) *[]mysql.UserMissionCards {
	var userMissionCardDTO []mysql.UserMissionCards
	for _, cardID := range missionEntity.Cards {
		userMissionCardDTO = append(userMissionCardDTO, mysql.UserMissionCards{
			UserMissionID: userMissionID,
			CardID:        cardID,
		})
	}
	return &userMissionCardDTO
}

func CreateUserBirdCardDTO(importSingleCard entity.WSImportSingleCardEntity) *mysql.UserBirdCards {
	userBirdCardDTO := mysql.UserBirdCards{
		UserID: int(importSingleCard.UserID),
		RoomID: int(importSingleCard.RoomID),
		CardID: int(importSingleCard.CardID),
		State:  "picked",
	}
	return &userBirdCardDTO
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
	newErr := repository.MatchRedisSessionSet(context.TODO(), newSessionID, roomID)
	if newErr != nil {
		fmt.Println(newErr)
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

	// 방에 세션 추가
	entity.RoomSessions[roomID] = append(entity.RoomSessions[roomID], sessionID)
	fmt.Println(len(entity.RoomSessions[roomID]))
	// 핑/퐁 핸들링 시작
	go HandlePingPong(wsClient)

	// 메시지 처리 루프 시작
	go readMessages(ws, sessionID, roomID, userID)

}

// 메시지 읽기 및 처리
func readMessages(ws *websocket.Conn, sessionID string, roomID uint, userID uint) {
	client := entity.WSClients[sessionID]
	fmt.Println(client)
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
	}
}

// 클라이언트에 메시지 전송
func sendMessageToClients(roomID uint, msg *entity.WSMessage) {
	// 로그 메시지 생성
	utils.LogError(msg.Message)

	// 메시지 암호화
	encryptedMessage, err := utils.EncryptAES(msg.Message)
	if err != nil {
		fmt.Printf("Failed to encrypt message: %v\n", err)
		return
	}
	msg.Message = encryptedMessage

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

	// 메시지 암호화
	encryptedMessage, err := utils.EncryptAES(msg.Message)
	if err != nil {
		fmt.Printf("Failed to encrypt message: %v\n", err)
		return
	}
	msg.Message = encryptedMessage

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

				// encrypt the message
				encryptedMessage, err := utils.EncryptAES(message)
				if err != nil {
					fmt.Println("Error encrypting message:", err)
					continue
				}

				// Set the encrypted message
				msg.Message = encryptedMessage

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

func CreateMessage(roomInfoMsg *entity.RoomInfo) (string, error) {
	// 구조체를 JSON 문자열로 변환 (마샬링)
	jsonData, err := json.Marshal(roomInfoMsg)
	if err != nil {
		return "", fmt.Errorf("JSON 마샬링 에러: %s", err)
	}

	return string(jsonData), nil
}

func CreateErrorMessage(errCode int, errType, errMsg string) *entity.ErrorInfo {
	result := &entity.ErrorInfo{
		Code: errCode,
		Type: errType,
		Msg:  errMsg,
	}
	return result
}

func cleanGameInfo(ctx context.Context, userID uint) *entity.ErrorInfo {
	var errInfo *entity.ErrorInfo
	err := mysql.GormMysqlDB.Transaction(func(tx *gorm.DB) error {
		// user_bird_cards 제거
		errInfo := repository.DeleteAllUserBirdCards(ctx, tx, userID)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}
		// frog_room_users 제거
		errInfo = repository.DeleteAllRoomUsers(ctx, tx, userID)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}
		// rooms 제거
		errInfo = repository.DeleteAllRooms(ctx, tx, userID)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}
		// user_missions 제거
		errInfo = repository.DeleteAllUserMissions(ctx, tx, userID)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}
		// user_items 제거
		errInfo = repository.DeleteAllUserItems(ctx, tx, userID)
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
