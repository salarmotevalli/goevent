package eventhandler

import (
	"event-manager/config"
	"event-manager/param/eventparam"
	"event-manager/pkg/httpmsg"
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
	return EventHandler{
		eventSvc: es,
	}
}

func (h EventHandler) CreateEvent(c echo.Context) error {
	claim := c.Get(config.AuthMiddlewareContextKey).(*authservice.Claims)

	var req eventparam.CreateEventRequest
	err := c.Bind(&req)
	if err != nil {
		msg, code := httpmsg.Error(err)
		return c.JSON(code, echo.Map{
			"message": msg})
	}

	req.OwnerID = claim.UserID

	res, err := h.eventSvc.CreateNewEvent(req)
	if err != nil {
		msg, code := httpmsg.Error(err)
		return c.JSON(code, echo.Map{
			"message": msg})
	}

	return c.JSON(http.StatusCreated, res)
}

func (h EventHandler) IndexEvent(c echo.Context) error {
	claim := c.Get(config.AuthMiddlewareContextKey).(*authservice.Claims)

	var req eventparam.GetAllEventRequest
	req.UserID = claim.UserID

	events, err := h.eventSvc.GetAllEvents(req)
	if err != nil {
		msg, code := httpmsg.Error(err)
		return c.JSON(code, echo.Map{
			"message": msg})
	}

	return c.JSON(http.StatusOK, events)
}

func (h EventHandler) ShowEvent(c echo.Context) error {
	// claim := c.Get(config.AuthMiddlewareContextKey).(*authservice.Claims)

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		msg, code := httpmsg.Error(err)
		return c.JSON(code, echo.Map{
			"message": msg})
	}

	var req eventparam.GetEventRequest
	req.EventID = uint(id)

	event, err := h.eventSvc.GetEvent(req)
	if err != nil {
		msg, code := httpmsg.Error(err)
		return c.JSON(code, echo.Map{
			"message": msg})
	}

	return c.JSON(http.StatusOK, event)
}

func (h EventHandler) UpdateEvent(c echo.Context) error {
	// TODO: check user is owner of event
	// claim := c.Get(config.AuthMiddlewareContextKey).(*authservice.Claims)

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		msg, code := httpmsg.Error(err)
		return c.JSON(code, echo.Map{
			"message": msg})
	}

	var req eventparam.UpdateEventRequest
	c.Bind(&req)
	req.ID = uint(id)

	res, sErr := h.eventSvc.UpdateEvent(req)
	if sErr != nil {
		msg, code := httpmsg.Error(err)
		return c.JSON(code, echo.Map{
			"message": msg})
	}

	return c.JSON(http.StatusAccepted, res)
}

func (h EventHandler) DeleteEvent(c echo.Context) error {
	// TODO: check user is owner of event
	// claim := c.Get(config.AuthMiddlewareContextKey).(*authservice.Claims)

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		msg, code := httpmsg.Error(err)
		return c.JSON(code, echo.Map{
			"message": msg})
	}

	_, sErr := h.eventSvc.DeleteEvent(
		eventparam.DeleteEventRequest{EventID: uint(id)})

	if sErr != nil {
		msg, code := httpmsg.Error(err)
		return c.JSON(code, echo.Map{
			"message": msg})
	}

	return c.JSON(http.StatusNoContent, struct{}{})
}
