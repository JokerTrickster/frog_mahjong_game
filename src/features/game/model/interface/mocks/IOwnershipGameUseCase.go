// Code generated by mockery v2.42.2. DO NOT EDIT.

package mocks

import (
	context "context"
	request "main/features/game/model/request"

	mock "github.com/stretchr/testify/mock"
)

// IOwnershipGameUseCase is an autogenerated mock type for the IOwnershipGameUseCase type
type IOwnershipGameUseCase struct {
	mock.Mock
}

// Ownership provides a mock function with given fields: c, req
func (_m *IOwnershipGameUseCase) Ownership(c context.Context, req *request.ReqOwnership) error {
	ret := _m.Called(c, req)

	if len(ret) == 0 {
		panic("no return value specified for Ownership")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *request.ReqOwnership) error); ok {
		r0 = rf(c, req)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewIOwnershipGameUseCase creates a new instance of IOwnershipGameUseCase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIOwnershipGameUseCase(t interface {
	mock.TestingT
	Cleanup(func())
}) *IOwnershipGameUseCase {
	mock := &IOwnershipGameUseCase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
