package httpserver

import (
	"event-manager/service/userservice"
	// "fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s Server) Login(c echo.Context) error {
	var request userservice.LoginRequest
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if len(request.Password) < 8 || len(request.Username) < 3 {
		return echo.NewHTTPError(http.StatusUnprocessableEntity)
	}

	result, err := s.userSvc.Login(request)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusAccepted, result)
}

func (s Server) Register(c echo.Context) error {
	var request userservice.RegisterRequest
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if len(request.Password) < 8 || len(request.Username) < 3 {
		return echo.NewHTTPError(http.StatusUnprocessableEntity)
	}

	result, err := s.userSvc.Register(request)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusCreated, result)
}
