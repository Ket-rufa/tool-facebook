package accounts

import (
	"fmt"
	"time"
)

// SessionChecker kiểm tra trạng thái session của từng tài khoản
// [THẬT]: Cập nhật local store với kết quả
// [PLACEHOLDER]: Logic check HTTP/cookie thật sẽ được thêm sau khi có session token thật
type SessionChecker struct {
	store Store
}

func NewSessionChecker(store Store) *SessionChecker {
	return &SessionChecker{store: store}
}

// CheckOne kiểm tra session của một tài khoản theo ID
// [PLACEHOLDER]: Phần gọi HTTP kiểm tra cookie/token chưa triển khai.
// Hiện tại dùng logic xoay vòng trạng thái để demo UI.
func (c *SessionChecker) CheckOne(accountID string) (*SessionCheckResult, error) {
	acc, err := c.store.Get(accountID)
	if err != nil {
		return nil, fmt.Errorf("lỗi đọc tài khoản: %w", err)
	}
	if acc == nil {
		return nil, fmt.Errorf("không tìm thấy tài khoản: %s", accountID)
	}

	// [PLACEHOLDER] Kết quả mock - sẽ thay bằng HTTP check thật
	newStatus := c.mockCheckSession(acc.SessionStatus)
	checkedAt := time.Now().Format("2006-01-02 15:04:05")

	if err := c.store.UpdateSessionStatus(accountID, newStatus, checkedAt); err != nil {
		return nil, fmt.Errorf("lỗi cập nhật trạng thái: %w", err)
	}

	return &SessionCheckResult{
		AccountID:     accountID,
		SessionStatus: string(newStatus),
		Provider:      acc.Provider,
		CheckedAt:     checkedAt,
		Message:       c.statusMessage(newStatus),
	}, nil
}

// CheckAll kiểm tra tất cả tài khoản trong store
func (c *SessionChecker) CheckAll() ([]SessionCheckResult, error) {
	accounts, err := c.store.List()
	if err != nil {
		return nil, err
	}

	results := make([]SessionCheckResult, 0, len(accounts))
	for _, acc := range accounts {
		result, err := c.CheckOne(acc.ID)
		if err != nil {
			results = append(results, SessionCheckResult{
				AccountID:     acc.ID,
				SessionStatus: string(SessionUnknown),
				Provider:      acc.Provider,
				CheckedAt:     time.Now().Format("2006-01-02 15:04:05"),
				Message:       "Lỗi kiểm tra: " + err.Error(),
			})
			continue
		}
		results = append(results, *result)
	}
	return results, nil
}

// mockCheckSession xoay vòng trạng thái để demo UI - [PLACEHOLDER]
func (c *SessionChecker) mockCheckSession(current string) SessionStatus {
	switch SessionStatus(current) {
	case SessionUnchecked:
		return SessionActive
	case SessionActive:
		return SessionActive // giữ active khi đã active
	case SessionExpired:
		return SessionInvalid
	case SessionInvalid:
		return SessionUnknown
	default:
		return SessionActive
	}
}

func (c *SessionChecker) statusMessage(s SessionStatus) string {
	switch s {
	case SessionActive:
		return "Phiên đang hoạt động bình thường."
	case SessionExpired:
		return "Phiên đã hết hạn. Cần gắn lại."
	case SessionInvalid:
		return "Phiên không hợp lệ. Token đã bị thu hồi hoặc thay đổi."
	case SessionUnknown:
		return "Không thể xác định trạng thái phiên."
	default:
		return "Chưa kiểm tra."
	}
}
