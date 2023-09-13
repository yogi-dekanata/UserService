package handler

import (
	"context"
	"errors"
	"github.com/SawitProRecruitment/UserService/commons"
	pwd "github.com/SawitProRecruitment/UserService/commons/mocks"
	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/SawitProRecruitment/UserService/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestUserService_RegisterNewUser(t *testing.T) {
	type args struct {
		ctx context.Context
		req *generated.UserRegisterRequest
	}
	tests := []struct {
		name               string
		args               args
		wantErr            assert.ErrorAssertionFunc
		repCreateUserErr   error
		repGenerateHashErr error
	}{
		{
			name: "Success RegisterNewUser",

			args: args{
				ctx: context.Background(),
				req: &generated.UserRegisterRequest{
					PhoneNumber: "12345",
					Password:    "password",
					FullName:    "John Doe",
				},
			},
			wantErr:          assert.NoError,
			repCreateUserErr: nil,
		},
		{
			name: "RegisterNewUser Error",

			args: args{
				ctx: context.Background(),
				req: &generated.UserRegisterRequest{
					PhoneNumber: "12345",
					Password:    "password",
					FullName:    "John Doe",
				},
			},
			wantErr:          assert.Error,
			repCreateUserErr: errors.New("simulate err"),
		},
		{
			name: "GenerateHash Error",

			args: args{
				ctx: context.Background(),
				req: &generated.UserRegisterRequest{
					PhoneNumber: "12345",
					Password:    "password",
					FullName:    "John Doe",
				},
			},
			wantErr:            assert.Error,
			repGenerateHashErr: errors.New("simulate err"),
			repCreateUserErr:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.RepositoryInterface)
			mockPwd := new(pwd.PasswordManagerInterface)

			mockRepo.On("CreateUser", tt.args.ctx, mock.Anything).Return(111, tt.repCreateUserErr)
			mockPwd.On("CreateSalt").Return("randomsalt")
			mockPwd.On("GenerateHash", "password", "randomsalt").Return("hashedpassword", tt.repGenerateHashErr)

			us := &UserService{
				Repo: mockRepo,
				Pwd:  mockPwd,
			}
			tt.wantErr(t, us.RegisterNewUser(tt.args.ctx, tt.args.req))
		})
	}
}

func TestUserService_FetchUserByPhoneNumber(t *testing.T) {
	type args struct {
		ctx         context.Context
		phoneNumber string
	}
	tests := []struct {
		name          string
		args          args
		wantErr       assert.ErrorAssertionFunc
		repGetUserErr error
	}{
		{
			name: "Success FetchUserByPhoneNumber",

			args: args{
				ctx:         context.Background(),
				phoneNumber: "12345",
			},
			wantErr:       assert.NoError,
			repGetUserErr: nil,
		},
		{
			name: "FetchUserByPhoneNumber Error",

			args: args{
				ctx:         context.Background(),
				phoneNumber: "12345",
			},
			wantErr:       assert.Error,
			repGetUserErr: errors.New("simulate err"),
		},
		{
			name: "No User Found",

			args: args{
				ctx:         context.Background(),
				phoneNumber: "12345",
			},
			wantErr:       assert.NoError,
			repGetUserErr: errors.New(commons.ErrorNoData),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.RepositoryInterface)

			mockRepo.On("GetUser", tt.args.ctx, mock.Anything).Return(&repository.UserModel{}, tt.repGetUserErr)

			us := &UserService{
				Repo: mockRepo,
			}
			_, err := us.FetchUserByPhoneNumber(tt.args.ctx, tt.args.phoneNumber)
			tt.wantErr(t, err)
		})
	}
}
