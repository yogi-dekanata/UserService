package handler

import (
	"fmt"
	"github.com/SawitProRecruitment/UserService/commons"
	"github.com/go-playground/validator/v10"
	"net/http"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/labstack/echo/v4"
)

func (s *Server) GetHealth(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, generated.SuccessResponse{Message: "HI"})
}

func (s *Server) PostRegister(ctx echo.Context) error {
	userRegisterRequest, err := bindAndValidateRequest(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, err := s.UserService.FetchUserByPhoneNumber(ctx.Request().Context(), userRegisterRequest.PhoneNumber)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "System error")
	}

	if user != nil {
		return echo.NewHTTPError(http.StatusConflict, "Phone number already exists")
	}

	if err := s.UserService.RegisterNewUser(ctx.Request().Context(), userRegisterRequest); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "System error")
	}

	return ctx.JSON(http.StatusNoContent, generated.SuccessResponse{Message: "User successfully created"})
}

func bindAndValidateRequest(ctx echo.Context) (*generated.UserRegisterRequest, error) {
	req := &generated.UserRegisterRequest{}
	if err := ctx.Bind(req); err != nil {
		return nil, fmt.Errorf(commons.InValidData, err)
	}

	if err := commons.Validate.Struct(req); err != nil {
		return nil, fmt.Errorf(commons.InValidData, err.(validator.ValidationErrors)[0])
	}
	return req, nil
}

func (s *Server) PostLogin(ctx echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (s *Server) GetUserId(ctx echo.Context, id int, params generated.GetUserIdParams) error {
	//TODO implement me
	panic("implement me")
}

func (s *Server) PatchUserIdEdit(ctx echo.Context, id int, params generated.PatchUserIdEditParams) error {
	//TODO implement me
	panic("implement me")
}
