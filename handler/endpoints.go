package handler

import (
	"fmt"
	"github.com/SawitProRecruitment/UserService/commons"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/gommon/log"
	"net/http"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/labstack/echo/v4"
)

func (s *Server) GetHealth(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, generated.SuccessResponse{Message: "HI"})
}

func (s *Server) PostRegister(ctx echo.Context) error {
	userRegisterRequest := &generated.UserRegisterRequest{}
	err := bindAndValidate(ctx, userRegisterRequest)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, err := s.FetchUserByPhoneNumber(ctx.Request().Context(), userRegisterRequest.PhoneNumber)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, commons.ErrSystemError)
	}

	if user != nil {
		return echo.NewHTTPError(http.StatusConflict, commons.ErrUserExists)
	}

	if err := s.RegisterNewUser(ctx.Request().Context(), userRegisterRequest); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, commons.ErrSystemError)
	}

	return ctx.JSON(http.StatusNoContent, generated.SuccessResponse{Message: "User successfully created"})
}

func (s *Server) PostLogin(ctx echo.Context) error {
	loginRequest := &generated.LoginRequest{}
	err := bindAndValidate(ctx, loginRequest)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, token, err := s.PerformLogin(ctx.Request().Context(), loginRequest)
	if err != nil {
		if err.Error() == commons.ErrorInvalidPassword {
			return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
		} else if err.Error() == "user not found" {
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, commons.ErrSystemError)
	}

	return ctx.JSON(http.StatusOK, generated.LoginResponse{
		Jwt:    token,
		UserId: user.ID,
	})
}

func (s *Server) GetUserId(ctx echo.Context, id int, params generated.GetUserIdParams) error {

	data, err := s.Jwt.ParseToken(params.Authorization)
	if err != nil {
		log.Errorf("ParseToken, error when creating token err:%s", err.Error())
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: commons.ErrSystemError,
		})
	}

	if data.ID != id {
		return ctx.JSON(http.StatusForbidden, generated.ErrorResponse{
			Message: "Forbidden",
		})
	}

	user, err := s.FetchUserById(ctx.Request().Context(), data.ID)
	if err != nil {
		if err.Error() == commons.ErrorNoRow {
			log.Warn(fmt.Sprintf("GetUser, not found userId:%d err:%s", data.ID, err.Error()))
			return ctx.JSON(http.StatusForbidden, generated.ErrorResponse{
				Message: "Forbidden",
			})
		}
		log.Errorf("GetUser, found error get user from db err:%s", err.Error())
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: commons.ErrSystemError,
		})
	}

	return ctx.JSON(http.StatusOK, generated.UserResponse{
		UserId:      &user.ID,
		FullName:    &user.FullName,
		PhoneNumber: &user.PhoneNumber,
	})
}

func (s *Server) PatchUserIdEdit(ctx echo.Context, id int, params generated.PatchUserIdEditParams) error {
	data, err := s.Jwt.ParseToken(params.Authorization)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: commons.ErrSystemError})
	}

	if data.ID != id {
		return ctx.JSON(http.StatusForbidden, generated.ErrorResponse{Message: "Forbidden"})
	}

	userEditRequest := &generated.UserEditRequest{}
	err = bindAndValidate(ctx, userEditRequest)

	if err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: err.Error()})
	}

	err = s.EditUser(ctx.Request().Context(), data.ID, userEditRequest)
	if err != nil {
		if err.Error() == commons.ErrUserExists {
			return ctx.JSON(http.StatusConflict, generated.ErrorResponse{Message: err.Error()})
		}
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: commons.ErrSystemError})
	}

	return ctx.JSON(http.StatusNoContent, generated.SuccessResponse{Message: "success update user"})
}

func bindAndValidate(ctx echo.Context, req interface{}) error {
	// Bind the request
	if err := ctx.Bind(req); err != nil {
		return fmt.Errorf("invalid data: %v", err)
	}

	// Validate the request
	if err := commons.Validate.Struct(req); err != nil {
		return fmt.Errorf("invalid data %v", err.(validator.ValidationErrors))
	}

	return nil
}
