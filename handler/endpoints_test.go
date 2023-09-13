package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/SawitProRecruitment/UserService/handler"
	"github.com/SawitProRecruitment/UserService/handler/mocks"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestPostRegister(t *testing.T) {
	// Initialize the Echo framework
	e := echo.New()

	t.Run("Bad Request - Failed to Validate", func(t *testing.T) {
		mockUserService := new(mocks.UserServiceInterface)

		reqBody := map[string]interface{}{
			"PhoneNumber": "",
		}
		reqBodyBytes, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(reqBodyBytes))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockUserService.On("RegisterNewUser", mock.Anything, mock.Anything).Return(nil, nil)
		s := &handler.Server{UserService: mockUserService}

		err := s.PostRegister(c)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusBadRequest, err.(*echo.HTTPError).Code)
	})

	t.Run("Success", func(t *testing.T) {
		mockUserService := new(mocks.UserServiceInterface)

		reqBody := map[string]interface{}{
			"PhoneNumber": "1234567890",
		}
		reqBodyBytes, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(reqBodyBytes))
		req.Header.Set("Content-Type", "application/json") // set header
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockUserService.On("FetchUserByPhoneNumber", mock.Anything, mock.Anything).Return(nil, nil)
		mockUserService.On("RegisterNewUser", mock.Anything, mock.Anything).Return(nil)
		s := &handler.Server{UserService: mockUserService}

		err := s.PostRegister(c)
		assert.Nil(t, err)                              // Check that the error is nil
		assert.Equal(t, http.StatusNoContent, rec.Code) // Check that the status code is 204
	})

	t.Run("FetchUserByPhoneNumber err", func(t *testing.T) {
		mockUserService := new(mocks.UserServiceInterface)

		reqBody := map[string]interface{}{
			"PhoneNumber": "1234567890",
		}
		reqBodyBytes, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(reqBodyBytes))
		req.Header.Set("Content-Type", "application/json") // set header
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockUserService.On("FetchUserByPhoneNumber", mock.Anything, mock.Anything).Return(nil, errors.New("simulate err"))
		s := &handler.Server{UserService: mockUserService}

		err := s.PostRegister(c)
		httpError, _ := err.(*echo.HTTPError)

		assert.Equal(t, http.StatusInternalServerError, httpError.Code) // Check that the status code is 204
	})

	t.Run("Phone number already exists", func(t *testing.T) {
		mockUserService := new(mocks.UserServiceInterface)

		reqBody := map[string]interface{}{
			"PhoneNumber": "1234567890",
		}
		reqBodyBytes, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(reqBodyBytes))
		req.Header.Set("Content-Type", "application/json") // set header
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockUserService.On("FetchUserByPhoneNumber", mock.Anything, mock.Anything).Return(&repository.UserModel{
			ID:          11,
			PhoneNumber: "111",
			FullName:    "11",
			Password:    "11",
			SaltKey:     "111",
			CreatedAt:   time.Time{},
			UpdatedAt:   time.Time{},
		}, nil)
		s := &handler.Server{UserService: mockUserService}

		err := s.PostRegister(c)
		httpError, _ := err.(*echo.HTTPError)

		assert.Equal(t, http.StatusConflict, httpError.Code) // Check that the status code is 204
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		mockUserService := new(mocks.UserServiceInterface)

		reqBody := map[string]interface{}{
			"PhoneNumber": "1111",
		}
		reqBodyBytes, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(reqBodyBytes))
		req.Header.Set("Content-Type", "application/json") // set header
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockUserService.On("FetchUserByPhoneNumber", mock.Anything, mock.Anything).Return(nil, nil)
		mockUserService.On("RegisterNewUser", mock.Anything, mock.Anything).Return(errors.New("simulate err"))
		s := &handler.Server{UserService: mockUserService}
		err := s.PostRegister(c)

		httpError, _ := err.(*echo.HTTPError)

		assert.Equal(t, http.StatusInternalServerError, httpError.Code) // Check that the status code is 500

	})

}
