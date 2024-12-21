package httpserver

import (
	"errors"
	"event-manager/config"
	"event-manager/delivery/httpserver/eventhandler"
	internalMiddleware "event-manager/delivery/httpserver/middleware"
	"event-manager/delivery/httpserver/userhandler"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	config       config.Config
	eventHandler eventhandler.EventHandler
	userHandler  userhandler.UserHandler
}

func New(cnf config.Config,
	uh userhandler.UserHandler,
	eh eventhandler.EventHandler) Server {
	return Server{
		config:       cnf,
		eventHandler: eh,
		userHandler:  uh,
	}
}

func (s Server) Serve() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("health-check", s.healthCheck)

	// auth
	userGroup := e.Group("/users")
	userGroup.POST("/register", s.userHandler.Register)
	userGroup.POST("/login", s.userHandler.Login)

	// event
	eventGroup := e.Group("/events")
	eventAuth := eventGroup.Group("/", internalMiddleware.Auth(s.userHandler.AuthSvc, s.config.AuthConfig))

	eventGroup.GET("/", s.eventHandler.IndexEvent)
	eventGroup.GET("/:id", s.eventHandler.ShowEvent)
	eventAuth.POST("/", s.eventHandler.CreateEvent)
	eventAuth.PUT("/:id", s.eventHandler.UpdateEvent)
	eventAuth.DELETE("/:id", s.eventHandler.DeleteEvent)

	if err := e.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to start server", "error", err)
	}
}
