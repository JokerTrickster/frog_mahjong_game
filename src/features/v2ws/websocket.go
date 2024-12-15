package v2ws

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"main/features/v2ws/model/entity"
	"main/utils"
	"time"

	"github.com/gorilla/websocket"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	// 클라이언트에 메시지를 쓸 수 있는 시간입니다.
	WriteWait = 10 * time.Second

	// 클라이언트로부터 다음 퐁 메시지를 읽을 수 있는 시간입니다.
	PongWait = 30 * time.Second

	// 핑을 보낼 수 있는 기간입니다. (PongWait 보다 작아야 된다.)
	PingPeriod = (PongWait * 9) / 10
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
			// 연결이 이미 닫혀있는지 확인
			if wsClient.IsClosed() {
				return
			}
			if err := ws.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(WriteWait)); err != nil {
				fmt.Println("Error sending ping:", err)
				AbnormalErrorHandling(wsClient.RoomID, wsClient.UserID)
				return
			}

		}

	}
}

func ErrorHandling(msg *entity.WSMessage, roomID uint, userID uint, err *entity.RoomInfo) {
	// 에러 처리
	if clients, ok := entity.WSClients[roomID]; ok {
		for client := range clients {
			//이벤트 요청한 유저에게 에러 메시지 전송
			if clients[client].UserID == userID {
				message, err := CreateMessage(err)
				if err != nil {
					fmt.Println(err)
				}
				msg.Message = message
				err = client.WriteJSON(msg)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}
}
