// Code generated by mockery v2.50.2. DO NOT EDIT.

//go:build !compile

package core

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockResourceRepo is an autogenerated mock type for the ResourceRepo type
type MockResourceRepo struct {
	mock.Mock
}

type MockResourceRepo_Expecter struct {
	mock *mock.Mock
}

func (_m *MockResourceRepo) EXPECT() *MockResourceRepo_Expecter {
	return &MockResourceRepo_Expecter{mock: &_m.Mock}
}

// GetCodeStyle provides a mock function with given fields: ctx, categories
func (_m *MockResourceRepo) GetCodeStyle(ctx context.Context, categories []string) ([]Rule, error) {
	ret := _m.Called(ctx, categories)

	if len(ret) == 0 {
		panic("no return value specified for GetCodeStyle")
	}

	var r0 []Rule
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, []string) ([]Rule, error)); ok {
		return rf(ctx, categories)
	}
	if rf, ok := ret.Get(0).(func(context.Context, []string) []Rule); ok {
		r0 = rf(ctx, categories)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]Rule)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, []string) error); ok {
		r1 = rf(ctx, categories)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockResourceRepo_GetCodeStyle_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetCodeStyle'
type MockResourceRepo_GetCodeStyle_Call struct {
	*mock.Call
}

// GetCodeStyle is a helper method to define mock.On call
//   - ctx context.Context
//   - categories []string
func (_e *MockResourceRepo_Expecter) GetCodeStyle(ctx interface{}, categories interface{}) *MockResourceRepo_GetCodeStyle_Call {
	return &MockResourceRepo_GetCodeStyle_Call{Call: _e.mock.On("GetCodeStyle", ctx, categories)}
}

func (_c *MockResourceRepo_GetCodeStyle_Call) Run(run func(ctx context.Context, categories []string)) *MockResourceRepo_GetCodeStyle_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]string))
	})
	return _c
}

func (_c *MockResourceRepo_GetCodeStyle_Call) Return(_a0 []Rule, _a1 error) *MockResourceRepo_GetCodeStyle_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockResourceRepo_GetCodeStyle_Call) RunAndReturn(run func(context.Context, []string) ([]Rule, error)) *MockResourceRepo_GetCodeStyle_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockResourceRepo creates a new instance of MockResourceRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockResourceRepo(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockResourceRepo {
	mock := &MockResourceRepo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
