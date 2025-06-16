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
	"github.com/kadekchresna/payroll/internal/v1/compensation/dto"

	usecase_interface "github.com/kadekchresna/payroll/internal/v1/compensation/usecase/interface/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
)

func TestOvertimeHandler_Create(t *testing.T) {

	body := OvertimeRequest{
		Date:  "2025-06-01",
		Hours: 2,
	}

	parsedDate, _ := time.Parse("2006-01-02", body.Date)

	type args struct {
		req OvertimeRequest
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		beforeFunc func(uc *usecase_interface.MockIOvertimeUsecase) *OvertimeHandler
	}{
		{
			name: "Create-Success",
			args: args{
				OvertimeRequest{
					Date:  "2025-06-01",
					Hours: 2,
				},
			},
			wantErr: false,
			beforeFunc: func(uc *usecase_interface.MockIOvertimeUsecase) *OvertimeHandler {
				overtime := &dto.CreateOvertimeRequest{
					UserID:     0,
					EmployeeID: 0,
					Date:       parsedDate,
					Hours:      body.Hours,
				}
				uc.EXPECT().CreateOvertime(mock.Anything, overtime).Return(nil).Once()

				return &OvertimeHandler{
					uc: uc,
				}
			},
		},
		{
			name: "Create-Failed",
			args: args{
				OvertimeRequest{
					Date:  "2025-06-01",
					Hours: 2,
				},
			},
			wantErr: true,
			beforeFunc: func(uc *usecase_interface.MockIOvertimeUsecase) *OvertimeHandler {
				overtime := &dto.CreateOvertimeRequest{
					UserID:     0,
					EmployeeID: 0,
					Date:       parsedDate,
					Hours:      body.Hours,
				}
				uc.EXPECT().CreateOvertime(mock.Anything, overtime).Return(errors.New("FATAL ERROR")).Once()

				return &OvertimeHandler{
					uc: uc,
				}
			},
		},
		{
			name: "Create-FailedErrorParseDate",
			args: args{
				OvertimeRequest{
					Date:  "2025-06-012",
					Hours: 2,
				},
			},
			wantErr: true,
			beforeFunc: func(uc *usecase_interface.MockIOvertimeUsecase) *OvertimeHandler {

				return &OvertimeHandler{
					uc: uc,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := new(usecase_interface.MockIOvertimeUsecase)
			e := echo.New()

			jsonBody, _ := json.Marshal(tt.args.req)
			req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(jsonBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)

			h := tt.beforeFunc(uc)
			if err := h.Create(c); (err != nil) != tt.wantErr {
				t.Errorf("OvertimeHandler.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewOvertimeHandler(t *testing.T) {
	type args struct {
		e      *echo.Group
		config config.Config
		uc     *usecase_interface.MockIOvertimeUsecase
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "NewOvertimeHandler-Success",
			args: args{
				e:      echo.New().Group("v1"),
				config: config.Config{},
				uc:     usecase_interface.NewMockIOvertimeUsecase(t),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			NewOvertimeHandler(tt.args.e, tt.args.config, tt.args.uc)
		})
	}
}
