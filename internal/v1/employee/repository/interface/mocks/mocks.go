// Code generated by mockery; DO NOT EDIT.
// github.com/vektra/mockery
// template: testify

package repository_interface_test

import (
	"context"

	"github.com/kadekchresna/payroll/internal/v1/employee/model"
	mock "github.com/stretchr/testify/mock"
)

// NewMockIEmployeeRepository creates a new instance of MockIEmployeeRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockIEmployeeRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockIEmployeeRepository {
	mock := &MockIEmployeeRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

// MockIEmployeeRepository is an autogenerated mock type for the IEmployeeRepository type
type MockIEmployeeRepository struct {
	mock.Mock
}

type MockIEmployeeRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *MockIEmployeeRepository) EXPECT() *MockIEmployeeRepository_Expecter {
	return &MockIEmployeeRepository_Expecter{mock: &_m.Mock}
}

// Create provides a mock function for the type MockIEmployeeRepository
func (_mock *MockIEmployeeRepository) Create(ctx context.Context, e *model.Employee) error {
	ret := _mock.Called(ctx, e)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, *model.Employee) error); ok {
		r0 = returnFunc(ctx, e)
	} else {
		r0 = ret.Error(0)
	}
	return r0
}

// MockIEmployeeRepository_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type MockIEmployeeRepository_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - ctx context.Context
//   - e *model.Employee
func (_e *MockIEmployeeRepository_Expecter) Create(ctx interface{}, e interface{}) *MockIEmployeeRepository_Create_Call {
	return &MockIEmployeeRepository_Create_Call{Call: _e.mock.On("Create", ctx, e)}
}

func (_c *MockIEmployeeRepository_Create_Call) Run(run func(ctx context.Context, e *model.Employee)) *MockIEmployeeRepository_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		var arg0 context.Context
		if args[0] != nil {
			arg0 = args[0].(context.Context)
		}
		var arg1 *model.Employee
		if args[1] != nil {
			arg1 = args[1].(*model.Employee)
		}
		run(
			arg0,
			arg1,
		)
	})
	return _c
}

func (_c *MockIEmployeeRepository_Create_Call) Return(err error) *MockIEmployeeRepository_Create_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *MockIEmployeeRepository_Create_Call) RunAndReturn(run func(ctx context.Context, e *model.Employee) error) *MockIEmployeeRepository_Create_Call {
	_c.Call.Return(run)
	return _c
}

// Delete provides a mock function for the type MockIEmployeeRepository
func (_mock *MockIEmployeeRepository) Delete(ctx context.Context, id int) error {
	ret := _mock.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, int) error); ok {
		r0 = returnFunc(ctx, id)
	} else {
		r0 = ret.Error(0)
	}
	return r0
}

// MockIEmployeeRepository_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type MockIEmployeeRepository_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - ctx context.Context
//   - id int
func (_e *MockIEmployeeRepository_Expecter) Delete(ctx interface{}, id interface{}) *MockIEmployeeRepository_Delete_Call {
	return &MockIEmployeeRepository_Delete_Call{Call: _e.mock.On("Delete", ctx, id)}
}

func (_c *MockIEmployeeRepository_Delete_Call) Run(run func(ctx context.Context, id int)) *MockIEmployeeRepository_Delete_Call {
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

func (_c *MockIEmployeeRepository_Delete_Call) Return(err error) *MockIEmployeeRepository_Delete_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *MockIEmployeeRepository_Delete_Call) RunAndReturn(run func(ctx context.Context, id int) error) *MockIEmployeeRepository_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// GetByID provides a mock function for the type MockIEmployeeRepository
func (_mock *MockIEmployeeRepository) GetByID(ctx context.Context, id int) (*model.Employee, error) {
	ret := _mock.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetByID")
	}

	var r0 *model.Employee
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, int) (*model.Employee, error)); ok {
		return returnFunc(ctx, id)
	}
	if returnFunc, ok := ret.Get(0).(func(context.Context, int) *model.Employee); ok {
		r0 = returnFunc(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Employee)
		}
	}
	if returnFunc, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = returnFunc(ctx, id)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockIEmployeeRepository_GetByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByID'
type MockIEmployeeRepository_GetByID_Call struct {
	*mock.Call
}

// GetByID is a helper method to define mock.On call
//   - ctx context.Context
//   - id int
func (_e *MockIEmployeeRepository_Expecter) GetByID(ctx interface{}, id interface{}) *MockIEmployeeRepository_GetByID_Call {
	return &MockIEmployeeRepository_GetByID_Call{Call: _e.mock.On("GetByID", ctx, id)}
}

