package handler

import (
	"github.com/SawitProRecruitment/UserService/commons"
	"github.com/SawitProRecruitment/UserService/middleware"
	"github.com/SawitProRecruitment/UserService/repository"
)

type Server struct {
	Repository repository.RepositoryInterface
	Jwt        middleware.JwtInterface
	Pwd        commons.PasswordManagerInterface
	Middleware middleware.IMiddleware
}

type NewServerOptions struct {
	Repository repository.RepositoryInterface
	Jwt        middleware.JwtInterface
	Pwd        commons.PasswordManagerInterface
	Middleware middleware.IMiddleware
}

func NewServer(opts NewServerOptions) *Server {
	return &Server{
		Repository: opts.Repository,
		Jwt:        opts.Jwt,
		Pwd:        opts.Pwd,
		Middleware: opts.Middleware,
	}
}
