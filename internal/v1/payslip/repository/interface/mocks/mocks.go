// Code generated by mockery; DO NOT EDIT.
// github.com/vektra/mockery
// template: testify

package repository_interface_test

import (
	"context"

	"github.com/kadekchresna/payroll/internal/v1/payslip/model"
	mock "github.com/stretchr/testify/mock"
)

// NewMockIPayslipRepository creates a new instance of MockIPayslipRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockIPayslipRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockIPayslipRepository {
	mock := &MockIPayslipRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

// MockIPayslipRepository is an autogenerated mock type for the IPayslipRepository type
type MockIPayslipRepository struct {
	mock.Mock
}

type MockIPayslipRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *MockIPayslipRepository) EXPECT() *MockIPayslipRepository_Expecter {
	return &MockIPayslipRepository_Expecter{mock: &_m.Mock}
}

// Create provides a mock function for the type MockIPayslipRepository
func (_mock *MockIPayslipRepository) Create(ctx context.Context, p *model.Payslip) (int, error) {
	ret := _mock.Called(ctx, p)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 int
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, *model.Payslip) (int, error)); ok {
		return returnFunc(ctx, p)
	}
	if returnFunc, ok := ret.Get(0).(func(context.Context, *model.Payslip) int); ok {
		r0 = returnFunc(ctx, p)
	} else {
		r0 = ret.Get(0).(int)
	}
	if returnFunc, ok := ret.Get(1).(func(context.Context, *model.Payslip) error); ok {
		r1 = returnFunc(ctx, p)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockIPayslipRepository_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type MockIPayslipRepository_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - ctx context.Context
//   - p *model.Payslip
func (_e *MockIPayslipRepository_Expecter) Create(ctx interface{}, p interface{}) *MockIPayslipRepository_Create_Call {
	return &MockIPayslipRepository_Create_Call{Call: _e.mock.On("Create", ctx, p)}
}

func (_c *MockIPayslipRepository_Create_Call) Run(run func(ctx context.Context, p *model.Payslip)) *MockIPayslipRepository_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		var arg0 context.Context
		if args[0] != nil {
			arg0 = args[0].(context.Context)
		}
		var arg1 *model.Payslip
		if args[1] != nil {
			arg1 = args[1].(*model.Payslip)
		}
		run(
			arg0,
			arg1,
		)
	})
	return _c
}

func (_c *MockIPayslipRepository_Create_Call) Return(n int, err error) *MockIPayslipRepository_Create_Call {
	_c.Call.Return(n, err)
	return _c
}

func (_c *MockIPayslipRepository_Create_Call) RunAndReturn(run func(ctx context.Context, p *model.Payslip) (int, error)) *MockIPayslipRepository_Create_Call {
	_c.Call.Return(run)
	return _c
}

// GetByID provides a mock function for the type MockIPayslipRepository
func (_mock *MockIPayslipRepository) GetByID(ctx context.Context, id int) (*model.Payslip, error) {
	ret := _mock.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetByID")
	}

	var r0 *model.Payslip
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, int) (*model.Payslip, error)); ok {
		return returnFunc(ctx, id)
	}
	if returnFunc, ok := ret.Get(0).(func(context.Context, int) *model.Payslip); ok {
		r0 = returnFunc(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Payslip)
		}
	}
	if returnFunc, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = returnFunc(ctx, id)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockIPayslipRepository_GetByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByID'
type MockIPayslipRepository_GetByID_Call struct {
	*mock.Call
}

// GetByID is a helper method to define mock.On call
//   - ctx context.Context
//   - id int
func (_e *MockIPayslipRepository_Expecter) GetByID(ctx interface{}, id interface{}) *MockIPayslipRepository_GetByID_Call {
	return &MockIPayslipRepository_GetByID_Call{Call: _e.mock.On("GetByID", ctx, id)}
}

func (_c *MockIPayslipRepository_GetByID_Call) Run(run func(ctx context.Context, id int)) *MockIPayslipRepository_GetByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		var arg0 context.Context
		if args[0] != nil {
			arg0 = args[0].(context.Context)
		}
		var arg1 int
		if args[1] != nil {
			arg1 = args[1].(int)
		}
		run(
			arg0,
			arg1,
		)
	})
	return _c
}

func (_c *MockIPayslipRepository_GetByID_Call) Return(payslip *model.Payslip, err error) *MockIPayslipRepository_GetByID_Call {
	_c.Call.Return(payslip, err)
	return _c
}

func (_c *MockIPayslipRepository_GetByID_Call) RunAndReturn(run func(ctx context.Context, id int) (*model.Payslip, error)) *MockIPayslipRepository_GetByID_Call {
	_c.Call.Return(run)
	return _c
}