func (_c *MockIEmployeeRepository_GetByID_Call) Run(run func(ctx context.Context, id int)) *MockIEmployeeRepository_GetByID_Call {
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

func (_c *MockIEmployeeRepository_GetByID_Call) Return(employee *model.Employee, err error) *MockIEmployeeRepository_GetByID_Call {
	_c.Call.Return(employee, err)
	return _c
}

func (_c *MockIEmployeeRepository_GetByID_Call) RunAndReturn(run func(ctx context.Context, id int) (*model.Employee, error)) *MockIEmployeeRepository_GetByID_Call {
	_c.Call.Return(run)
	return _c
}

// GetByUserID provides a mock function for the type MockIEmployeeRepository
func (_mock *MockIEmployeeRepository) GetByUserID(ctx context.Context, userID int) (*model.Employee, error) {
	ret := _mock.Called(ctx, userID)

	if len(ret) == 0 {
		panic("no return value specified for GetByUserID")
	}

	var r0 *model.Employee
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, int) (*model.Employee, error)); ok {
		return returnFunc(ctx, userID)
	}
	if returnFunc, ok := ret.Get(0).(func(context.Context, int) *model.Employee); ok {
		r0 = returnFunc(ctx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Employee)
		}
	}
	if returnFunc, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = returnFunc(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockIEmployeeRepository_GetByUserID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByUserID'
type MockIEmployeeRepository_GetByUserID_Call struct {
	*mock.Call
}

// GetByUserID is a helper method to define mock.On call
//   - ctx context.Context
//   - userID int
func (_e *MockIEmployeeRepository_Expecter) GetByUserID(ctx interface{}, userID interface{}) *MockIEmployeeRepository_GetByUserID_Call {
	return &MockIEmployeeRepository_GetByUserID_Call{Call: _e.mock.On("GetByUserID", ctx, userID)}
}

func (_c *MockIEmployeeRepository_GetByUserID_Call) Run(run func(ctx context.Context, userID int)) *MockIEmployeeRepository_GetByUserID_Call {
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

func (_c *MockIEmployeeRepository_GetByUserID_Call) Return(employee *model.Employee, err error) *MockIEmployeeRepository_GetByUserID_Call {
	_c.Call.Return(employee, err)
	return _c
}

func (_c *MockIEmployeeRepository_GetByUserID_Call) RunAndReturn(run func(ctx context.Context, userID int) (*model.Employee, error)) *MockIEmployeeRepository_GetByUserID_Call {
	_c.Call.Return(run)
	return _c
}

// ListAll provides a mock function for the type MockIEmployeeRepository
func (_mock *MockIEmployeeRepository) ListAll(ctx context.Context) ([]model.Employee, error) {
	ret := _mock.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for ListAll")
	}

	var r0 []model.Employee
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(context.Context) ([]model.Employee, error)); ok {
		return returnFunc(ctx)
	}
	if returnFunc, ok := ret.Get(0).(func(context.Context) []model.Employee); ok {
		r0 = returnFunc(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Employee)
		}
	}
	if returnFunc, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = returnFunc(ctx)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockIEmployeeRepository_ListAll_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListAll'
type MockIEmployeeRepository_ListAll_Call struct {
	*mock.Call
}

// ListAll is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockIEmployeeRepository_Expecter) ListAll(ctx interface{}) *MockIEmployeeRepository_ListAll_Call {
	return &MockIEmployeeRepository_ListAll_Call{Call: _e.mock.On("ListAll", ctx)}
}

func (_c *MockIEmployeeRepository_ListAll_Call) Run(run func(ctx context.Context)) *MockIEmployeeRepository_ListAll_Call {
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

func (_c *MockIEmployeeRepository_ListAll_Call) Return(employees []model.Employee, err error) *MockIEmployeeRepository_ListAll_Call {
	_c.Call.Return(employees, err)
	return _c
}

func (_c *MockIEmployeeRepository_ListAll_Call) RunAndReturn(run func(ctx context.Context) ([]model.Employee, error)) *MockIEmployeeRepository_ListAll_Call {
	_c.Call.Return(run)
	return _c
}

// Update provides a mock function for the type MockIEmployeeRepository
func (_mock *MockIEmployeeRepository) Update(ctx context.Context, e *model.Employee) error {
	ret := _mock.Called(ctx, e)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, *model.Employee) error); ok {
		r0 = returnFunc(ctx, e)
	} else {
		r0 = ret.Error(0)
	}
	return r0
}

// MockIEmployeeRepository_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type MockIEmployeeRepository_Update_Call struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - ctx context.Context
//   - e *model.Employee
func (_e *MockIEmployeeRepository_Expecter) Update(ctx interface{}, e interface{}) *MockIEmployeeRepository_Update_Call {
	return &MockIEmployeeRepository_Update_Call{Call: _e.mock.On("Update", ctx, e)}
}

func (_c *MockIEmployeeRepository_Update_Call) Run(run func(ctx context.Context, e *model.Employee)) *MockIEmployeeRepository_Update_Call {
	_c.Call.Run(func(args mock.Arguments) {
		var arg0 context.Context
		if args[0] != nil {
			arg0 = args[0].(context.Context)
		}
		var arg1 *model.Employee
		if args[1] != nil {
			arg1 = args[1].(*model.Employee)
		}
		run(
			arg0,
			arg1,
		)
	})
	return _c
}

func (_c *MockIEmployeeRepository_Update_Call) Return(err error) *MockIEmployeeRepository_Update_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *MockIEmployeeRepository_Update_Call) RunAndReturn(run func(ctx context.Context, e *model.Employee) error) *MockIEmployeeRepository_Update_Call {
	_c.Call.Return(run)
	return _c
}
