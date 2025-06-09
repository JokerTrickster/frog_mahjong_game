package find_it

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"main/features/find_it/model/entity"
	"main/features/find_it/repository"
	"main/utils"
	_aws "main/utils/aws"
	"main/utils/db/mysql"
	"time"

	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

func CreateRoundImages(roomID uint, imagesDTO []*mysql.FindItImages) []*mysql.FindItRoundImages {
	roundImagesDTO := []*mysql.FindItRoundImages{}
	round := 1
	for _, imageDTO := range imagesDTO {
		roundImage := &mysql.FindItRoundImages{
			RoomID:     int(roomID),
			ImageSetId: int(imageDTO.ID),
			Round:      round,
		}
		round++
		roundImagesDTO = append(roundImagesDTO, roundImage)
	}
	return roundImagesDTO
}

func Deepcopy(src entity.MessageInfo) entity.MessageInfo {
	var dst entity.MessageInfo
	b, _ := json.Marshal(src)
	json.Unmarshal(b, &dst)
	return dst
}

func CalcPlayTurn(playTurn, playerCount int) int {
	return (playTurn % playerCount) + 1
}
func CreateRoomSetting(roomID uint) *mysql.FindItRoomSettings {
	roomSetting := &mysql.FindItRoomSettings{
		RoomID:             int(roomID),
		Timer:              60,
		Lifes:              3,
		ItemHintCount:      3,
		ItemTimerStopCount: 3,
		Round:              1,
	}
	return roomSetting
}
func CreateMatchRoomUserDTO(userID uint, roomID uint) *mysql.GameRoomUsers {
	roomUser := &mysql.GameRoomUsers{
		RoomID:      int(roomID),
		UserID:      int(userID),
		PlayerState: "wait",
	}
	return roomUser
}

