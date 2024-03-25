package usecase

import (
	"context"
	_errors "main/features/auth/model/errors"
	"main/features/auth/model/interface/mocks"
	"main/features/auth/model/request"
	"main/features/auth/model/response"
	"main/utils"
	"main/utils/db/mysql"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"gopkg.in/go-playground/assert.v1"
)

// Write code to test the signinAuthUseCase.go file

func TestSigninAuthUseCase_Signin(t *testing.T) {
	// table driven test
	tests := []struct {
		name string
		req  request.ReqSignin
		want response.ResSignin
		err  error
	}{
		{"success1", request.ReqSignin{Email: "ryan@breathings.co.kr", Password: "asdASD123"}, response.ResSignin{AccessToken: "test", RefreshToken: "test"}, nil},
		{"fail1", request.ReqSignin{Email: "ryan@breathings.co.kr", Password: "asdasdasd"}, response.ResSignin{}, utils.ErrorMsg(context.TODO(), utils.ErrUserNotExist, utils.Trace(), _errors.ErrUserNotFound.Error(), utils.ErrFromClient)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//given
			mockSigninAuthRepository := new(mocks.ISigninAuthRepository)
			if tt.err != nil {
				mockSigninAuthRepository.On("FindOneAndUpdateUser", mock.Anything, mock.Anything, mock.Anything).Return(mysql.Users{}, tt.err) //mock
			} else {
				mockSigninAuthRepository.On("FindOneAndUpdateUser", mock.Anything, mock.Anything, mock.Anything).Return(mysql.Users{Email: tt.req.Email, Password: tt.req.Password}, nil) //mock
			}
			us := NewSigninAuthUseCase(mockSigninAuthRepository, 8*time.Second)

			//when
			_, err := us.Signin(context.TODO(), &tt.req)
			//then
			assert.Equal(t, tt.err, err)
		})
	}
}
