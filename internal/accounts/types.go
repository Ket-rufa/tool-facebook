package accounts

import "time"

// SessionStatus là enum trạng thái phiên đăng nhập
type SessionStatus string

const (
	SessionActive   SessionStatus = "active"
	SessionExpired  SessionStatus = "expired"
	SessionInvalid  SessionStatus = "invalid"
	SessionUnknown  SessionStatus = "unknown"
	SessionUnchecked SessionStatus = "unchecked"
)

// AccountType loại tài khoản
type AccountType string

const (
	AccountTypeProfile AccountType = "Profile"
	AccountTypePage    AccountType = "Page"
	AccountTypeTest    AccountType = "Test"
)

// AttachFlowState trạng thái của flow attach tài khoản
type AttachFlowState string

const (
	FlowCreated       AttachFlowState = "created"
	FlowWaitingAuth   AttachFlowState = "waiting_auth"
	FlowAuthenticated AttachFlowState = "authenticated"
	FlowCompleted     AttachFlowState = "completed"
	FlowFailed        AttachFlowState = "failed"
	FlowCancelled     AttachFlowState = "cancelled"
)

// AccountProfile là thông tin tài khoản đã gắn session
type AccountProfile struct {
	ID            string `json:"id"`
	AccountID     string `json:"accountId"`
	DisplayName   string `json:"displayName"`
	Avatar        string `json:"avatar"`
	AccountType   string `json:"accountType"`
	Provider      string `json:"provider"`
	SessionID     string `json:"sessionId"`
	SessionStatus string `json:"sessionStatus"`
	AttachedAt    string `json:"attachedAt"`
	LastCheckedAt string `json:"lastCheckedAt"`
	IsDefault     bool   `json:"isDefault"`
	Note          string `json:"note"`
}

// AttachFlow đại diện cho một phiên gắn tài khoản đang diễn ra
type AttachFlow struct {
	FlowID        string          `json:"flowId"`
	State         AttachFlowState `json:"state"`
	Message       string          `json:"message"`
	AccountPreview *AccountProfile `json:"accountPreview,omitempty"`
	CreatedAt     time.Time       `json:"createdAt"`
	UpdatedAt     time.Time       `json:"updatedAt"`
}

// AttachFlowStatusResponse là response trả về cho frontend
type AttachFlowStatusResponse struct {
	FlowID         string          `json:"flowId"`
	State          string          `json:"state"`
	Message        string          `json:"message"`
	AccountPreview *AccountProfile `json:"accountPreview,omitempty"`
	UpdatedAt      string          `json:"updatedAt"`
}

// SessionCheckResult kết quả kiểm tra session
type SessionCheckResult struct {
	AccountID     string `json:"accountId"`
	SessionStatus string `json:"sessionStatus"`
	Provider      string `json:"provider"`
	CheckedAt     string `json:"checkedAt"`
	Message       string `json:"message"`
}
