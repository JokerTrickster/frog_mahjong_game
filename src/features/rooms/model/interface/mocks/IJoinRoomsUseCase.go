// Code generated by mockery v2.42.2. DO NOT EDIT.

package mocks

import (
	context "context"
	request "main/features/rooms/model/request"

	mock "github.com/stretchr/testify/mock"

	response "main/features/rooms/model/response"
)

// IJoinRoomsUseCase is an autogenerated mock type for the IJoinRoomsUseCase type
type IJoinRoomsUseCase struct {
	mock.Mock
}

// Join provides a mock function with given fields: c, uID, email, req
func (_m *IJoinRoomsUseCase) Join(c context.Context, uID uint, email string, req *request.ReqJoin) (response.ResJoinRoom, error) {
	ret := _m.Called(c, uID, email, req)

	if len(ret) == 0 {
		panic("no return value specified for Join")
	}

	var r0 response.ResJoinRoom
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint, string, *request.ReqJoin) (response.ResJoinRoom, error)); ok {
		return rf(c, uID, email, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint, string, *request.ReqJoin) response.ResJoinRoom); ok {
		r0 = rf(c, uID, email, req)
	} else {
		r0 = ret.Get(0).(response.ResJoinRoom)
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint, string, *request.ReqJoin) error); ok {
		r1 = rf(c, uID, email, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewIJoinRoomsUseCase creates a new instance of IJoinRoomsUseCase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIJoinRoomsUseCase(t interface {
	mock.TestingT
	Cleanup(func())
}) *IJoinRoomsUseCase {
	mock := &IJoinRoomsUseCase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}