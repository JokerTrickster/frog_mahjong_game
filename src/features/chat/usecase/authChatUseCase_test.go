package usecase

import (
	"context"
	"testing"
	"time"

	"main/features/chat/model/interface/mocks"

	"github.com/stretchr/testify/mock"
	"gopkg.in/go-playground/assert.v1"
	"gorm.io/gorm"

	"main/utils/db/mysql"
)

// TestAuthChatUseCase_Auth 함수는 AuthChatUseCase 의 Auth 메서드를 테스트합니다.
// 테이블 기반 테스트를 사용하여 여러 시나리오를 테스트합니다.
// given-when-then 패턴을 사용하여 테스트를 작성합니다.
// 구조체 필드 :  *mysql.Users, *mysql.Chats, wantErr error
// 반환값 : string, error
// 테스트 케이스:
// 1. userID 가 존재해서 정상적으로 응답한 경우
// 테스트 경로: src/features/chats/usecase/authChatUseCase_test.go

func TestAuthChatUseCase_Auth(t *testing.T) {
	tests := []struct {
		name    string
		userDTO *mysql.Users
		chatDTO *mysql.Chats
		wantErr error
	}{
		{
			name: "userID 가 존재해서 정상적으로 응답한 경우",
			userDTO: &mysql.Users{
				Model: gorm.Model{
					ID: 1,
				},
				Name:  "test",
				Email: "ryan@gmail.com",
			},
			chatDTO: &mysql.Chats{
				Model: gorm.Model{
					ID: 1,
				},
				UserID: 1,
				Secret: "test",
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			mockAuthChatRepository := new(mocks.IAuthChatRepository)
			mockAuthChatRepository.On("FindOneUserInfo", mock.Anything, mock.Anything).Return(tt.userDTO, tt.wantErr) //mock
			mockAuthChatRepository.On("InsertOneChat", mock.Anything, mock.Anything).Return(tt.wantErr)               //mock
			us := NewAuthChatUseCase(mockAuthChatRepository, 8*time.Second)
			// when
			_, err := us.Auth(context.TODO(), 1)
			// then
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
