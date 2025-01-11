package v2ws

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"main/features/v2ws/model/entity"
	"main/features/v2ws/repository"
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
	// 클라이언트에 메시지를 쓸 수 있는 시간입니다.
	WriteWait = 5 * time.Second // 3~5초

	// 클라이언트로부터 다음 퐁 메시지를 읽을 수 있는 시간입니다.
	PongWait = 20 * time.Second // 15~30초

	// 핑을 보낼 수 있는 주기입니다. (PongWait보다 짧아야 함)
	PingPeriod = 7 * time.Second // PongWait의 1/3~1/2
)

func WSHandleMessages(gameName string) {
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

			err = utils.V2MQCH.PublishWithContext(ctx,
				"",              // exchange
				utils.V2MQ.Name, // routing key
				false,           // mandatory
				false,           // immediate
				amqp.Publishing{
					ContentType: "application/json",
					Body:        msgBytes,
				})
			if err != nil {
				log.Printf("Failed to publish a message: %v", err)
			}
		}
	}()
	// Consume messages
	go func() {
		msgs, err := utils.V2MQCH.Consume(
			utils.V2MQ.Name, // Queue name
			"",              // Consumer tag
			false,           // Auto-ack (manual ack)
			false,           // Exclusive
			false,           // No-local
			false,           // No-wait
			nil,             // Arguments
		)
		if err != nil {
			log.Printf("Failed to register consumer for %s: %v", gameName, err)
		}

		log.Printf("Waiting for messages for game: %s", gameName)
		for msg := range msgs {
			processMessage(gameName, msg)
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

	var errInfo *entity.ErrorInfo
	// 이벤트 처리
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
	case "MISSION":
		errInfo = MissionEventWebsocket(&msg)
	case "RANDOM":
		errInfo = RandomEventWebsocket(&msg)
	case "ITEM_CHANGE":
		errInfo = ItemChangeEventWebsocket(&msg)
	default:
		log.Printf("Unknown event: %s", msg.Event)
		d.Nack(false, false) // 알 수 없는 이벤트 -> 재처리하지 않음
		return
	}
	if errInfo != nil {
		SendErrorMessage(&msg, errInfo)
		d.Ack(false)
	}
	// Acknowledge message after successful processing
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
				// Notify all users in the same room about the disconnection
				// false이면
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
				newErr := repository.RedisSessionDelete(context.TODO(), sessionID)
				if newErr != nil {
					fmt.Printf("Failed to delete session: %v\n", newErr.Msg)
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
