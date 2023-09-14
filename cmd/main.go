package main

import (
	"log"
	"os"

	"github.com/SawitProRecruitment/UserService/commons"
	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/handler"
	"github.com/SawitProRecruitment/UserService/middleware"
	"github.com/SawitProRecruitment/UserService/repository"

	"github.com/labstack/echo/v4"
)

func main() {
	e := initEcho()

	server, err := newServer()
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}

	generated.RegisterHandlers(e, server)

	customGroup := e.Group("")
	customGroup.Use(server.Middleware.Auth)

	log.Fatal(e.Start(":8080"))
}

func initEcho() *echo.Echo {
	e := echo.New()
	return e
}

func newServer() (*handler.Server, error) {
	dbDsn := os.Getenv("DATABASE_URL")
	repo := repository.NewRepository(repository.NewRepositoryOptions{Dsn: dbDsn})

	jwtMiddleware := &middleware.Jwt{}
	middle := middleware.NewMiddleware(jwtMiddleware, repo)

	passwordManager := &commons.PasswordManager{}
	serverOptions := handler.NewServerOptions{
		Middleware: middle,
		Repository: repo,
		Pwd:        passwordManager,
	}

	return handler.NewServer(serverOptions), nil
}
