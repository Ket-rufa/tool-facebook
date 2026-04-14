package actiontest

type SessionHandler struct{}

func NewSessionHandler() *SessionHandler {
	return &SessionHandler{}
}

// GetSessionStatus returns a mock valid session
func (h *SessionHandler) GetSessionStatus() SessionStatus {
	return SessionStatus{
		IsActive: true,
		Provider: "local",
		Message:  "Phiên hoạt động",
	}
}
