package eventparam

type GetAllEventRequest struct {
	UserID uint `json:"user_id"`
}

type GetAllEventResponse struct {
	Events []EventInfo `json:"events"`
}
