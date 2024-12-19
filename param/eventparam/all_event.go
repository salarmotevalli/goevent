package eventparam

type GetAllEventRequest struct{}

type GetAllEventResponse struct {
	Events []EventInfo `json:"events"`
}