// GetTotalTakeHomePayAllEmployees provides a mock function for the type MockIPayslipRepository
func (_mock *MockIPayslipRepository) GetTotalTakeHomePayAllEmployees(ctx context.Context) (float64, error) {
	ret := _mock.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetTotalTakeHomePayAllEmployees")
	}

	var r0 float64
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(context.Context) (float64, error)); ok {
		return returnFunc(ctx)
	}
	if returnFunc, ok := ret.Get(0).(func(context.Context) float64); ok {
		r0 = returnFunc(ctx)
	} else {
		r0 = ret.Get(0).(float64)
	}
	if returnFunc, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = returnFunc(ctx)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockIPayslipRepository_GetTotalTakeHomePayAllEmployees_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetTotalTakeHomePayAllEmployees'
type MockIPayslipRepository_GetTotalTakeHomePayAllEmployees_Call struct {
	*mock.Call
}

// GetTotalTakeHomePayAllEmployees is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockIPayslipRepository_Expecter) GetTotalTakeHomePayAllEmployees(ctx interface{}) *MockIPayslipRepository_GetTotalTakeHomePayAllEmployees_Call {
	return &MockIPayslipRepository_GetTotalTakeHomePayAllEmployees_Call{Call: _e.mock.On("GetTotalTakeHomePayAllEmployees", ctx)}
}

func (_c *MockIPayslipRepository_GetTotalTakeHomePayAllEmployees_Call) Run(run func(ctx context.Context)) *MockIPayslipRepository_GetTotalTakeHomePayAllEmployees_Call {
	_c.Call.Run(func(args mock.Arguments) {
		var arg0 context.Context
		if args[0] != nil {
			arg0 = args[0].(context.Context)
		}
		run(
			arg0,
		)
	})
	return _c
}

func (_c *MockIPayslipRepository_GetTotalTakeHomePayAllEmployees_Call) Return(f float64, err error) *MockIPayslipRepository_GetTotalTakeHomePayAllEmployees_Call {
	_c.Call.Return(f, err)
	return _c
}

func (_c *MockIPayslipRepository_GetTotalTakeHomePayAllEmployees_Call) RunAndReturn(run func(ctx context.Context) (float64, error)) *MockIPayslipRepository_GetTotalTakeHomePayAllEmployees_Call {
	_c.Call.Return(run)
	return _c
}

// GetTotalTakeHomePayPerEmployee provides a mock function for the type MockIPayslipRepository
func (_mock *MockIPayslipRepository) GetTotalTakeHomePayPerEmployee(ctx context.Context) ([]model.TotalTakeHomePayPerEmployee, error) {
	ret := _mock.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetTotalTakeHomePayPerEmployee")
	}

	var r0 []model.TotalTakeHomePayPerEmployee
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(context.Context) ([]model.TotalTakeHomePayPerEmployee, error)); ok {
		return returnFunc(ctx)
	}
	if returnFunc, ok := ret.Get(0).(func(context.Context) []model.TotalTakeHomePayPerEmployee); ok {
		r0 = returnFunc(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.TotalTakeHomePayPerEmployee)
		}
	}
	if returnFunc, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = returnFunc(ctx)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockIPayslipRepository_GetTotalTakeHomePayPerEmployee_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetTotalTakeHomePayPerEmployee'
type MockIPayslipRepository_GetTotalTakeHomePayPerEmployee_Call struct {
	*mock.Call
}

// GetTotalTakeHomePayPerEmployee is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockIPayslipRepository_Expecter) GetTotalTakeHomePayPerEmployee(ctx interface{}) *MockIPayslipRepository_GetTotalTakeHomePayPerEmployee_Call {
	return &MockIPayslipRepository_GetTotalTakeHomePayPerEmployee_Call{Call: _e.mock.On("GetTotalTakeHomePayPerEmployee", ctx)}
}

func (_c *MockIPayslipRepository_GetTotalTakeHomePayPerEmployee_Call) Run(run func(ctx context.Context)) *MockIPayslipRepository_GetTotalTakeHomePayPerEmployee_Call {
	_c.Call.Run(func(args mock.Arguments) {
		var arg0 context.Context
		if args[0] != nil {
			arg0 = args[0].(context.Context)
		}
		run(
			arg0,
		)
	})
	return _c
}

func (_c *MockIPayslipRepository_GetTotalTakeHomePayPerEmployee_Call) Return(totalTakeHomePayPerEmployees []model.TotalTakeHomePayPerEmployee, err error) *MockIPayslipRepository_GetTotalTakeHomePayPerEmployee_Call {
	_c.Call.Return(totalTakeHomePayPerEmployees, err)
	return _c
}

func (_c *MockIPayslipRepository_GetTotalTakeHomePayPerEmployee_Call) RunAndReturn(run func(ctx context.Context) ([]model.TotalTakeHomePayPerEmployee, error)) *MockIPayslipRepository_GetTotalTakeHomePayPerEmployee_Call {
	_c.Call.Return(run)
	return _c
}
