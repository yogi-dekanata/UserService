package main

import (
	"github.com/SawitProRecruitment/UserService/commons"
	"os"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/handler"
	"github.com/SawitProRecruitment/UserService/middleware"
	"github.com/SawitProRecruitment/UserService/repository"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func main() {
	e := echo.New()
	srv, err := newServer()
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}
	var server generated.ServerInterface = srv

	generated.RegisterHandlers(e, server)

	customGroup := e.Group("")
	customGroup.Use(srv.Middleware.Auth)

	e.Logger.Fatal(e.Start(":8080"))
}

func newServer() (*handler.Server, error) {
	dbDsn := os.Getenv("DATABASE_URL")

	var repo repository.RepositoryInterface = repository.NewRepository(repository.NewRepositoryOptions{
		Dsn: dbDsn,
	})

	jwt := &middleware.Jwt{
		PrivateKey: nil,
		PublicKey:  nil,
	}
	m := middleware.NewMiddleware(jwt, repo)

	// Initialize UserService here
	userService := handler.NewUserService(repo, &commons.PasswordManager{})

	opts := handler.NewServerOptions{
		UserService: userService,
		Middleware:  m,
	}

	return handler.NewServer(opts), nil
}
