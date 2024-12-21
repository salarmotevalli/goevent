package eventservice

import (
	"event-manager/entity"
	"event-manager/param/eventparam"
	"event-manager/pkg/richerror"
)

type EventRepo interface {
	GetAllEventsFor(uint) ([]entity.Event, error)
	GetEventByID(uint) (entity.Event, bool, error)
	CreateEvent(entity.Event) (entity.Event, error)
	UpdateEvent(entity.Event) error
	DeleteEvent(uint) error
}

type EventService struct {
	repo EventRepo
}

func New(er EventRepo) EventService {
	return EventService{
		repo: er,
	}
}

func (s *EventService) GetAllEvents(req eventparam.GetAllEventRequest) (eventparam.GetAllEventResponse, error) {
	const op = "eventservice.GetAllEvents"

	events, err := s.repo.GetAllEventsFor(req.UserID)
	if err != nil {
		return eventparam.GetAllEventResponse{}, richerror.New(op).WithErr(err).
			WithMeta(map[string]interface{}{"req": req})
	}

	eventInfos := make([]eventparam.EventInfo, 0)
	for _, e := range events {
		var ei eventparam.EventInfo
		ei.FillFromEventEntity(e)
		eventInfos = append(eventInfos, ei)
	}

	return eventparam.GetAllEventResponse{Events: eventInfos}, nil
}

func (s *EventService) GetEvent(req eventparam.GetEventRequest) (eventparam.GetEventResponse, error) {
	const op = "eventservice.GetEvent"

	event, exist, err := s.repo.GetEventByID(req.EventID)
	if err != nil {
		return eventparam.GetEventResponse{}, richerror.New(op).WithErr(err).
			WithMeta(map[string]interface{}{"req": req})
	}

	if !exist {
		return eventparam.GetEventResponse{}, richerror.New(op).WithKind(richerror.KindNotFound).
			WithMeta(map[string]interface{}{"req": req})
	}

	var ei eventparam.EventInfo
	ei.FillFromEventEntity(event)

	return eventparam.GetEventResponse{
		Event: ei,
	}, nil
}

func (s *EventService) CreateNewEvent(req eventparam.CreateEventRequest) (eventparam.CreateEventResponse, error) {
	const op = "eventservice.CreateNewEvent"
	
	e := entity.Event{
		OwnerID:  req.OwnerID,
		Title:    req.Title,
		Location: req.Location,
		StartAt:  req.StartAt,
		Status:   entity.EvenetActiveStatus,
	}

	event, err := s.repo.CreateEvent(e)
	if err != nil {
		return eventparam.CreateEventResponse{}, richerror.New(op).WithErr(err).
			WithMeta(map[string]interface{}{"req": req})
	}

	var eventInfo eventparam.EventInfo
	eventInfo.FillFromEventEntity(event)

	return eventparam.CreateEventResponse{
		Event: eventInfo,
	}, nil
}

func (s *EventService) UpdateEvent(req eventparam.UpdateEventRequest) (eventparam.UpdateEventResponse, error) {
	const op = "eventservice.UpdateEvent"

	event, exist, err := s.repo.GetEventByID(req.ID)
	if err != nil {
		return eventparam.UpdateEventResponse{}, richerror.New(op).WithErr(err).
			WithMeta(map[string]interface{}{"req": req})
	}

	if !exist {
		return eventparam.UpdateEventResponse{}, richerror.New(op).WithKind(richerror.KindNotFound).
			WithMeta(map[string]interface{}{"req": req})
	}

	// update entity
	event.Title = req.Title
	event.Location = req.Location
	event.StartAt = req.StartAt

	rErr := s.repo.UpdateEvent(event)
	if rErr != nil {
		return eventparam.UpdateEventResponse{}, richerror.New(op).WithErr(rErr).
			WithMeta(map[string]interface{}{"req": req})
	}

	var eventInfo eventparam.EventInfo
	eventInfo.FillFromEventEntity(event)

	return eventparam.UpdateEventResponse{
		Event: eventInfo,
	}, nil
}

func (s *EventService) DeleteEvent(req eventparam.DeleteEventRequest) (eventparam.DeleteEventResponse, error) {
	const op = "eventservice.DeleteEvent"

	_, exist, err := s.repo.GetEventByID(req.EventID)
	if err != nil {
		return eventparam.DeleteEventResponse{}, richerror.New(op).WithErr(err).
			WithMeta(map[string]interface{}{"req": req})
	}

	// check is there or notfound
	if !exist {
		return eventparam.DeleteEventResponse{}, richerror.New(op).WithKind(richerror.KindNotFound).
			WithMeta(map[string]interface{}{"req": req})
	}

	rErr := s.repo.DeleteEvent(req.EventID)
	if rErr != nil {
		return eventparam.DeleteEventResponse{}, richerror.New(op).WithErr(rErr).
			WithMeta(map[string]interface{}{"req": req})
	}

	return eventparam.DeleteEventResponse{}, nil
}
