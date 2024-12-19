package v2ws

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"main/features/v2ws/model/entity"
	"main/utils"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	amqp "github.com/rabbitmq/amqp091-go"
)

/*
	웹소켓 핑퐁 관리

PingPeriod : 서버가 클라이언트에게 Ping 메시지를 전송하는 주기 (PongWait보다 작아야 함 일반적으로 PongWait의 50 ~ 90%)
PongWait : 서버가 클라이언트로부터 Pong 메시지를 수신해야 하는 최대 대기 시간 (Pong 응답을 보내지 않으면 연결이 끊겼다고 판단)
WriteWait : 서버가 클라이언트에 데이터를 쓸 수 있는 최대 시간이다.
reconnectTime : 클라이언트가 연결을 잃었을 때 다시 연결을 시도할 수 있는 시간 (PongWait보다 크거나 같아야 된다. )
*/
const (
	WriteWait  = 10 * time.Second
	PongWait   = 30 * time.Second    // 30초마다 퐁 메시지를 수신
	PingPeriod = (PongWait * 5) / 10 // 6초마다 핑 메시지 전송
)

func WSHandleMessages() {
	// 웹소켓 메시지를 큐에 넣기
	go func() {
		for {
			msg := <-entity.WSBroadcast

			// 로그 생성
			logging := utils.Log{}
			logging.V2MakeWSLog(msg)
			utils.LogInfo(logging)

			// RabbitMQ에 메시지 발행
			msgBytes, err := json.Marshal(msg)
			if err != nil {
				log.Printf("Failed to marshal WSMessage: %v", err)
				continue
			}

			ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

			err = utils.MQCH.PublishWithContext(ctx,
				"",            // exchange
				utils.MQ.Name, // routing key
				false,         // mandatory
				false,         // immediate
				amqp.Publishing{
					ContentType: "application/json",
					Body:        msgBytes,
				})
			if err != nil {
				log.Printf("Failed to publish a message: %v", err)
			}
		}
	}()
	// 큐에서 메시지 소비 및 처리
	go func() {
		msgs, err := utils.MQCH.Consume(
			utils.MQ.Name, // Queue
			"",            // Consumer
			false,         // Auto-Ack (수동 Ack 사용)
			false,         // Exclusive
			false,         // No-local
			false,         // No-wait
			nil,           // Args
		)
		if err != nil {
			log.Fatalf("Failed to register a consumer: %v", err)
		}

		for d := range msgs {
			processMessage(d)
		}
	}()
}
func processMessage(d amqp.Delivery) {
	var msg entity.WSMessage

	// JSON 파싱
	err := json.Unmarshal(d.Body, &msg)
	if err != nil {
		log.Printf("Failed to unmarshal JSON: %v", err)
		d.Nack(false, false) // 메시지 처리 실패 -> 재처리하지 않음
		return
	}

	log.Printf(" [x] Received: %s", msg)

	// 이벤트 처리
	switch msg.Event {
	case "QUIT_GAME":
		CloseEventWebsocket(&msg)
	case "START":
		StartEventWebsocket(&msg)
	case "DISCARD":
		DiscardCardsEventWebsocket(&msg)
	case "IMPORT_SINGLE_CARD":
		ImportSingleCardEventWebsocket(&msg)
	case "GAME_OVER":
		GameOverEventWebsocket(&msg)
	case "CHAT":
		ChatEventWebsocket(&msg)
	case "REQUEST_WIN":
		RequestWinEventWebsocket(&msg)
	case "TIME_OUT_DISCARD":
		TimeOutDiscardCardsEventWebsocket(&msg)
	case "MATCH":
		MatchEventWebsocket(&msg)
	case "CANCEL_MATCH":
		CancelMatchEventWebsocket(&msg)
	case "PLAY_TOGETHER":
		PlayTogetherEventWebsocket(&msg)
	case "JOIN_PLAY":
		JoinPlayEventWebsocket(&msg)
	case "MISSION":
		MissionEventWebsocket(&msg)
	case "RANDOM":
		RandomEventWebsocket(&msg)
	case "ITEM_CHANGE":
		ItemChangeEventWebsocket(&msg)
	default:
		log.Printf("Unknown event: %s", msg.Event)
		d.Nack(false, false) // 알 수 없는 이벤트 -> 재처리하지 않음
		return
	}

	// 처리 성공 시 Ack
	d.Ack(false)
}

