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

func TestReimbursementHandler_Create(t *testing.T) {

	body := ReimbursementRequest{
		Date:        "2025-06-01",
		Amount:      10000,
		Description: "desc",
	}

	parsedDate, _ := time.Parse("2006-01-02", body.Date)

	type args struct {
		req ReimbursementRequest
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		beforeFunc func(uc *usecase_interface.MockIReimbursementUsecase) *ReimbursementHandler
	}{
		{
			name: "Create-Success",
			args: args{
				req: body,
			},
			wantErr: false,
			beforeFunc: func(uc *usecase_interface.MockIReimbursementUsecase) *ReimbursementHandler {

				uc.EXPECT().CreateReimbursement(mock.Anything, &dto.CreateReimbursementRequest{
					UserID:      0,
					EmployeeID:  0,
					Date:        parsedDate,
					Description: body.Description,
					Amount:      body.Amount,
				}).Return(nil).Once()

				return &ReimbursementHandler{
					uc: uc,
				}
			},
		},
		{
			name: "Create-Failed",
			args: args{
				req: body,
			},
			wantErr: true,
			beforeFunc: func(uc *usecase_interface.MockIReimbursementUsecase) *ReimbursementHandler {

				uc.EXPECT().CreateReimbursement(mock.Anything, &dto.CreateReimbursementRequest{
					UserID:      0,
					EmployeeID:  0,
					Date:        parsedDate,
					Description: body.Description,
					Amount:      body.Amount,
				}).Return(errors.New("FATAL")).Once()

				return &ReimbursementHandler{
					uc: uc,
				}
			},
		},
		{
			name: "Create-FailedParseDate",
			args: args{
				req: func() ReimbursementRequest {
					return ReimbursementRequest{
						Date: "2025-06-012",
					}
				}(),
			},
			wantErr: true,
			beforeFunc: func(uc *usecase_interface.MockIReimbursementUsecase) *ReimbursementHandler {

				uc.EXPECT().CreateReimbursement(mock.Anything, &dto.CreateReimbursementRequest{
					UserID:      0,
					EmployeeID:  0,
					Date:        parsedDate,
					Description: body.Description,
					Amount:      body.Amount,
				}).Return(errors.New("FATAL")).Once()

				return &ReimbursementHandler{
					uc: uc,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := new(usecase_interface.MockIReimbursementUsecase)
			e := echo.New()

			jsonBody, _ := json.Marshal(tt.args.req)
			req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(jsonBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)

			h := tt.beforeFunc(uc)
			if err := h.Create(c); (err != nil) != tt.wantErr {
				t.Errorf("ReimbursementHandler.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewReimbursementHandler(t *testing.T) {
	type args struct {
		e      *echo.Group
		config config.Config
		uc     *usecase_interface.MockIReimbursementUsecase
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "NewReimbursementHandler-Success",
			args: args{
				e:      echo.New().Group("v1"),
				config: config.Config{},
				uc:     usecase_interface.NewMockIReimbursementUsecase(t),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			NewReimbursementHandler(tt.args.e, tt.args.config, tt.args.uc)
		})
	}
}
