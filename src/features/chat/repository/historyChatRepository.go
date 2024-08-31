package repository

import (
	"context"
	"main/features/chat/model/entity"
	_interface "main/features/chat/model/interface"
	"main/utils"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewHistoryChatRepository(gormDB *gorm.DB) _interface.IHistoryChatRepository {
	return &HistoryChatRepository{GormDB: gormDB}
}

func (d *HistoryChatRepository) FindChatHistory(ctx context.Context, entitySQL *entity.HistoryEntitySQL) ([]*mysql.Chats, error) {
	var chats []*mysql.Chats
	// 페이지네이션 처리
	err := d.GormDB.Model(&mysql.Chats{}).Where("deleted_at IS NULL and room_id = ?", entitySQL.RoomID).Limit(entitySQL.PageSize).Offset((entitySQL.Page - 1) * entitySQL.PageSize).Find(&chats).Error
	if err != nil {
		return nil, utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), err.Error(), utils.ErrFromMysqlDB)
	}
	return chats, nil
}

func (d *HistoryChatRepository) CountChatHistory(ctx context.Context, entitySQL *entity.HistoryEntitySQL) (int, error) {
	var count int64
	err := d.GormDB.Model(&mysql.Chats{}).Where("deleted_at IS NULL and room_id = ?", entitySQL.RoomID).Count(&count).Error
	if err != nil {
		return 0, utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), err.Error(), utils.ErrFromMysqlDB)
	}

	return int(count), nil
}
