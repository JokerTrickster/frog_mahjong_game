// Code generated by mockery v2.42.2. DO NOT EDIT.

package mocks

import (
	context "context"
	entity "main/features/auth/model/entity"

	mock "github.com/stretchr/testify/mock"

	mysql "main/utils/db/mysql"
)

// IGoogleOauthCallbackAuthRepository is an autogenerated mock type for the IGoogleOauthCallbackAuthRepository type
type IGoogleOauthCallbackAuthRepository struct {
	mock.Mock
}

// CreateUser provides a mock function with given fields: ctx, user
func (_m *IGoogleOauthCallbackAuthRepository) CreateUser(ctx context.Context, user *mysql.Users) (*mysql.Users, error) {
	ret := _m.Called(ctx, user)

	if len(ret) == 0 {
		panic("no return value specified for CreateUser")
	}

	var r0 *mysql.Users
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *mysql.Users) (*mysql.Users, error)); ok {
		return rf(ctx, user)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *mysql.Users) *mysql.Users); ok {
		r0 = rf(ctx, user)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mysql.Users)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *mysql.Users) error); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteToken provides a mock function with given fields: ctx, uID
func (_m *IGoogleOauthCallbackAuthRepository) DeleteToken(ctx context.Context, uID uint) error {
	ret := _m.Called(ctx, uID)

	if len(ret) == 0 {
		panic("no return value specified for DeleteToken")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint) error); ok {
		r0 = rf(ctx, uID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindOneAndUpdateUser provides a mock function with given fields: ctx, googleOauthCallbackSQLQuery
func (_m *IGoogleOauthCallbackAuthRepository) FindOneAndUpdateUser(ctx context.Context, googleOauthCallbackSQLQuery *entity.GoogleOauthCallbackSQLQuery) (*mysql.Users, error) {
	ret := _m.Called(ctx, googleOauthCallbackSQLQuery)

	if len(ret) == 0 {
		panic("no return value specified for FindOneAndUpdateUser")
	}

	var r0 *mysql.Users
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *entity.GoogleOauthCallbackSQLQuery) (*mysql.Users, error)); ok {
		return rf(ctx, googleOauthCallbackSQLQuery)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *entity.GoogleOauthCallbackSQLQuery) *mysql.Users); ok {
		r0 = rf(ctx, googleOauthCallbackSQLQuery)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mysql.Users)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *entity.GoogleOauthCallbackSQLQuery) error); ok {
		r1 = rf(ctx, googleOauthCallbackSQLQuery)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SaveToken provides a mock function with given fields: ctx, uID, accessToken, refreshToken, refreshTknExpiredAt
func (_m *IGoogleOauthCallbackAuthRepository) SaveToken(ctx context.Context, uID uint, accessToken string, refreshToken string, refreshTknExpiredAt int64) error {
	ret := _m.Called(ctx, uID, accessToken, refreshToken, refreshTknExpiredAt)

	if len(ret) == 0 {
		panic("no return value specified for SaveToken")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint, string, string, int64) error); ok {
		r0 = rf(ctx, uID, accessToken, refreshToken, refreshTknExpiredAt)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewIGoogleOauthCallbackAuthRepository creates a new instance of IGoogleOauthCallbackAuthRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIGoogleOauthCallbackAuthRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *IGoogleOauthCallbackAuthRepository {
	mock := &IGoogleOauthCallbackAuthRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
