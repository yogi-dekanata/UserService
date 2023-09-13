package handler

import (
	"github.com/SawitProRecruitment/UserService/middleware"
)

// Server ...
type Server struct {
	//Jwt        commons.JwtInterface
	Middleware  middleware.IMiddleware
	UserService UserServiceInterface // Use the interface here
}

// NewServerOptions ...
type NewServerOptions struct {
	//Jwt        commons.JwtInterface
	Middleware  middleware.IMiddleware
	UserService UserServiceInterface // Use the interface here
}

// NewServer ...
func NewServer(opts NewServerOptions) *Server {
	return &Server{
		//Jwt:        opts.Jwt,
		Middleware:  opts.Middleware,
		UserService: opts.UserService, // No need to change this
	}
}
