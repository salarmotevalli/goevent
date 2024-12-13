package eventparam

import "time"

type UpdateEventRequest struct {
	ID       uint
	Title    string
	UserID   uint
	Location string
	StartAt  time.Time
}

type UpdateEventResponse struct {
	Event EventInfo `json:"event"`
} 