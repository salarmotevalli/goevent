package eventparam

type ShhowEventRequest struct{}

type GetEventResponse struct {
	Event EventInfo `json:"event"`
}
