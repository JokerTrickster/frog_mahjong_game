package ws

import (
	"fmt"
	"main/features/ws/model/entity"
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
