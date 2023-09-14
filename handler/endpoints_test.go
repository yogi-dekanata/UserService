package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/SawitProRecruitment/UserService/commons"
	pwdMocks "github.com/SawitProRecruitment/UserService/commons/mocks"
	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/handler"
	"github.com/SawitProRecruitment/UserService/middleware"
	authMocks "github.com/SawitProRecruitment/UserService/middleware/mocks"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/SawitProRecruitment/UserService/repository/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetHealth(t *testing.T) {
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)

	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	s := &handler.Server{}

	if assert.NoError(t, s.GetHealth(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

	}
}

func TestPostRegister(t *testing.T) {
	e := echo.New()

	t.Run("Bad Request - Failed to Validate", func(t *testing.T) {
		mockRepo := new(mocks.RepositoryInterface)

		reqBody := map[string]interface{}{
			"PhoneNumber": "",
		}
		reqBodyBytes, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(reqBodyBytes))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		s := &handler.Server{Repository: mockRepo}

		err := s.PostRegister(c)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusBadRequest, err.(*echo.HTTPError).Code)
	})

	t.Run("User Already Exists", func(t *testing.T) {
		mockRepo := new(mocks.RepositoryInterface)
		reqBody := map[string]interface{}{"PhoneNumber": "+628222667727", "fullName": "LOLTOS", "password": "@Python12345@"}
		reqBodyBytes, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(reqBodyBytes))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		mockRepo.On("GetUser", mock.Anything, mock.Anything).Return(&repository.UserModel{ID: 111}, nil).Once()
		s := &handler.Server{Repository: mockRepo}
		err := s.PostRegister(c)
		if assert.Error(t, err) {
			httpErr, ok := err.(*echo.HTTPError)
			if assert.True(t, ok) {
				assert.Equal(t, http.StatusConflict, httpErr.Code)
			}
		}
		mockRepo.AssertExpectations(t)
	})

	t.Run("GetUser Err", func(t *testing.T) {
		mockRepo := new(mocks.RepositoryInterface)
		reqBody := map[string]interface{}{"PhoneNumber": "+628222667727", "fullName": "LOLTOS", "password": "@Python12345@"}
		reqBodyBytes, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(reqBodyBytes))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		mockRepo.On("GetUser", mock.Anything, mock.Anything).Return(nil, errors.New("simulate err"))
		s := &handler.Server{Repository: mockRepo}
		err := s.PostRegister(c)
		if assert.Error(t, err) {
			httpErr, ok := err.(*echo.HTTPError)
			if assert.True(t, ok) {
				assert.Equal(t, http.StatusInternalServerError, httpErr.Code)
			}
		}
		mockRepo.AssertExpectations(t)
	})

	t.Run("Register success", func(t *testing.T) {
		mockRepo := new(mocks.RepositoryInterface)
		mockPwd := new(pwdMocks.PasswordManagerInterface)
		reqBody := map[string]interface{}{"PhoneNumber": "+628222667727", "fullName": "LOLTOS", "password": "@Python12345@"}
		reqBodyBytes, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(reqBodyBytes))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockRepo.On("GetUser", mock.Anything, mock.Anything).Return(nil, nil)
		mockPwd.On("CreateSalt").Return("okCreate")
		mockPwd.On("GenerateHash", mock.Anything, mock.Anything).Return("ok", nil)
		mockRepo.On("CreateUser", mock.Anything, mock.Anything).Return(11, nil)

		s := &handler.Server{
			Repository: mockRepo,
			Pwd:        mockPwd,
		}
		err := s.PostRegister(c)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusNoContent, rec.Code)
		}
	})

}

