package delivery_http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kadekchresna/payroll/config"
	"github.com/kadekchresna/payroll/internal/v1/payroll/dto"
	usecase_interface "github.com/kadekchresna/payroll/internal/v1/payroll/usecase/interface/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
)

func TestPayrollHandler_GetPayrollSummary(t *testing.T) {

	tests := []struct {
		name       string
		wantErr    bool
		beforeFunc func(uc *usecase_interface.MockIPayrollUsecase) *PayrollHandler
	}{
		{
			name:    "GetPayrollSummary-Success",
			wantErr: false,
			beforeFunc: func(uc *usecase_interface.MockIPayrollUsecase) *PayrollHandler {

				uc.EXPECT().GetPayrollSummary(mock.Anything).Return(&dto.GetPayrollSummary{
					EmployeePayrollSummary: []dto.EmployeePayroll{
						{
							EmployeeID:       1,
							TotalTakeHomePay: 10000,
						},
					},
					TotalTakeHomePayAll: 1000,
				}, nil).Once()
				return &PayrollHandler{
					uc: uc,
				}
			},
		},
		{
			name:    "GetPayrollSummary-Failed",
			wantErr: true,
			beforeFunc: func(uc *usecase_interface.MockIPayrollUsecase) *PayrollHandler {

				uc.EXPECT().GetPayrollSummary(mock.Anything).Return(nil, errors.New("FATAL ERROR")).Once()
				return &PayrollHandler{
					uc: uc,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := new(usecase_interface.MockIPayrollUsecase)
			e := echo.New()

			req := httptest.NewRequest(http.MethodGet, "/payroll/summary", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)

			h := tt.beforeFunc(uc)
			if err := h.GetPayrollSummary(c); (err != nil) != tt.wantErr {
				t.Errorf("PayrollHandler.GetPayrollSummary() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPayrollHandler_GetPayroll(t *testing.T) {

	req := GetPayrollRequest{
		PayslipID: 1,
	}
	type args struct {
		req GetPayrollRequest
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		beforeFunc func(uc *usecase_interface.MockIPayrollUsecase) *PayrollHandler
	}{
		{
			name: "GetPayroll-Success",
			args: args{
				req: req,
			},
			wantErr: false,
			beforeFunc: func(uc *usecase_interface.MockIPayrollUsecase) *PayrollHandler {
				uc.EXPECT().GetPayrollByID(mock.Anything, &dto.GetEmployeePayrollRequest{
					EmployeeID: 0,
					PayslipID:  req.PayslipID,
				}).Return(&dto.GetEmployeePayrollResponse{}, nil).Once()
				return &PayrollHandler{
					uc: uc,
				}
			},
		},
		{
			name: "GetPayroll-Failed",
			args: args{
				req: req,
			},
			wantErr: true,
			beforeFunc: func(uc *usecase_interface.MockIPayrollUsecase) *PayrollHandler {
				uc.EXPECT().GetPayrollByID(mock.Anything, &dto.GetEmployeePayrollRequest{
					EmployeeID: 0,
					PayslipID:  req.PayslipID,
				}).Return(nil, errors.New("FATAL ERROR")).Once()
				return &PayrollHandler{
					uc: uc,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := new(usecase_interface.MockIPayrollUsecase)
			e := echo.New()

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/payroll?payslip_id=%d", tt.args.req.PayslipID), nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)

			h := tt.beforeFunc(uc)
			if err := h.GetPayroll(c); (err != nil) != tt.wantErr {
				t.Errorf("PayrollHandler.GetPayroll() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPayrollHandler_Create(t *testing.T) {

	type args struct {
		req PayrollRequest
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		beforeFunc func(uc *usecase_interface.MockIPayrollUsecase) *PayrollHandler
	}{
		{
			name: "Create-Success",
			args: args{
				req: PayrollRequest{
					AttendancePeriodID: 1,
				},
			},
			wantErr: false,
			beforeFunc: func(uc *usecase_interface.MockIPayrollUsecase) *PayrollHandler {
				uc.EXPECT().CreatePayroll(mock.Anything, dto.CreatePayrollRequest{
					AttendancePeriodID: 1,
				}).Return(nil).Once()
				return &PayrollHandler{
					uc: uc,
				}
			},
		},
		{
			name: "Create-Failed",
			args: args{
				req: PayrollRequest{
					AttendancePeriodID: 1,
				},
			},
			wantErr: true,
			beforeFunc: func(uc *usecase_interface.MockIPayrollUsecase) *PayrollHandler {
				uc.EXPECT().CreatePayroll(mock.Anything, dto.CreatePayrollRequest{
					AttendancePeriodID: 1,
				}).Return(errors.New("FATAL ERROR")).Once()
				return &PayrollHandler{
					uc: uc,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := new(usecase_interface.MockIPayrollUsecase)
			e := echo.New()

			jsonBody, _ := json.Marshal(tt.args.req)
			req := httptest.NewRequest(http.MethodPost, "/payroll", bytes.NewReader(jsonBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)

			h := tt.beforeFunc(uc)
			if err := h.Create(c); (err != nil) != tt.wantErr {
				t.Errorf("PayrollHandler.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewPayrollHandler(t *testing.T) {
	type args struct {
		e      *echo.Group
		config config.Config
		uc     *usecase_interface.MockIPayrollUsecase
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "NewPayrollHandler-Success",
			args: args{
				e:      echo.New().Group("v1"),
				config: config.Config{},
				uc:     usecase_interface.NewMockIPayrollUsecase(t),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			NewPayrollHandler(tt.args.e, tt.args.config, tt.args.uc)
		})
	}
}
