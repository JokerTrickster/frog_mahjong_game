package ws

import (
	"context"
	"main/features/ws/model/entity"
	"main/features/ws/repository"
)

func CloseEventWebsocket(msg *entity.WSMessage) *entity.ErrorInfo {
	//유저 상태를 변경한다. (대기실로 이동)
	ctx := context.Background()
	uID := msg.UserID
	rID := msg.RoomID

	//비즈니스 로직
	var errInfo *entity.ErrorInfo
	errInfo = cleanGameInfo(ctx, uID)
	if errInfo != nil {
		return errInfo
	}
	// 유저 정보를 업데이트 한다.
	errInfo = repository.CloseFindOneAndUpdateUser(ctx, uID)
	if errInfo != nil {
		return errInfo
	}

	// 정상적으로 연결을 끊는다.
	if sessionIDs, ok := entity.RoomSessions[rID]; ok {
		for _, sessionID := range sessionIDs {
			if client, exists := entity.WSClients[sessionID]; exists && client.UserID == uID {
				closeAndRemoveClient(client, sessionID, rID)
			}
		}
	}
	return nil
}
