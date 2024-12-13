package httpserver

import (
	"encoding/json"
	"errors"
	"event-manager/config"
	"event-manager/service/authservice"
	"event-manager/service/eventservice"
	"event-manager/service/userservice"
	internalMiddleware "event-manager/delivery/httpserver/middleware"
	"log"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	config   config.Config
	userSvc  userservice.UserService
	authSvc  authservice.AuthService
	eventSvc eventservice.EventService
}

func New(cnf config.Config,
	userservice userservice.UserService,
	authservice authservice.AuthService,
	eventservice eventservice.EventService) Server {
	return Server{
		config:   cnf,
		userSvc:  userservice,
		authSvc:  authservice,
		eventSvc: eventservice,
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
	userGroup.POST("/register", s.Register)
	userGroup.POST("/login", s.Login)

	// event
	eventGroup := e.Group("/events", internalMiddleware.Auth(s.authSvc, s.config.AuthConfig))
	eventGroup.GET("/", s.IndexEvent)
	eventGroup.GET("/:id", s.ShowEvent)
	eventGroup.POST("/", s.CreateEvent)
	eventGroup.PUT("/:id", s.UpdateEvent)
	eventGroup.DELETE("/:id", s.DeleteEvent)

	// logRoutes(e.Routes())

	if err := e.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to start server", "error", err)
	}

}

func logRoutes(routes []*echo.Route) {
	data, err := json.MarshalIndent(routes, "", "  ")
	if err != nil {
		panic(err)
	}

	log.Println(string(data))
}
