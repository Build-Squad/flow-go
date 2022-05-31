// Code generated by mockery v2.12.1. DO NOT EDIT.

package mock

import (
	flow "github.com/onflow/flow-go/model/flow"
	mock "github.com/stretchr/testify/mock"

	testing "testing"
)

// Finalizer is an autogenerated mock type for the Finalizer type
type Finalizer struct {
	mock.Mock
}

// MakeFinal provides a mock function with given fields: blockID
func (_m *Finalizer) MakeFinal(blockID flow.Identifier) error {
	ret := _m.Called(blockID)

	var r0 error
	if rf, ok := ret.Get(0).(func(flow.Identifier) error); ok {
		r0 = rf(blockID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MakeValid provides a mock function with given fields: blockID
func (_m *Finalizer) MakeValid(blockID flow.Identifier) error {
	ret := _m.Called(blockID)

	var r0 error
	if rf, ok := ret.Get(0).(func(flow.Identifier) error); ok {
		r0 = rf(blockID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewFinalizer creates a new instance of Finalizer. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewFinalizer(t testing.TB) *Finalizer {
	mock := &Finalizer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
