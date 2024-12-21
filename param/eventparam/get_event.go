package eventparam

type GetEventRequest struct {
	EventID uint `json:"event_id"`
}

type GetEventResponse struct {
	Event EventInfo `json:"event"`
}
