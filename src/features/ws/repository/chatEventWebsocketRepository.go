package repository

import (
	"context"
	"fmt"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func ChatInsertOneChat(ctx context.Context, tx *gorm.DB, chatDTO *mysql.Chats) error {
	err := tx.Create(&chatDTO).Error
	if err != nil {
		return fmt.Errorf("채팅 저장 실패 %v", err)
	}
	return nil
}