// HandlePingPong manages PING/PONG messages to keep the connection alive.
func HandlePingPong(wsClient *entity.WSClient) {
	ws := wsClient.Conn

	// Set initial deadline for Pong
	ws.SetReadDeadline(time.Now().Add(PongWait))
	ws.SetPongHandler(func(string) error {
		// Update the deadline on Pong receipt
		ws.SetReadDeadline(time.Now().Add(PongWait))
		return nil
	})

	ticker := time.NewTicker(PingPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:

			// Send Ping message
			if err := ws.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(WriteWait)); err != nil {
				fmt.Printf("Error sending ping for session %s: %v\n", wsClient.SessionID, err)
				wsClient.Closed = true
				// Handle abnormal connection termination
				AbnormalErrorHandling(wsClient.RoomID, wsClient.SessionID)
				return
			}
		}
	}
}

// ErrorHandling processes errors and sends them to the corresponding client.
func ErrorHandling(msg *entity.WSMessage, roomID uint, userID uint, roomError *entity.RoomInfo) {
	// Retrieve all sessionIDs for the room
	if sessionIDs, ok := entity.RoomSessions[roomID]; ok {
		for _, sessionID := range sessionIDs {
			// Find the client associated with the sessionID
			if client, exists := entity.WSClients[sessionID]; exists && client.UserID == userID {
				// Create an error message
				message, err := CreateMessage(roomError)
				if err != nil {
					fmt.Println("Error creating error message:", err)
					continue
				}
				msg.Message = message

				// Attempt to send the error message
				err = client.Conn.WriteJSON(msg)
				if err != nil {
					fmt.Printf("Error sending message to session %s (user %d): %v\n", sessionID, userID, err)

					// Mark the client as closed (instead of immediate removal)
					client.Closed = true

					// Optionally retry sending the message (if needed)
					// Retry logic can be implemented here

					// Remove the client only after retries or severe errors
					closeAndRemoveClient(client, sessionID, roomID)
				}
			}
		}
	}

	// If the room has no active sessions, delete it
	if len(entity.RoomSessions[roomID]) == 0 {
		delete(entity.RoomSessions, roomID)
		fmt.Printf("Room %d deleted as it has no active sessions.\n", roomID)
	}
}

// closeAndRemoveClient safely closes a client and removes it from session lists
func closeAndRemoveClient(client *entity.WSClient, sessionID string, roomID uint) {
	// Close the connection if not already closed
	if !client.Closed {
		client.Conn.Close()
		client.Closed = true
	}

	// Remove from WSClients and RoomSessions
	delete(entity.WSClients, sessionID)
	removeSessionFromRoom(roomID, sessionID)

	fmt.Printf("Client with session %s removed from room %d.\n", sessionID, roomID)
}

// Generate a new sessionID
func generateSessionID() string {
	return uuid.New().String() // Generate a new UUID
}

// Add a sessionID to the room
func addSessionToRoom(roomID uint, sessionID string) {
	entity.RoomSessions[roomID] = append(entity.RoomSessions[roomID], sessionID)
}

// Remove a sessionID from the room
func removeSessionFromRoom(roomID uint, sessionID string) {
	sessions := entity.RoomSessions[roomID]
	for i, id := range sessions {
		if id == sessionID {
			// Remove sessionID from the room
			entity.RoomSessions[roomID] = append(sessions[:i], sessions[i+1:]...)
			break
		}
	}
}