func TestPostLogin(t *testing.T) {
	e := echo.New()

	t.Run("Bad Request - Failed to Validate", func(t *testing.T) {
		mockRepo := new(mocks.RepositoryInterface)

		reqBody := map[string]interface{}{
			"PhoneNumber": "2323244",
		}
		reqBodyBytes, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(reqBodyBytes))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		s := &handler.Server{Repository: mockRepo}

		err := s.PostLogin(c)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusBadRequest, err.(*echo.HTTPError).Code)
	})

	t.Run("User Not Found", func(t *testing.T) {
		mockRepo := new(mocks.RepositoryInterface)

		reqBody := map[string]interface{}{"PhoneNumber": "+628222667727", "fullName": "LOLTOS", "password": "@Python12345@"}
		reqBodyBytes, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(reqBodyBytes))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockRepo.On("GetUser", mock.Anything, mock.Anything).Return(nil, errors.New("user not found"))

		s := &handler.Server{Repository: mockRepo}

		err := s.PostLogin(c)
		if assert.Error(t, err) {
			httpErr, ok := err.(*echo.HTTPError)
			if assert.True(t, ok) {
				assert.Equal(t, http.StatusNotFound, httpErr.Code)
			}
		}
	})

	t.Run("Success Login", func(t *testing.T) {
		mockRepo := new(mocks.RepositoryInterface)
		mockPwd := new(pwdMocks.PasswordManagerInterface)
		mockJwt := new(authMocks.JwtInterface)

		reqBody := map[string]interface{}{"PhoneNumber": "+628222667727", "fullName": "LOLTOS", "password": "@Python12345@"}
		reqBodyBytes, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(reqBodyBytes))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockRepo.On("GetUser", mock.Anything, mock.Anything).Return(&repository.UserModel{
			ID:          111,
			PhoneNumber: "111",
			FullName:    "111",
			Password:    "11",
			SaltKey:     "111",
			CreatedAt:   time.Time{},
			UpdatedAt:   time.Time{},
		}, nil)
		mockPwd.On("VerifyPassword", mock.Anything, mock.Anything, mock.Anything).Return(true)
		mockJwt.On("CreateToken", mock.Anything, mock.Anything).Return("ok", nil)
		s := &handler.Server{Repository: mockRepo, Pwd: mockPwd, Jwt: mockJwt}

		err := s.PostLogin(c)
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})
}

func TestGetUserId(t *testing.T) {
	e := echo.New()

	t.Run("Error on ParseToken", func(t *testing.T) {
		mockJwt := &authMocks.JwtInterface{}

		req := httptest.NewRequest(http.MethodGet, "/user/1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/user/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		mockJwt.On("ParseToken", mock.Anything).Return(nil, errors.New("some error"))

		s := &handler.Server{Jwt: mockJwt}

		err := s.GetUserId(c, 1, generated.GetUserIdParams{
			Authorization: "some-token",
		})

		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}
	})

	t.Run("Forbidden", func(t *testing.T) {
		mockJwt := &authMocks.JwtInterface{}

		req := httptest.NewRequest(http.MethodGet, "/user/1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/user/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		mockJwt.On("ParseToken", mock.Anything).Return(&middleware.JwtParsedPayload{
			ID:     111,
			Expire: 111,
		}, nil)

		s := &handler.Server{Jwt: mockJwt}

		err := s.GetUserId(c, 1, generated.GetUserIdParams{
			Authorization: "some-token",
		})

		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusForbidden, rec.Code)
		}
	})

	t.Run("FetchUserById err", func(t *testing.T) {
		mockRepo := new(mocks.RepositoryInterface)
		mockJwt := &authMocks.JwtInterface{}

		req := httptest.NewRequest(http.MethodGet, "/user/1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/user/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		mockJwt.On("ParseToken", mock.Anything).Return(&middleware.JwtParsedPayload{
			ID:     1,
			Expire: 111,
		}, nil)

		mockRepo.On("GetUser", mock.Anything, mock.Anything).Return(nil, errors.New(commons.ErrorNoRow))

		s := &handler.Server{Jwt: mockJwt, Repository: mockRepo}

		err := s.GetUserId(c, 1, generated.GetUserIdParams{
			Authorization: "some-token",
		})

		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusForbidden, rec.Code)
		}
	})

	t.Run("Success GetUserId ", func(t *testing.T) {
		mockRepo := new(mocks.RepositoryInterface)
		mockJwt := &authMocks.JwtInterface{}

		req := httptest.NewRequest(http.MethodGet, "/user/1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/user/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		mockJwt.On("ParseToken", mock.Anything).Return(&middleware.JwtParsedPayload{
			ID:     1,
			Expire: 111,
		}, nil)
		mockRepo.On("GetUser", mock.Anything, mock.Anything).Return(&repository.UserModel{
			ID:          111,
			PhoneNumber: "111",
			FullName:    "111",
			Password:    "111",
			SaltKey:     "111",
			CreatedAt:   time.Time{},
			UpdatedAt:   time.Time{},
		}, nil)

		s := &handler.Server{Jwt: mockJwt, Repository: mockRepo}

		err := s.GetUserId(c, 1, generated.GetUserIdParams{
			Authorization: "some-token",
		})

		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	})

}

