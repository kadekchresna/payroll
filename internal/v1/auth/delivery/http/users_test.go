package delivery_http

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kadekchresna/payroll/internal/v1/auth/dto"

	usecase_interface "github.com/kadekchresna/payroll/internal/v1/auth/usecase/interface/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
)

func TestUsersHandler_Register(t *testing.T) {
	reqBody := dto.CreateUserRequest{
		Username: "testuser",
		Password: "securepass",
	}

	tests := []struct {
		name       string
		wantErr    bool
		beforeFunc func(uc *usecase_interface.MockIUserUsecase) *UsersHandler
	}{
		{
			name: "Register-Success",
			beforeFunc: func(uc *usecase_interface.MockIUserUsecase) *UsersHandler {
				uc.EXPECT().Create(mock.Anything, reqBody).Return(nil).Once()
				return &UsersHandler{
					uc: uc,
				}
			},
		},
		{
			name: "Register-Failed",
			beforeFunc: func(uc *usecase_interface.MockIUserUsecase) *UsersHandler {
				uc.EXPECT().Create(mock.Anything, reqBody).Return(errors.New("FATAL ERROR")).Once()
				return &UsersHandler{
					uc: uc,
				}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := new(usecase_interface.MockIUserUsecase)
			e := echo.New()

			jsonBody, _ := json.Marshal(reqBody)
			req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(jsonBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)

			h := tt.beforeFunc(uc)
			if err := h.Register(c); (err != nil) != tt.wantErr {
				t.Errorf("UsersHandler.Register() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUsersHandler_Login(t *testing.T) {

	reqBody := dto.LoginUserRequest{
		Username: "testuser",
		Password: "securepass",
	}
	res := dto.LoginUserResponse{
		AccessToken:  "token",
		RefreshToken: "token",
	}

	tests := []struct {
		name       string
		wantErr    bool
		beforeFunc func(uc *usecase_interface.MockIUserUsecase) *UsersHandler
	}{
		{
			name: "Login-Success",
			beforeFunc: func(uc *usecase_interface.MockIUserUsecase) *UsersHandler {
				uc.EXPECT().Login(mock.Anything, reqBody).Return(&res, nil).Once()
				return &UsersHandler{
					uc: uc,
				}
			},
		},
		{
			name: "Login-Failed",
			beforeFunc: func(uc *usecase_interface.MockIUserUsecase) *UsersHandler {
				uc.EXPECT().Login(mock.Anything, reqBody).Return(nil, errors.New("FATAL ERROR")).Once()
				return &UsersHandler{
					uc: uc,
				}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := new(usecase_interface.MockIUserUsecase)
			e := echo.New()

			jsonBody, _ := json.Marshal(reqBody)
			req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(jsonBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)

			h := tt.beforeFunc(uc)
			if err := h.Login(c); (err != nil) != tt.wantErr {
				t.Errorf("UsersHandler.Login() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewUsersHandler(t *testing.T) {
	type args struct {
		e  *echo.Group
		uc *usecase_interface.MockIUserUsecase
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "NewUsersHandler-Success",
			args: args{
				e:  echo.New().Group("v1"),
				uc: new(usecase_interface.MockIUserUsecase),
			},
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			NewUsersHandler(tt.args.e, tt.args.uc)
		})
	}
}
