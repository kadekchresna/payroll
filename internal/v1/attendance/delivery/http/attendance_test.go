package delivery_http

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kadekchresna/payroll/config"
	"github.com/kadekchresna/payroll/internal/v1/attendance/dto"

	usecase_interface "github.com/kadekchresna/payroll/internal/v1/attendance/usecase/interface/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
)

func TestAttendanceHandler_Create(t *testing.T) {
	req := AttendanceRequest{
		EmployeeID: 1,
		Date:       "2025-06-01",
	}

	parsedDate, _ := time.Parse("2006-01-02", req.Date)

	type args struct {
		body AttendanceRequest
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		beforeFunc func(uc *usecase_interface.MockIAttendanceUsecase) AttendanceHandler
	}{
		{
			name: "Create-Success",
			args: args{
				body: req,
			},
			wantErr: false,
			beforeFunc: func(uc *usecase_interface.MockIAttendanceUsecase) AttendanceHandler {
				uc.EXPECT().CreateAttendance(mock.Anything, &dto.CreateAttendanceRequest{
					EmployeeID: 0,
					Date:       parsedDate,
					UserID:     0,
				}).Return(nil).Once()
				return AttendanceHandler{
					uc: uc,
				}
			},
		},
		{
			name: "Create-Failed",
			args: args{
				body: req,
			},
			wantErr: true,
			beforeFunc: func(uc *usecase_interface.MockIAttendanceUsecase) AttendanceHandler {
				uc.EXPECT().CreateAttendance(mock.Anything, &dto.CreateAttendanceRequest{
					EmployeeID: 0,
					Date:       parsedDate,
					UserID:     0,
				}).Return(errors.New("FATAL ERROR")).Once()
				return AttendanceHandler{
					uc: uc,
				}
			},
		},
		{
			name: "Create-FailedParseDate",
			args: args{
				body: AttendanceRequest{
					Date: "2025-06-123",
				},
			},
			wantErr: true,
			beforeFunc: func(uc *usecase_interface.MockIAttendanceUsecase) AttendanceHandler {
				uc.EXPECT().CreateAttendance(mock.Anything, &dto.CreateAttendanceRequest{
					EmployeeID: 0,
					Date:       parsedDate,
					UserID:     0,
				}).Return(nil).Once()
				return AttendanceHandler{
					uc: uc,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := new(usecase_interface.MockIAttendanceUsecase)
			e := echo.New()

			jsonBody, _ := json.Marshal(tt.args.body)
			req := httptest.NewRequest(http.MethodPost, "/attendances", bytes.NewReader(jsonBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)

			h := tt.beforeFunc(uc)
			if err := h.Create(c); (err != nil) != tt.wantErr {
				t.Errorf("AttendanceHandler.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewAttendanceHandler(t *testing.T) {
	type args struct {
		e      *echo.Group
		config config.Config
		uc     *usecase_interface.MockIAttendanceUsecase
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "NewAttendanceHandler-Success",
			args: args{
				e:      echo.New().Group("v1"),
				config: config.Config{},
				uc:     usecase_interface.NewMockIAttendanceUsecase(t),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			NewAttendanceHandler(tt.args.e, tt.args.config, tt.args.uc)
		})
	}
}