func TestPatchUserIdEdit(t *testing.T) {

	e := echo.New()

	t.Run("Error on ParseToken", func(t *testing.T) {
		mockJwt := new(authMocks.JwtInterface)
		mockRepo := new(mocks.RepositoryInterface)
		req := httptest.NewRequest(echo.PUT, "/", nil)
		req.Header.Set("Authorization", "someValidToken")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		mockJwt.On("ParseToken", mock.Anything).Return(&middleware.JwtParsedPayload{
			ID:     1,
			Expire: 111,
		}, errors.New("simulate err"))
		mockRepo.On("GetUser", mock.Anything, 1, mock.Anything).Return(nil)

		s := &handler.Server{Jwt: mockJwt, Repository: mockRepo}
		err := s.PatchUserIdEdit(c, 1, generated.PatchUserIdEditParams{
			Authorization: "some-token",
		})

		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}
	})

	t.Run("Status Forbidden", func(t *testing.T) {
		mockJwt := new(authMocks.JwtInterface)
		mockRepo := new(mocks.RepositoryInterface)
		req := httptest.NewRequest(echo.PUT, "/", nil)
		req.Header.Set("Authorization", "someValidToken")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		mockJwt.On("ParseToken", mock.Anything).Return(&middleware.JwtParsedPayload{
			ID:     111,
			Expire: 111,
		}, nil)
		mockRepo.On("GetUser", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		s := &handler.Server{Jwt: mockJwt, Repository: mockRepo}
		err := s.PatchUserIdEdit(c, 1, generated.PatchUserIdEditParams{
			Authorization: "some-token",
		})

		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusForbidden, rec.Code)
		}
	})

	t.Run("Err Body request", func(t *testing.T) {
		mockJwt := new(authMocks.JwtInterface)
		mockRepo := new(mocks.RepositoryInterface)

		req := httptest.NewRequest(echo.PATCH, "/", nil)
		req.Header.Set("Authorization", "someValidToken")

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		mockJwt.On("ParseToken", mock.Anything).Return(&middleware.JwtParsedPayload{
			ID:     1,
			Expire: 111,
		}, nil)
		mockRepo.On("GetUser", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		s := &handler.Server{Jwt: mockJwt, Repository: mockRepo}
		err := s.PatchUserIdEdit(c, 1, generated.PatchUserIdEditParams{
			Authorization: "some-token",
		})

		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("Get user phone number duplicate", func(t *testing.T) {
		mockJwt := new(authMocks.JwtInterface)
		mockRepo := new(mocks.RepositoryInterface)

		reqBody := map[string]interface{}{"PhoneNumber": "+628222667727", "fullName": "LOLTOS", "password": "@Python12345@"}

		bodyBytes, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(echo.PATCH, "/", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Authorization", "someValidToken")
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		mockJwt.On("ParseToken", mock.Anything).Return(&middleware.JwtParsedPayload{
			ID:     1,
			Expire: 111,
		}, nil)

		mockRepo.On("GetUser", mock.Anything, mock.Anything).Return(&repository.UserModel{
			ID:          111,
			PhoneNumber: "111",
			FullName:    "111",
			Password:    "111",
			SaltKey:     "111",
			CreatedAt:   time.Time{},
			UpdatedAt:   time.Time{},
		}, nil)
		s := &handler.Server{Jwt: mockJwt, Repository: mockRepo}
		err := s.PatchUserIdEdit(c, 1, generated.PatchUserIdEditParams{
			Authorization: "some-token",
		})

		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}
	})

	t.Run("Success Update User", func(t *testing.T) {
		mockJwt := new(authMocks.JwtInterface)
		mockRepo := new(mocks.RepositoryInterface)

		reqBody := map[string]interface{}{"PhoneNumber": "+628222667727", "fullName": "LOLTOS", "password": "@Python12345@"}

		bodyBytes, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(echo.PATCH, "/", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Authorization", "someValidToken")
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		mockJwt.On("ParseToken", mock.Anything).Return(&middleware.JwtParsedPayload{
			ID:     1,
			Expire: 111,
		}, nil)

		mockRepo.On("UpdateUser", mock.Anything, mock.Anything).Return(nil)
		mockRepo.On("GetUser", mock.Anything, mock.Anything).Return(&repository.UserModel{
			ID:          1,
			PhoneNumber: "111",
			FullName:    "111",
			Password:    "111",
			SaltKey:     "111",
			CreatedAt:   time.Time{},
			UpdatedAt:   time.Time{},
		}, nil)
		s := &handler.Server{Jwt: mockJwt, Repository: mockRepo}
		err := s.PatchUserIdEdit(c, 1, generated.PatchUserIdEditParams{
			Authorization: "some-token",
		})

		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusNoContent, rec.Code)
		}
	})

}
