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
		switch msg.Event {
		case "JOIN": // 방 참여
			JoinEventWebsocket(&msg)
		case "QUIT_GAME": // 방 나가기
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
		case "LOAN":
			LoanEventWebsocket(&msg)
		case "FAILED_LOAN":
			FailedLoanEventWebsocket(&msg)
		case "SUCCESS_LOAN":
			SuccessLoanEventWebsocket(&msg)
		case "REQUEST_WIN":
			RequestWinEventWebsocket(&msg)
		case "GAME_OVER":
			GameOverEventWebsocket(&msg)
		case "ROOM_OUT":
			RoomOutEventWebsocket(&msg)
		case "CHAT":
			ChatEventWebsocket(&msg)
		case "TIME_OUT_DISCARD":
			TimeOutDiscardCardsEventWebsocket(&msg)
		case "MATCH":
			MatchEventWebsocket(&msg)
		case "CANCEL_MATCH":
			CancelMatchEventWebsocket(&msg)
		case "PLAY_TOGETHER":
			PlayTogetherEventWebsocket(&msg)
		}
	}
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
