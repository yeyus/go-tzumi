package commands

type SignalStatusRequest struct{}

func NewSignalStatusRequest() *SignalStatusRequest {
	return &SignalStatusRequest{}
}

func (c *SignalStatusRequest) Serialize() string {
	return "<msg type=\"signal_status_req\"></msg>"
}

func (c *SignalStatusRequest) ToHuman() string {
	return "SignalStatusRequest"
}
