package main

import (
	"github.com/SawitProRecruitment/UserService/commons"
	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/handler"
	"github.com/SawitProRecruitment/UserService/middleware"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/labstack/echo/v4"
	"log"
	"os"
)

func main() {
	e := echo.New()

	server, err := initializeServer()
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}

	registerRoutes(e, server)

	log.Fatal(e.Start(":8080"))
}

func initializeServer() (*handler.Server, error) {
	dbDsn := getEnv("DATABASE_URL", "")

	//dbDsn := "host=localhost port=5432 user=yogidekanata dbname=users password=newpassword sslmode=disable"

	repo := repository.NewRepository(repository.NewRepositoryOptions{Dsn: dbDsn})

	privateKey, publicKey, err := readKeys("private_key.pem", "public_key.pem")
	if err != nil {
		return nil, err
	}

	jwtMiddleware := &middleware.Jwt{PrivateKey: privateKey, PublicKey: publicKey}
	middlewareInstance := middleware.NewMiddleware(jwtMiddleware, repo)

	return handler.NewServer(handler.NewServerOptions{
		Middleware: middlewareInstance,
		Repository: repo,
		Pwd:        &commons.PasswordManager{},
		Jwt:        jwtMiddleware,
	}), nil
}

func registerRoutes(e *echo.Echo, server *handler.Server) {
	e.GET("/health", server.GetHealth) // Assume HealthCheckHandler is the method you use to handle health checks
	authGroup := e.Group("/auth")
	authGroup.Use(server.Middleware.Auth)
	generated.RegisterHandlers(e, server)

}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

func readKeys(privateKeyPath, publicKeyPath string) ([]byte, []byte, error) {
	privateKey, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, nil, err
	}

	publicKey, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, nil, err
	}

	return privateKey, publicKey, nil
}
