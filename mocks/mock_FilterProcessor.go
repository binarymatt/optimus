// Code generated by mockery v2.28.2. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	optimusv1 "github.com/binarymatt/optimus/gen/optimus/v1"
)

// MockFilterProcessor is an autogenerated mock type for the FilterProcessor type
type MockFilterProcessor struct {
	mock.Mock
}

type MockFilterProcessor_Expecter struct {
	mock *mock.Mock
}

func (_m *MockFilterProcessor) EXPECT() *MockFilterProcessor_Expecter {
	return &MockFilterProcessor_Expecter{mock: &_m.Mock}
}

// Process provides a mock function with given fields: _a0, _a1
func (_m *MockFilterProcessor) Process(_a0 context.Context, _a1 *optimusv1.LogEvent) (*optimusv1.LogEvent, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *optimusv1.LogEvent
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *optimusv1.LogEvent) (*optimusv1.LogEvent, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *optimusv1.LogEvent) *optimusv1.LogEvent); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*optimusv1.LogEvent)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *optimusv1.LogEvent) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockFilterProcessor_Process_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Process'
type MockFilterProcessor_Process_Call struct {
	*mock.Call
}

// Process is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 *optimusv1.LogEvent
func (_e *MockFilterProcessor_Expecter) Process(_a0 interface{}, _a1 interface{}) *MockFilterProcessor_Process_Call {
	return &MockFilterProcessor_Process_Call{Call: _e.mock.On("Process", _a0, _a1)}
}

func (_c *MockFilterProcessor_Process_Call) Run(run func(_a0 context.Context, _a1 *optimusv1.LogEvent)) *MockFilterProcessor_Process_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*optimusv1.LogEvent))
	})
	return _c
}

func (_c *MockFilterProcessor_Process_Call) Return(_a0 *optimusv1.LogEvent, _a1 error) *MockFilterProcessor_Process_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockFilterProcessor_Process_Call) RunAndReturn(run func(context.Context, *optimusv1.LogEvent) (*optimusv1.LogEvent, error)) *MockFilterProcessor_Process_Call {
	_c.Call.Return(run)
	return _c
}

// Setup provides a mock function with given fields:
func (_m *MockFilterProcessor) Setup() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockFilterProcessor_Setup_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Setup'
type MockFilterProcessor_Setup_Call struct {
	*mock.Call
}

// Setup is a helper method to define mock.On call
func (_e *MockFilterProcessor_Expecter) Setup() *MockFilterProcessor_Setup_Call {
	return &MockFilterProcessor_Setup_Call{Call: _e.mock.On("Setup")}
}

func (_c *MockFilterProcessor_Setup_Call) Run(run func()) *MockFilterProcessor_Setup_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockFilterProcessor_Setup_Call) Return(_a0 error) *MockFilterProcessor_Setup_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockFilterProcessor_Setup_Call) RunAndReturn(run func() error) *MockFilterProcessor_Setup_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockFilterProcessor interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockFilterProcessor creates a new instance of MockFilterProcessor. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockFilterProcessor(t mockConstructorTestingTNewMockFilterProcessor) *MockFilterProcessor {
	mock := &MockFilterProcessor{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
