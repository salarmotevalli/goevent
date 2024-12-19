package httpserver

import (
	"event-manager/config"
	"event-manager/param/eventparam"
	"event-manager/service/authservice"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (s Server) CreateEvent(c echo.Context) error {
	claim := c.Get(config.AuthMiddlewareContextKey).(*authservice.Claims)

	var req eventparam.CreateEventRequest
	err := c.Bind(&req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	req.OwnerID = claim.UserID

	// Validation

	res, err := s.eventSvc.CreateNewEvent(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusCreated, res)
}

func (s Server) IndexEvent(c echo.Context) error {
	claim := c.Get(config.AuthMiddlewareContextKey).(*authservice.Claims)

	events, err := s.eventSvc.GetAllEvents(claim.UserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	log.Println(events)

	return c.JSON(http.StatusOK, events)
}

func (s Server) ShowEvent(c echo.Context) error {
	// TODO: check user is owner of event
	// claim := c.Get(config.AuthMiddlewareContextKey).(*authservice.Claims)

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	// TODO: retrive user id from jwt token
	event, err := s.eventSvc.GetEvent(uint(id))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, event)
}

func (s Server) UpdateEvent(c echo.Context) error {
	// TODO: check user is owner of event
	// claim := c.Get(config.AuthMiddlewareContextKey).(*authservice.Claims)

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	var req eventparam.UpdateEventRequest
	c.Bind(&req)
	req.ID = uint(id)

	// TODO: validation

	res, sErr := s.eventSvc.UpdateEvent(req)
	if sErr != nil {
		return echo.NewHTTPError(http.StatusBadRequest, sErr)
	}

	return c.JSON(http.StatusAccepted, res)
}

func (s Server) DeleteEvent(c echo.Context) error {
	// TODO: check user is owner of event
	// claim := c.Get(config.AuthMiddlewareContextKey).(*authservice.Claims)

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	_, sErr := s.eventSvc.DeleteEvent(
		eventparam.DeleteEventRequest{EventID: uint(id)})

	if sErr != nil {
		return echo.NewHTTPError(http.StatusBadRequest, sErr)
	}

	return c.JSON(http.StatusNoContent, struct{}{})
}
