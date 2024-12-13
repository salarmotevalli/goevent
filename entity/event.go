package entity

import "time"

type Event struct {
	id       uint
	ownerID  uint
	title    string
	location string
	startAt  time.Time
	status   EventStatus
}

// setter & getter
func (e *Event) ID() uint      { return e.id }
func (e *Event) SetID(id uint) { e.id = id }

func (e *Event) SetOwner(u uint)       { e.ownerID = u }
func (e *Event) OwnerID() uint         { return e.ownerID }
func (e *Event) BelongsTo(u uint) bool { return e.ownerID == u }

func (e *Event) SetTitle(s string) { e.title = s }
func (e *Event) Title() string     { return e.title }

func (e *Event) SetLocation(l string) { e.location = l }
func (e *Event) Location() string     { return e.location }

func (e *Event) SetStartAt(t time.Time) { e.startAt = t }
func (e *Event) StartAt() time.Time     { return e.startAt }

func (e *Event) SetStatus(s EventStatus) { e.status = s }
func (e *Event) Activate()               { e.SetStatus(EvenetActiveStatus) }
func (e *Event) Deactivate()             { e.SetStatus(EvenetDeactiveStatus) }
func (e *Event) Status() EventStatus     { return e.status }

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
