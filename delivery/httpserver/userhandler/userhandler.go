package userhandler

import (
	"event-manager/param/userparam"
	"event-manager/service/authservice"
	"event-manager/service/userservice"
	"github.com/labstack/echo/v4"
	"net/http"
)

type UserHandler struct {
	UserSvc userservice.UserService
	AuthSvc authservice.AuthService
}

func New(us userservice.UserService, as authservice.AuthService) UserHandler {
	return UserHandler{
		UserSvc: us,
		AuthSvc: as,
	}
}

func (h UserHandler) Login(c echo.Context) error {
	var request userparam.LoginRequest
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if len(request.Password) < 8 || len(request.Username) < 3 {
		return echo.NewHTTPError(http.StatusUnprocessableEntity)
	}

	result, err := h.UserSvc.Login(request)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusAccepted, result)
}

func (h UserHandler) Register(c echo.Context) error {
	var request userparam.RegisterUserRequest
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if len(request.Password) < 8 || len(request.Username) < 3 {
		return echo.NewHTTPError(http.StatusUnprocessableEntity)
	}

	result, err := h.UserSvc.Register(request)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusCreated, result)
}
