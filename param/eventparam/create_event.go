package eventparam

import (
	"time"
)

type CreateEventRequest struct {
	Title    string    `json:"title"`
	OwnerID  uint      `json:"-"`
	Location string    `json:"location"`
	StartAt  time.Time `json:"start_at"`
}

type CreateEventResponse struct {
	Event EventInfo `json:"event"`
} 