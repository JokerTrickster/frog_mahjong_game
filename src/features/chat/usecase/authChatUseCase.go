package usecase

import (
	"context"
	_interface "main/features/chat/model/interface"
	"main/utils/db/mysql"
	"time"
)

type AuthChatUseCase struct {
	Repository     _interface.IAuthChatRepository
	ContextTimeout time.Duration
}

func NewAuthChatUseCase(repo _interface.IAuthChatRepository, timeout time.Duration) _interface.IAuthChatUseCase {
	return &AuthChatUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *AuthChatUseCase) Auth(c context.Context, userID uint) (string, error) {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()

	userInfo, err := d.Repository.FindOneUserInfo(ctx, userID)
	if err != nil {
		return "", err
	}
	secret := GenerateSecret(userID)
	chatDTO := &mysql.Chats{
		UserID: int(userID),
		Secret: secret,
		Name:   userInfo.Name,
	}

	err = d.Repository.InsertOneChat(ctx, chatDTO)
	if err != nil {
		return "", err
	}
	return secret, nil

}
