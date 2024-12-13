package middleware

import (
	cfg "event-manager/config"
	"event-manager/service/authservice"
	"github.com/labstack/echo/v4"
	mw "github.com/labstack/echo-jwt/v4"
)

func Auth(service authservice.AuthService, config authservice.Config) echo.MiddlewareFunc {
	return mw.WithConfig(mw.Config{
		ContextKey: cfg.AuthMiddlewareContextKey,
		SigningKey: []byte(config.SignKey),
		// TODO - as sign method string t	o config...
		SigningMethod: "HS256",
		ParseTokenFunc: func(c echo.Context, auth string) (interface{}, error) {
			claims, err := service.VerifyToken(auth)
			if err != nil {
				return nil, err
			}

			return claims, nil
		},
	})
}