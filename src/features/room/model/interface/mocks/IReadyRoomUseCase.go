// Code generated by mockery v2.40.1. DO NOT EDIT.

package mocks

import (
	context "context"
	request "main/features/room/model/request"

	mock "github.com/stretchr/testify/mock"
)

// IReadyRoomUseCase is an autogenerated mock type for the IReadyRoomUseCase type
type IReadyRoomUseCase struct {
	mock.Mock
}

// Ready provides a mock function with given fields: c, uID, req
func (_m *IReadyRoomUseCase) Ready(c context.Context, uID uint, req *request.ReqReady) error {
	ret := _m.Called(c, uID, req)

	if len(ret) == 0 {
		panic("no return value specified for Ready")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint, *request.ReqReady) error); ok {
		r0 = rf(c, uID, req)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewIReadyRoomUseCase creates a new instance of IReadyRoomUseCase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIReadyRoomUseCase(t interface {
	mock.TestingT
	Cleanup(func())
}) *IReadyRoomUseCase {
	mock := &IReadyRoomUseCase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}