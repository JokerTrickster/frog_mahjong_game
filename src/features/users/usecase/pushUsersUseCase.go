package usecase

import (
	"context"
	"log"
	_interface "main/features/users/model/interface"
	"main/features/users/model/request"
	"main/utils"
	"time"

	"firebase.google.com/go/messaging"
)

type PushUsersUseCase struct {
	Repository     _interface.IPushUsersRepository
	ContextTimeout time.Duration
}

func NewPushUsersUseCase(repo _interface.IPushUsersRepository, timeout time.Duration) _interface.IPushUsersUseCase {
	return &PushUsersUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *PushUsersUseCase) Push(c context.Context, req *request.ReqPushUsers) error {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()

	// 어드민 유저인지 체크한다.
	if req.Role != "frogadmin" {
		return utils.ErrorMsg(ctx, utils.ErrBadParameter, utils.Trace(), utils.HandleError("only food admin can send message", req), utils.ErrFromClient)
	}
	// 1. 알람 여부를 체크한다.
	users, err := d.Repository.FindUsersForNotifications(ctx)
	if err != nil {
		return err
	}
	//TODO 추후 고루틴으로 처리할 예정
	for _, user := range users {

		// 2. 푸시 토큰을 가져온다.
		token, err := d.Repository.FindOnePushToken(ctx, uint(user.ID))
		if err != nil {
			return err
		}
		if token == "" {
			continue
		}
		// 3. 푸시를 보낸다.
		// 메시지 생성
		message := &messaging.Message{
			Token: token,
			Notification: &messaging.Notification{
				Title: req.Title,
				Body:  req.Message,
			},
		}

		// 메시지 전송
		response, err := utils.MessageClient.Send(ctx, message)
		if err != nil {
			log.Printf("error sending message: %v", err)
			return err
		}

		log.Printf("Successfully sent message: %s", response)
	}

	// 4. 메시지를 저장한다.

	return nil
}
