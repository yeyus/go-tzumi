package commands

type CheckLockRequest struct{}

func NewCheckLockRequest() *CheckLockRequest {
	return &CheckLockRequest{}
}

func (c *CheckLockRequest) Serialize() string {
	return "<msg type=\"Check_Lock_req\"></msg>"
}

func (c *CheckLockRequest) ToHuman() string {
	return "CheckLockRequest"
}
