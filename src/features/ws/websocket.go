package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"main/features/ws/model/entity"
	"main/features/ws/repository"
	"main/utils"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	// 클라이언트에 메시지를 쓸 수 있는 시간입니다.
	WriteWait = 5 * time.Second // 3~5초

	// 클라이언트로부터 다음 퐁 메시지를 읽을 수 있는 시간입니다.
	PongWait = 20 * time.Second // 15~30초

	// 핑을 보낼 수 있는 주기입니다. (PongWait보다 짧아야 함)
	PingPeriod = 7 * time.Second // PongWait의 1/3~1/2
)

func WSHandleMessages(gameName string) {
	rabbitManager := utils.GetRabbitMQManager()

	// 메시지 발행
	go func() {
		for {
			msg := <-entity.WSBroadcast
			msgBytes, err := json.Marshal(msg)
			if err != nil {
				log.Printf("Failed to marshal WSMessage: %v", err)
				continue
			}

			if err := rabbitManager.PublishMessage(gameName, msgBytes); err != nil {
				log.Printf("Failed to publish message to frog queue: %v", err)
			}
		}
	}()

	// 메시지 소비
	go func() {
		for {
			channel, err := rabbitManager.GetChannel(gameName)
			if err != nil {
				log.Printf("Failed to get channel for frog queue: %v", err)
				time.Sleep(5 * time.Second) // 재시도 대기
				continue
			}

			msgs, err := channel.Consume(
				gameName, // Queue name
				"",       // Consumer tag
				false,    // Auto-ack
				false,    // Exclusive
				false,    // No-local
				false,    // No-wait
				nil,      // Arguments
			)
			if err != nil {
				log.Printf("Failed to register consumer for frog queue: %v", err)
				time.Sleep(5 * time.Second) // 재시도 대기
				continue
			}

			log.Printf("Waiting for messages for game: %s (frog)", gameName)
			for msg := range msgs {
				processMessage(gameName, msg)
			}
		}
	}()
}

func processMessage(gameName string, d amqp.Delivery) {
	var msg entity.WSMessage
	// Parse JSON message
	err := json.Unmarshal(d.Body, &msg)
	if err != nil {
		log.Printf("Failed to unmarshal JSON for %s: %v", gameName, err)
		d.Nack(false, false) // Reject message, don't requeue
		return
	}
	utils.LogInfo(fmt.Sprintf("[FROG] Received message: %v \n", msg))
	var errInfo *entity.ErrorInfo
	// Handle events
	switch msg.Event {
	case "QUIT_GAME":
		errInfo = CloseEventWebsocket(&msg)
	case "START":
		errInfo = StartEventWebsocket(&msg)
	case "DISCARD":
		errInfo = DiscardCardsEventWebsocket(&msg)
	case "IMPORT_SINGLE_CARD":
		errInfo = ImportSingleCardEventWebsocket(&msg)
	case "GAME_OVER":
		errInfo = GameOverEventWebsocket(&msg)
	case "CHAT":
		errInfo = ChatEventWebsocket(&msg)
	case "REQUEST_WIN":
		errInfo = RequestWinEventWebsocket(&msg)
	case "TIME_OUT_DISCARD":
		errInfo = TimeOutDiscardCardsEventWebsocket(&msg)
	case "MATCH":
		errInfo = MatchEventWebsocket(&msg)
	case "CANCEL_MATCH":
		errInfo = CancelMatchEventWebsocket(&msg)
	case "PLAY_TOGETHER":
		errInfo = PlayTogetherEventWebsocket(&msg)
	case "JOIN_PLAY":
		errInfo = JoinPlayEventWebsocket(&msg)
	// New events for additional games
	case "DORA":
		errInfo = DoraEventWebsocket(&msg)
	case "IMPORT_CARDS":
		errInfo = ImportCardsEventWebsocket(&msg)
	case "LOAN":
		errInfo = LoanEventWebsocket(&msg)
	case "FAILED_LOAN":
		errInfo = FailedLoanEventWebsocket(&msg)
	case "SUCCESS_LOAN":
		errInfo = SuccessLoanEventWebsocket(&msg)
	default:
		log.Printf("Unknown event for %s: %s", gameName, msg.Event)
		d.Nack(false, false) // Reject message, don't requeue
		return
	}
	if errInfo != nil {
		SendErrorMessage(&msg, errInfo)
		d.Nack(false, false) // Reject message, don't requeue
	}

	// Acknowledge message after successful processing
	d.Ack(false)
}

// HandlePingPong manages PING/PONG messages to keep the connection alive.
func HandlePingPong(wsClient *entity.WSClient) {
	ws := wsClient.Conn

	// Setting up the Pong handler
	ws.SetReadDeadline(time.Now().Add(PongWait))
	ws.SetPongHandler(func(string) error {
		ws.SetReadDeadline(time.Now().Add(PongWait))
		return nil
	})

	ticker := time.NewTicker(PingPeriod)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:

			if err := ws.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(WriteWait)); err != nil {
				// 연결이 이미 닫혀있는지 확인
				if wsClient.Closed {
					AbnormalSendErrorMessage(wsClient.RoomID, wsClient.UserID, wsClient.SessionID)
					return
				}
				return
			}
		}
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
}

// Generate a new sessionID
func generateSessionID() string {
	return uuid.New().String() // Generate a new UUID
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

func disconnectClient(userID, roomID uint) {
	// RoomID에 연결된 모든 세션을 검색
	if sessionIDs, ok := entity.RoomSessions[roomID]; ok {
		for _, sessionID := range sessionIDs {
			// 특정 userID를 가진 클라이언트를 찾는다.
			if client, exists := entity.WSClients[sessionID]; exists && client.UserID == userID {
				// 클라이언트 연결 종료
				client.Conn.Close()
				client.Closed = true
				// redis 세션 id 삭제
				errInfo := repository.RedisSessionDelete(context.TODO(), sessionID)
				if errInfo != nil {
					fmt.Printf("Failed to delete session: %v\n", errInfo.Msg)
				}

				// 세션 및 클라이언트 데이터 정리
				delete(entity.WSClients, sessionID)
				removeSessionFromRoom(roomID, sessionID)

				fmt.Printf("User %d disconnected from room %d\n", userID, roomID)
				break
			}
		}
	} else {
		fmt.Printf("Room %d does not exist or has no active sessions\n", roomID)
	}
}
