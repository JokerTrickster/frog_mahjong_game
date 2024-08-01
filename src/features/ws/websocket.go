package ws

import (
	"fmt"
	"main/features/ws/model/entity"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// WriteWait is the time allowed to write a message to the client.
	WriteWait = 10 * time.Second

	// PongWait is the time allowed to read the next pong message from the client.
	PongWait = 60 * time.Second

	// PingPeriod is the time period to send pings. Must be less than PongWait.
	PingPeriod = (PongWait * 9) / 10
)

func WSHandleMessages() {

	for {
		msg := <-entity.WSBroadcast
		fmt.Println("현재 방 인원수 : ", len(entity.WSClients[msg.RoomID]))
		switch msg.Event {
		case "JOIN": // 방 참여
			JoinEventWebsocket(&msg)
		case "CLOSE": // 방 나가기
			CloseEventWebsocket(&msg)
		case "READY": // 게임 준비
			ReadyEventWebsocket(&msg)
		case "READY_CANCEL": // 게임 준비를 취소
			ReadyCancelEventWebsocket(&msg)
		case "START": // 게임 시작
			StartEventWebsocket(&msg)

		}
	}
}

// HandlePingPong manages PING/PONG messages to keep the connection alive.
func HandlePingPong(conn *websocket.Conn) {
	// Setting up the Pong handler
	conn.SetReadDeadline(time.Now().Add(PongWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(PongWait))
		return nil
	})

	ticker := time.NewTicker(PingPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(WriteWait)); err != nil {
				fmt.Println("Error sending ping:", err)
				return
			}
		}
	}
}
