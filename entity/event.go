package entity

import "time"

type Event struct {
	ID       uint
	OwnerID  uint
	Title    string
	Location string
	StartAt  time.Time
	Status   EventStatus
}

func (e *Event) BelongsTo(u uint) bool { return e.OwnerID == u }
func (e *Event) Activate()             { e.Status = EvenetActiveStatus }
func (e *Event) Deactivate()           { e.Status = EvenetDeactiveStatus }

type EventStatus uint8

const (
	EvenetActiveStatus EventStatus = iota + 1
	EvenetDeactiveStatus
)

func (s EventStatus) IsValid() bool {
	switch s {
	case EvenetActiveStatus, EvenetDeactiveStatus:
		return true
	default:
		return false
	}
}
