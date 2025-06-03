package sequence

import (
	"context"
	"fmt"
	"log"
	"main/features/sequence/model/entity"
	"main/features/sequence/repository"
	"main/utils"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
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
	WriteWait = 3 * time.Second // 3~5초

	// 클라이언트로부터 다음 퐁 메시지를 읽을 수 있는 시간입니다.
	PongWait = 10 * time.Second // 15~30초

	// 핑을 보낼 수 있는 주기입니다. (PongWait보다 짧아야 함)
	PingPeriod = 5 * time.Second // PongWait의 1/3~1/2
)

func processMessage(gameName string, msg entity.WSMessage) {
	utils.LogInfo(fmt.Sprintf("[sequence] Received message: %v \n", msg))
	var errInfo *entity.ErrorInfo
	// 이벤트 처리
	switch msg.Event {
	case "MATCH":
		errInfo = MatchEventWebsocket(&msg)
	case "START":
		errInfo = StartEventWebsocket(&msg)
	case "TIME_OUT":
		errInfo = TimeOutEventWebsocket(&msg)
	case "MATCH_CANCEL":
		errInfo = CancelMatchEventWebsocket(&msg)
	case "TOGETHER":
		errInfo = PlayTogetherEventWebsocket(&msg)
	case "JOIN":
		errInfo = JoinPlayEventWebsocket(&msg)
	case "USE_CARD":
		errInfo = UseCardEventWebsocket(&msg)
	case "REMOVE_CARD":
		errInfo = RemoveCardEventWebsocket(&msg)
	case "GAME_OVER":
		errInfo = GameOverEventWebsocket(&msg)
	default:
		log.Printf("Unknown event: %s", msg.Event)
		return
	}
	if errInfo != nil {
		SendErrorMessage(&msg, errInfo)
	}
}

func WSHandleMessages(gameName string) {
	go func() {
		for {
			msg := <-entity.WSBroadcast
			processMessage(gameName, msg)
		}
	}()
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
