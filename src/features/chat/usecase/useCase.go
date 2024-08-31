package usecase

import (
	"crypto/rand"
	"encoding/hex"
	"main/features/chat/model/entity"
	"main/features/chat/model/request"
	"main/features/chat/model/response"
	"main/utils/db/mysql"
)

func GenerateSecret(userID uint) string {
	// 32바이트의 랜덤 데이터 생성
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		panic(err)
	}

	// 바이트 배열을 16진수 문자열로 변환하여 출력
	keyString := hex.EncodeToString(key)
	return keyString
}

func CreateChatHistoryEntitySQL(req *request.ReqHistory) *entity.HistoryEntitySQL {

	entitySQL := &entity.HistoryEntitySQL{
		Page:     req.Page,
		PageSize: req.PageSize,
		RoomID:   req.RoomID,
	}

	return entitySQL
}

func CreateResHistoryChat(chats []*mysql.Chats, total int) (response.ResHistoryChat, error) {
	res := response.ResHistoryChat{
		Total: total,
	}
	history := make([]response.HistoryChat, 0)
	for _, chat := range chats {
		msg := response.HistoryChat{
			UserID:  uint(chat.UserID),
			Name:    chat.Name,
			Message: chat.Message,
			Created: chat.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		history = append(history, msg)
	}
	res.Chats = history

	return res, nil
}
