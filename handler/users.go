package handler

import (
	"context"
	"github.com/SawitProRecruitment/UserService/commons"
	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/labstack/gommon/log"
)

type UserServiceInterface interface {
	FetchUserByPhoneNumber(ctx context.Context, phoneNumber string) (*repository.UserModel, error)
	RegisterNewUser(ctx context.Context, req *generated.UserRegisterRequest) error
}

type UserService struct {
	UserService UserServiceInterface
	Repo        repository.RepositoryInterface
	Pwd         commons.PasswordManagerInterface
}

func NewUserService(repo repository.RepositoryInterface, pwd commons.PasswordManagerInterface) *UserService {
	return &UserService{
		Repo: repo,
		Pwd:  pwd,
	}
}

func (us *UserService) RegisterNewUser(ctx context.Context, req *generated.UserRegisterRequest) error {
	saltKey := us.Pwd.CreateSalt()
	hashedPass, err := us.Pwd.GenerateHash(req.Password, saltKey)
	if err != nil {
		log.Errorf("error hashing password: %v", err)
		return err
	}

	_, err = us.Repo.CreateUser(ctx, repository.UserInput{
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

func (us *UserService) FetchUserByPhoneNumber(ctx context.Context, phoneNumber string) (*repository.UserModel, error) {
	user, err := us.Repo.GetUser(ctx, repository.GetUserInput{
		PhoneNumber: &phoneNumber,
	})

	if err != nil && err.Error() != commons.ErrorNoData {
		log.Errorf("error fetching user: %v", err)
		return nil, err
	}
	return user, nil
}
