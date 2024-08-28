package usecase

import (
	"context"
	"main/features/chat/model/interface/mocks"
	"main/utils/db/mysql"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"gopkg.in/go-playground/assert.v1"
	"gorm.io/gorm"
)

// TestMessageChatUseCase_Message 함수는 MessageChatUseCase 의 Message 메서드를 테스트합니다.
// 테이블 기반 테스트를 사용하여 여러 시나리오를 테스트합니다.
// given-when-then 패턴을 사용하여 테스트를 작성합니다.
// 구조체 필드 :  string, *mysql.Chats,wantErr error
// 반환값 : *mysql.Chats, error
// 테스트 케이스:
// 1. 채팅 정보를 찾을 수 있는 경우
// 테스트 경로: src/features/chats/usecase/messageChatUseCase_test.go

func TestMessageChatUseCase_Message(t *testing.T) {
	tests := []struct {
		name     string
		secret   string
		chatInfo *mysql.Chats
		wantErr  error
	}{
		{
			name:   "채팅 정보를 찾을 수 있는 경우",
			secret: "test",
			chatInfo: &mysql.Chats{
				Model: gorm.Model{
					ID: 1,
				},
				UserID:  1,
				Message: "test",
				Name:    "ryan",
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			mockMessageChatRepository := new(mocks.IMessageChatRepository)
			mockMessageChatRepository.On("FindOneChat", mock.Anything, mock.Anything).Return(tt.chatInfo, tt.wantErr) //mock
			us := NewMessageChatUseCase(mockMessageChatRepository, 8*time.Second)
			// when
			chatInfo, err := us.Message(context.Background(), tt.secret)
			// then
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, err)
				return
			}
			assert.Equal(t, tt.chatInfo, chatInfo)
		})
	}
}
