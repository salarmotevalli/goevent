package eventparam

import (
	"event-manager/entity"
	"time"
)

type EventInfo struct {
	ID       uint      `json:"id"`
	Title    string    `json:"title,omitempty"`
	Location string    `json:"location,omitempty"`
	StartAt  time.Time `json:"start_at,omitempty"`
}

func (e *EventInfo) FillFromEventEntity(event entity.Event) {
	e.ID = event.ID
	e.Title = event.Title
	e.Location = event.Location
	e.StartAt = event.StartAt
}
