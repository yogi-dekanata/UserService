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
	Middleware middleware.IMiddlewareInterface
}

type NewServerOptions struct {
	Repository repository.RepositoryInterface
	Jwt        middleware.JwtInterface
	Pwd        commons.PasswordManagerInterface
	Middleware middleware.IMiddlewareInterface
}

func NewServer(opts NewServerOptions) *Server {
	return &Server{
		Repository: opts.Repository,
		Jwt:        opts.Jwt,
		Pwd:        opts.Pwd,
		Middleware: opts.Middleware,
	}
}
