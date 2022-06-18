package utils

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// GenerateJWT ...
func GenerateJWT(uid uuid.UUID) string {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = uid.String()
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	t, _ := token.SignedString(JWTSecret)
	return t
}

// GetCurrentUserIdFromToken ...
func GetCurrentUserIdFromToken(c echo.Context) (uuid.UUID, error) {
	bearer := c.Request().Header.Get("Authorization")
	if bearer == "" {
		return uuid.UUID{}, errors.New("access token is required")
	}

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(bearer[7:], claims, func(token *jwt.Token) (interface{}, error) {
		return JWTSecret, nil
	})
	if err != nil {
		return uuid.UUID{}, errors.WithStack(err)
	}

	if token == nil || !token.Valid {
		return uuid.UUID{}, errors.New("access token is not validated")
	}

	uid := claims["uid"].(string)
	userId, err := uuid.Parse(uid)
	if err != nil {
		return uuid.UUID{}, errors.WithStack(err)
	}

	return userId, nil
}

func GetCurrentUserIdOrNilFromToken(c echo.Context) (*uuid.UUID, error) {
	bearer := c.Request().Header.Get("Authorization")
	if bearer == "" {
		return nil, nil
	}
	// fmt.Println(bearer)
	// fmt.Println(bearer[7:])

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(bearer[7:], claims, func(token *jwt.Token) (interface{}, error) {
		return JWTSecret, nil
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if token == nil || !token.Valid {
		return nil, errors.New("access token is not validated")
	}

	uid := claims["uid"].(string)
	userId, err := uuid.Parse(uid)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &userId, nil
}
