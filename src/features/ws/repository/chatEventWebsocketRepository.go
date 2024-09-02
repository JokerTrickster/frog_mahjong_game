package repository

import (
	"context"
	"fmt"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

// 채팅 저장 후 채팅 ID 반환
func ChatInsertOneChat(ctx context.Context, tx *gorm.DB, chatDTO *mysql.Chats) (uint, error) {
	if err := tx.WithContext(ctx).Create(chatDTO).Error; err != nil {
		return 0, fmt.Errorf("챗 정보 저장 실패: %v", err)
	}

	return chatDTO.ID, nil
}