func CreateMessageInfoMSG(ctx context.Context, preloadUsers []entity.PreloadUsers, playTurn int, MessageInfoError *entity.ErrorInfo, selectCardID int) *entity.MessageInfo {
	MessageInfoMsg := entity.MessageInfo{}
	timer := 120
	hintCount := 0
	timerStopCount := 0
	life := 0
	round := 0
	roomID := 0
	password := ""
	correctCount := 0
	imageID := 0
	roundCount := 0
	gameOver := false
	var startTime int64

	//유저 정보 저장
	for _, roomUser := range preloadUsers {
		user := entity.User{
			ID:        uint(roomUser.UserID),
			Name:      roomUser.User.Name,
			Email:     roomUser.User.Email,
			ProfileID: roomUser.User.ProfileID,
		}
		if roomUser.RoomSetting != nil {
			timer = roomUser.RoomSetting.Timer
			hintCount = roomUser.RoomSetting.ItemHintCount
			timerStopCount = roomUser.RoomSetting.ItemTimerStopCount
			life = roomUser.RoomSetting.Lifes
			round = roomUser.RoomSetting.Round
		}
		if roomUser.Room != nil {
			if roomUser.Room.Password != "" {
				password = roomUser.Room.Password
			}

			roomID = int(roomUser.RoomID)
			// 방장 여부 추가
			if roomUser.Room.OwnerID == int(roomUser.UserID) {
				user.IsOwner = true
			}
			//시작 시간 추가
			if !roomUser.Room.StartTime.IsZero() {
				// 시작 시간을 epoch time milliseconds로 변환 +3초 추가
				startTime = roomUser.Room.StartTime.UnixNano()/int64(time.Millisecond) + 5000
			}
		}
		// ✅ 맞힌 정보 저장 (x_position, y_position 추가)
		correctIDList := []int{}

		for _, userCorrect := range roomUser.UserCorrectPositions {
			if userCorrect.RoomID == roomID && userCorrect.Round == round && userCorrect.UserID == int(roomUser.UserID) {
				correctIDList = append(correctIDList, userCorrect.CorrectPositionID)
				correctCount++
			}
		}
		correctPositions, _ := repository.FindAllCorrectPositions(ctx, correctIDList)
		for _, correctPosition := range correctPositions {
			position := entity.Position{
				X: correctPosition.XPosition,
				Y: correctPosition.YPosition,
			}
			user.CorrectPositions = append(user.CorrectPositions, position)
		}

		if imageID == 0 {
			for _, roundImage := range roomUser.RoundImages {
				if roomUser.RoomSetting.Round == roundImage.Round {
					imageID = roundImage.ImageSetId
					break
				}
			}
			roundCount = len(roomUser.RoundImages)
		}

		MessageInfoMsg.Users = append(MessageInfoMsg.Users, &user)
	}
	// 이미지 정보를 가져온다.
	ImageInfo := entity.ImageInfo{}
	if imageID != 0 {
		roundImage, err := repository.FindOneRoundImage(ctx, imageID)
		if err != nil {
			MessageInfoError = err
		}
		//s3 이미지 URL
		normalSignedUrl, _ := _aws.ImageGetSignedURL(context.TODO(), roundImage.NormalImageUrl, _aws.ImgTypeFindIt)
		abnormalSignedUrl, _ := _aws.ImageGetSignedURL(context.TODO(), roundImage.AbnormalImageUrl, _aws.ImgTypeFindIt)
		ImageInfo.ID = imageID
		ImageInfo.NormalImageUrl = normalSignedUrl
		ImageInfo.AbnormalImageUrl = abnormalSignedUrl

	}

	//게임 정보 저장
	gameInfo := entity.GameInfo{
		AllReady:       true,
		Timer:          timer,
		IsFull:         true,
		RoomID:         uint(roomID),
		Password:       password,
		StartTime:      startTime,
		ItemTimerCount: timerStopCount,
		ItemHintCount:  hintCount,
		Round:          round,
		Life:           life,
		CorrectCount:   correctCount,
		ImageInfo:      &ImageInfo,
		TimerUsed:      false,
		HintPosition:   nil,
		RoundCount:     roundCount,
		GameOver:       gameOver,
	}

	MessageInfoMsg.GameInfo = &gameInfo
	if MessageInfoError != nil {
		MessageInfoMsg.ErrorInfo = MessageInfoError
	}
	return &MessageInfoMsg

}
func CreateChatMessage(chatInfoMsg *entity.ChatInfo) (string, error) {
	// 구조체를 JSON 문자열로 변환 (마샬링)
	jsonData, err := json.Marshal(chatInfoMsg)
	if err != nil {
		return "", fmt.Errorf("JSON 마샬링 에러: %s", err)
	}

	return string(jsonData), nil
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
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic in sendMessageToClients: %v", r)
		}
	}()
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
func SendMessageToClient(roomID uint, msg *entity.WSMessage) {
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

// SendErrorMessage processes errors and sends them to the corresponding client.
func SendErrorMessage(msg *entity.WSMessage, errMsg *entity.ErrorInfo) {
	MessageInfoMsg := &entity.MessageInfo{}
	MessageInfoMsg.ErrorInfo = errMsg
	utils.LogError(errMsg.Msg)
	// Retrieve all sessionIDs for the room
	if sessionIDs, ok := entity.RoomSessions[msg.RoomID]; ok {
		for _, sessionID := range sessionIDs {
			// Find the client associated with the sessionID
			if client, exists := entity.WSClients[sessionID]; exists && client.UserID == msg.UserID {
				// Create an error message
				message, err := CreateMessage(MessageInfoMsg)
				if err != nil {
					fmt.Println("Error creating error message:", err)
					continue
				}

				// // encrypt the message
				// encryptedMessage, err := utils.EncryptAES(message)
				// if err != nil {
				// 	fmt.Println("Error encrypting message:", err)
				// 	continue
				// }

				// // Set the encrypted message
				// msg.Message = encryptedMessage
				msg.Message = message

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

func CreateMessage(MessageInfoMsg *entity.MessageInfo) (string, error) {
	// 구조체를 JSON 문자열로 변환 (마샬링)
	jsonData, err := json.Marshal(MessageInfoMsg)
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
		// GameRooms 제거
		errInfo = repository.DeleteAllGameRooms(ctx, tx, userID)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}
		// GameRoomUsers 제거
		errInfo = repository.DeleteAllGameRoomUsers(ctx, tx, userID)
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
