package usecase

import (
	"context"
	_errors "main/features/users/model/errors"
	"main/features/users/model/interface/mocks"
	"main/utils"

	"main/utils/db/mysql"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"gopkg.in/go-playground/assert.v1"
	"gorm.io/gorm"
)

// TestGetUsersUseCase_Get 함수는 GetUsersUseCase 의 Get 메서드를 테스트합니다.
// 테이블 기반 테스트를 사용하여 여러 시나리오를 테스트합니다.
// given-when-then 패턴을 사용하여 테스트를 작성합니다.
// 구조체 필드 : name string, userID int,userDTO mysql.Users, wantErr error
// 반환값 : response.ResGetUser, error
// 테스트 케이스:
// 1. userID 가 존재해서 정상적으로 응답한 경우
// 2. userID 가 존재하지 않아서 에러가 발생한 경우
// 테스트 경로: src/features/users/usecase/getUsersUseCase_test.go

func TestGetUsersUseCase_Get(t *testing.T) {
	tests := []struct {
		name    string
		userID  int
		userDTO mysql.Users
		wantErr error
	}{
		{
			name:   "userID 가 존재해서 정상적으로 응답한 경우",
			userID: 1,
			userDTO: mysql.Users{
				Model: gorm.Model{
					ID: 1,
				},
				Name:  "test",
				Email: "test@gmail.com",
			},
			wantErr: nil,
		},
		{
			name:    "userID 가 존재하지 않아서 에러가 발생한 경우",
			userID:  2,
			userDTO: mysql.Users{},
			wantErr: utils.ErrorMsg(context.TODO(), utils.ErrBadParameter, utils.Trace(), _errors.ErrUserNotFound.Error(), utils.ErrFromClient),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			mockGetUsersRepository := new(mocks.IGetUsersRepository)
			mockGetUsersRepository.On("FindOneUser", mock.Anything, mock.Anything).Return(tt.userDTO, tt.wantErr) //mock
			us := NewGetUsersUseCase(mockGetUsersRepository, 8*time.Second)

			// when
			_, err := us.Get(context.TODO(), tt.userID)

			// then
			assert.Equal(t, tt.wantErr, err)

		})
	}
}
