package commands

type LoginRequest struct{}

func NewLoginRequest() *LoginRequest {
	return &LoginRequest{}
}

func (c *LoginRequest) Serialize() string {
	return "<msg type=\"login_req\"><params protocol=\"TCP\" port=\"8000\"/></msg>"
}

func (c *LoginRequest) ToHuman() string {
	return "LoginRequest"
}
