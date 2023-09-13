package middleware

import (
	"net/http"

	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

// UserJwtPayload ...
type UserJwtPayload struct {
	ID int
}

// JwtParsedPayload ...
type JwtParsedPayload struct {
	ID     int
	Expire int64
}

// Middleware ...
type Middleware struct {
	Jwt        JwtInterface
	Repository repository.RepositoryInterface
}

type Jwt struct {
	PrivateKey []byte `json:"private_key"`
	PublicKey  []byte `json:"public_key"`
}

// Payload ...
type Payload struct {
	UserID string
}

// JwtInterface ...
type JwtInterface interface {
	CreateToken(jwtData UserJwtPayload, expireInHour int) (string, error)
	ParseToken(tokenString string) (*JwtParsedPayload, error)
	IsValid(tokenString string) (bool, error)
}

// IMiddleware ...
type IMiddleware interface {
	Auth(next echo.HandlerFunc) echo.HandlerFunc
}

// NewMiddleware for creating new middleware
func NewMiddleware(jwt JwtInterface, repo repository.RepositoryInterface) *Middleware {
	return &Middleware{Jwt: jwt, Repository: repo}
}

// Auth function
func (m Middleware) Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		headers := c.Request().Header
		valueList, found := headers[http.CanonicalHeaderKey("Authorization")]
		if !found {
			return echo.NewHTTPError(http.StatusForbidden, "missing Authorization Header")
		}

		ok, err := m.Jwt.IsValid(valueList[0]) // Assuming IsValid returns an additional payload if valid
		if err != nil {
			log.Errorf("Auth Error: %s", err.Error())
			return echo.NewHTTPError(http.StatusForbidden, "invalid Authorization Token")
		}
		if !ok {
			return echo.NewHTTPError(http.StatusForbidden, "authorization Token is Not Valid")
		}

		return next(c)
	}
}

func (j Jwt) CreateToken(jwtData UserJwtPayload, expireInHour int) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (j Jwt) ParseToken(tokenString string) (*JwtParsedPayload, error) {
	//TODO implement me
	panic("implement me")
}

func (j Jwt) IsValid(tokenString string) (bool, error) {
	//TODO implement me
	panic("implement me")
}
