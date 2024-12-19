package eventhandler

import (
	"event-manager/config"
	"event-manager/param/eventparam"
	"event-manager/service/authservice"
	"event-manager/service/eventservice"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type EventHandler struct {
	eventSvc eventservice.EventService
}

func New(es eventservice.EventService) EventHandler {
	return EventHandler {
		eventSvc: es,
	}
}

func (h EventHandler) CreateEvent(c echo.Context) error {
	claim := c.Get(config.AuthMiddlewareContextKey).(*authservice.Claims)

	var req eventparam.CreateEventRequest
	err := c.Bind(&req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	req.OwnerID = claim.UserID

	// Validation

	res, err := h.eventSvc.CreateNewEvent(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusCreated, res)
}

func (h EventHandler) IndexEvent(c echo.Context) error {
	claim := c.Get(config.AuthMiddlewareContextKey).(*authservice.Claims)

	events, err := h.eventSvc.GetAllEvents(claim.UserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, events)
}

func (h EventHandler) ShowEvent(c echo.Context) error {
	// TODO: check user is owner of event
	// claim := c.Get(config.AuthMiddlewareContextKey).(*authservice.Claims)

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	// TODO: retrive user id from jwt token
	event, err := h.eventSvc.GetEvent(uint(id))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, event)
}

func (h EventHandler) UpdateEvent(c echo.Context) error {
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

	res, sErr := h.eventSvc.UpdateEvent(req)
	if sErr != nil {
		return echo.NewHTTPError(http.StatusBadRequest, sErr)
	}

	return c.JSON(http.StatusAccepted, res)
}

func (h EventHandler) DeleteEvent(c echo.Context) error {
	// TODO: check user is owner of event
	// claim := c.Get(config.AuthMiddlewareContextKey).(*authservice.Claims)

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	_, sErr := h.eventSvc.DeleteEvent(
		eventparam.DeleteEventRequest{EventID: uint(id)})

	if sErr != nil {
		return echo.NewHTTPError(http.StatusBadRequest, sErr)
	}

	return c.JSON(http.StatusNoContent, struct{}{})
}
