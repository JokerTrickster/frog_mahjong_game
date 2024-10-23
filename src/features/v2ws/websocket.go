package v2ws

import (
	"fmt"
	"main/features/v2ws/model/entity"
	"main/utils"
	"time"

	"github.com/gorilla/websocket"
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

	for {
		msg := <-entity.WSBroadcast
		logging := utils.Log{}
		logging.V2MakeWSLog(msg)
		utils.LogInfo(logging)
		switch msg.Event {
		case "QUIT_GAME": // 방 나가기
			CloseEventWebsocket(&msg)
		case "START": // 게임 시작
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
