package eventservice

import (
	"errors"
	"event-manager/entity"
	"event-manager/param/eventparam"
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

func (s *EventService) GetAllEvents(userId uint) (eventparam.GetAllEventResponse, error) {
	events, err := s.repo.GetAllEventsFor(userId)
	if err != nil {
		return eventparam.GetAllEventResponse{}, err
	}

	eventInfos := make([]eventparam.EventInfo, 0)
	for _, e := range events {
		var ei eventparam.EventInfo
		ei.FillFromEventEntity(e)
		eventInfos = append(eventInfos, ei)
	}

	return eventparam.GetAllEventResponse{Events: eventInfos}, nil
}

func (s *EventService) GetEvent(eventId uint) (eventparam.GetEventResponse, error) {
	event, exist, err := s.repo.GetEventByID(eventId)
	if err != nil {
		return eventparam.GetEventResponse{}, err
	}

	if !exist {
		return eventparam.GetEventResponse{}, errors.New("not found")
	}

	var ei eventparam.EventInfo
	ei.FillFromEventEntity(event)

	return eventparam.GetEventResponse{
		Event: ei,
	}, nil
}

func (s *EventService) CreateNewEvent(req eventparam.CreateEventRequest) (eventparam.CreateEventResponse, error) {
	var e entity.Event

	e.SetOwner(req.OwnerID)
	e.SetTitle(req.Title)
	e.SetLocation(req.Location)
	e.SetStartAt(req.StartAt)
	e.Activate()

	event, err := s.repo.CreateEvent(e)
	var eventInfo eventparam.EventInfo
	eventInfo.FillFromEventEntity(event)

	return eventparam.CreateEventResponse{
		Event: eventInfo,
	}, err
}

func (s *EventService) UpdateEvent(req eventparam.UpdateEventRequest) (eventparam.UpdateEventResponse, error) {
	// fetch event
	event, exist, err := s.repo.GetEventByID(req.ID)
	if err != nil {
		return eventparam.UpdateEventResponse{}, err
	}

	if !exist {
		return eventparam.UpdateEventResponse{}, errors.New("not found")
	}

	// update entity
	event.SetTitle(req.Title)
	event.SetLocation(req.Location)
	event.SetStartAt(req.StartAt)

	rErr := s.repo.UpdateEvent(event)
	var eventInfo eventparam.EventInfo
	eventInfo.FillFromEventEntity(event)

	return eventparam.UpdateEventResponse{
		Event: eventInfo,
	}, rErr
}

func (s *EventService) DeleteEvent(req eventparam.DeleteEventRequest) (eventparam.DeleteEventResponse, error) {
	_, exist, err := s.repo.GetEventByID(req.EventID)
	if err != nil {
		return eventparam.DeleteEventResponse{}, err
	}

	// check is there or notfound
	if !exist {
		return eventparam.DeleteEventResponse{}, errors.New("not found")
	}

	rErr := s.repo.DeleteEvent(req.EventID)

	return eventparam.DeleteEventResponse{}, rErr
}
