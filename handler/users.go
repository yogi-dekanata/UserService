package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/SawitProRecruitment/UserService/commons"
	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/middleware"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/labstack/gommon/log"
)

func (s *Server) RegisterNewUser(ctx context.Context, req *generated.UserRegisterRequest) error {
	saltKey := s.Pwd.CreateSalt()
	hashedPass, err := s.Pwd.GenerateHash(req.Password, saltKey)
	if err != nil {
		log.Errorf("error hashing password: %v", err)
		return err
	}

	_, err = s.Repository.CreateUser(ctx, repository.UserInput{
		PhoneNumber: req.PhoneNumber,
		Password:    hashedPass,
		FullName:    req.FullName,
		SaltKey:     saltKey,
	})
	if err != nil {
		log.Errorf("error creating user: %v", err)
		return err
	}
	return nil
}

func (s *Server) FetchUserByPhoneNumber(ctx context.Context, phoneNumber string) (*repository.UserModel, error) {
	user, err := s.Repository.GetUser(ctx, repository.GetUserInput{
		PhoneNumber: &phoneNumber,
	})

	if err != nil && err.Error() != commons.ErrorNoData {
		log.Errorf("error fetching user: %v", err)
		return nil, err
	}
	return user, nil
}

func (s *Server) PerformLogin(ctx context.Context, req *generated.LoginRequest) (*repository.UserModel, string, error) {
	user, err := s.Repository.GetUser(ctx, repository.GetUserInput{
		PhoneNumber: &req.PhoneNumber,
	})

	if err != nil {
		if err.Error() == commons.ErrorNoRow { // assuming commons.ErrorNoRow is a constant for "user not found"
			return nil, "", errors.New("user not found")
		}
		log.Errorf("Error when checking phone number from DB: %s", err.Error())
		return nil, "", err
	}

	// Validate the password
	ok := s.Pwd.VerifyPassword(req.Password, user.Password, user.SaltKey)
	if !ok {
		return nil, "", errors.New("invalid password")
	}

	// Create JWT Token
	token, err := s.Jwt.CreateToken(middleware.UserJwtPayload{
		ID: user.ID,
	}, 9)
	if err != nil {
		log.Errorf("CreateToken, error when creating token err:%s", err.Error())
		return nil, "", err
	}

	return user, token, nil
}

func (s *Server) FetchUserById(ctx context.Context, userId int) (*repository.UserModel, error) {
	user, err := s.Repository.GetUser(ctx, repository.GetUserInput{
		ID: &userId,
	})
	if err != nil {
		log.Errorf("FetchUserById, found error fetching user by id: %v", err)
		return nil, err
	}
	return user, nil
}

func (s *Server) EditUser(ctx context.Context, userId int, req *generated.UserEditRequest) error {
	// Check if the phone number is already taken
	user, err := s.Repository.GetUser(ctx, repository.GetUserInput{PhoneNumber: req.PhoneNumber})
	if err != nil && err.Error() != commons.ErrorNoRow {
		return err
	}
	if user != nil && user.ID != userId {
		return fmt.Errorf("phone number already exist")
	}

	// Update user
	err = s.Repository.UpdateUser(ctx, repository.UserInput{
		ID:          userId,
		PhoneNumber: *req.PhoneNumber,
		FullName:    *req.FullName,
	})
	if err != nil && err.Error() != commons.ErrorNoRow {
		return err
	}

	return nil
}
