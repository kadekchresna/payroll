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
	"github.com/kadekchresna/payroll/helper/jwt"
	"github.com/kadekchresna/payroll/internal/v1/attendance/dto"

	usecase_interface "github.com/kadekchresna/payroll/internal/v1/attendance/usecase/interface/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
)

func TestAttendancePeriodHandler_Create(t *testing.T) {

	req := AttendancePeriodRequest{
		PeriodStart: "2025-06-01",
		PeriodEnd:   "2025-06-30",
	}

	parsedPeriodStart, _ := time.Parse("2006-01-02", req.PeriodStart)
	parsedPeriodEnd, _ := time.Parse("2006-01-02", req.PeriodEnd)
	type args struct {
		body AttendancePeriodRequest
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		beforeFunc func(uc *usecase_interface.MockIAttendancePeriodUsecase) AttendancePeriodHandler
	}{
		{
			name: "Create-Success",
			args: args{
				body: req,
			},
			wantErr: false,
			beforeFunc: func(uc *usecase_interface.MockIAttendancePeriodUsecase) AttendancePeriodHandler {
				uc.EXPECT().CreateAttendancePeriod(mock.Anything, &dto.CreateAttendancePeriodRequest{
					PeriodStart: parsedPeriodStart,
					PeriodEnd:   parsedPeriodEnd,
				}).Return(nil).Once()
				return AttendancePeriodHandler{
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
			beforeFunc: func(uc *usecase_interface.MockIAttendancePeriodUsecase) AttendancePeriodHandler {
				uc.EXPECT().CreateAttendancePeriod(mock.Anything, &dto.CreateAttendancePeriodRequest{
					PeriodStart: parsedPeriodStart,
					PeriodEnd:   parsedPeriodEnd,
				}).Return(errors.New("FATAL ERROR")).Once()
				return AttendancePeriodHandler{
					uc: uc,
				}
			},
		},
		{
			name: "Create-FailedParsedPeriodStart",
			args: args{
				body: AttendancePeriodRequest{
					PeriodEnd:   req.PeriodEnd,
					PeriodStart: "2025-06-123",
				},
			},
			wantErr: true,
			beforeFunc: func(uc *usecase_interface.MockIAttendancePeriodUsecase) AttendancePeriodHandler {
				uc.EXPECT().CreateAttendancePeriod(mock.Anything, &dto.CreateAttendancePeriodRequest{
					PeriodStart: parsedPeriodStart,
					PeriodEnd:   parsedPeriodEnd,
				}).Return(nil).Once()
				return AttendancePeriodHandler{
					uc: uc,
				}
			},
		},
		{
			name: "Create-FailedParsedPeriodEnd",
			args: args{
				body: AttendancePeriodRequest{
					PeriodStart: req.PeriodStart,
					PeriodEnd:   "2025-06-123",
				},
			},
			wantErr: true,
			beforeFunc: func(uc *usecase_interface.MockIAttendancePeriodUsecase) AttendancePeriodHandler {
				uc.EXPECT().CreateAttendancePeriod(mock.Anything, &dto.CreateAttendancePeriodRequest{
					PeriodStart: parsedPeriodStart,
					PeriodEnd:   parsedPeriodEnd,
				}).Return(nil).Once()
				return AttendancePeriodHandler{
					uc: uc,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := new(usecase_interface.MockIAttendancePeriodUsecase)
			e := echo.New()

			jsonBody, _ := json.Marshal(tt.args.body)
			req := httptest.NewRequest(http.MethodPost, "/attendances-period", bytes.NewReader(jsonBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)

			c.Set(jwt.USER_ROLE_KEY, "admin")

			h := tt.beforeFunc(uc)
			if err := h.Create(c); (err != nil) != tt.wantErr {
				t.Errorf("AttendancePeriodHandler.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewAttendancePeriodHandler(t *testing.T) {
	type args struct {
		e      *echo.Group
		config config.Config
		uc     *usecase_interface.MockIAttendancePeriodUsecase
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "NewAttendancePeriodHandler-Success",
			args: args{
				e:      echo.New().Group("v1"),
				config: config.Config{},
				uc:     usecase_interface.NewMockIAttendancePeriodUsecase(t),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			NewAttendancePeriodHandler(tt.args.e, tt.args.config, tt.args.uc)
		})
	}
}
