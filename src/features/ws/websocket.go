package ws

import (
	"fmt"
	"main/features/ws/model/entity"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// 클라이언트에 메시지를 쓸 수 있는 시간입니다.
	WriteWait = 10 * time.Second

	// 클라이언트로부터 다음 퐁 메시지를 읽을 수 있는 시간입니다.
	PongWait = 40 * time.Second

	// 핑을 보낼 수 있는 기간입니다. (PongWait 보다 작아야 된다.)
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
		case "DORA":
			DoraEventWebsocket(&msg)
		case "IMPORT_CARDS":
			ImportCardsEventWebsocket(&msg)
		case "DISCARD":
			DiscardCardsEventWebsocket(&msg)
		case "IMPORT_SINGLE_CARD":
			ImportSingleCardEventWebsocket(&msg)
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
