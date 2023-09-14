package middleware

import (
	"fmt"
	"github.com/SawitProRecruitment/UserService/commons"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
	"time"
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

// IMiddlewareInterface ...
type IMiddlewareInterface interface {
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

func (j *Jwt) CreateToken(jwtData UserJwtPayload, expireInHour int) (string, error) {
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(j.PrivateKey)
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %w", err)
	}

	token := jwt.New(jwt.SigningMethodRS256)
	claims := token.Claims.(jwt.MapClaims)
	claims[commons.IDClaimKey] = fmt.Sprint(jwtData)
	claims[commons.ExpClaimKey] = time.Now().Add(time.Hour * time.Duration(expireInHour)).Unix()

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return tokenString, nil
}

func (j *Jwt) IsValid(tokenString string) (bool, error) {
	jwtData, err := j.ParseToken(tokenString)
	if err != nil {
		return false, fmt.Errorf("failed to validate token: %w", err)
	}
	return time.Unix(jwtData.Expire, 0).After(time.Now()), nil
}

func (j *Jwt) ParseToken(tokenString string) (*JwtParsedPayload, error) {
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(j.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	id, err := commons.ConvertInterfaceToInt(claims[commons.IDClaimKey])
	if err != nil {
		return nil, fmt.Errorf("failed to convert ID: %w", err)
	}

	exp, err := commons.ConvertInterfaceToInt64(claims[commons.ExpClaimKey])
	if err != nil {
		return nil, fmt.Errorf("failed to convert Expire time: %w", err)
	}

	return &JwtParsedPayload{ID: id, Expire: exp}, nil
}
