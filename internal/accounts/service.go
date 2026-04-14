package accounts

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

// AccountService là service chính - được Wails bind vào frontend
// Đây là điểm duy nhất frontend gọi vào backend accounts module.
type AccountService struct {
	store         Store
	flowManager   *AttachFlowManager
	windowHandler *AttachWindowHandler
	sessionCheck  *SessionChecker
}

// NewAccountService tạo service với data dir, tự tạo store
func NewAccountService(dataDir string) (*AccountService, error) {
	store, err := NewJSONStore(dataDir)
	if err != nil {
		return nil, fmt.Errorf("lỗi khởi tạo store: %w", err)
	}
	return NewAccountServiceWithStore(store), nil
}

// NewAccountServiceWithStore tạo service với store đã khởi tạo sẵn (để share store)
func NewAccountServiceWithStore(store *JSONStore) *AccountService {
	return &AccountService{
		store:         store,
		flowManager:   NewAttachFlowManager(),
		windowHandler: NewAttachWindowHandler(""),
		sessionCheck:  NewSessionChecker(store),
	}
}

// internalStore trả về Store để package nội bộ dùng (unexported, không generate Wails binding)
func (s *AccountService) internalStore() Store {
	return s.store
}

// ─────────────────────────────────────────────────
// Account Management Commands
// ─────────────────────────────────────────────────

// ListAccounts trả về toàn bộ tài khoản đã lưu
func (s *AccountService) ListAccounts() []AccountProfile {
	accounts, _ := s.store.List()
	if accounts == nil {
		return []AccountProfile{}
	}
	return accounts
}

// RemoveAccount gỡ tài khoản ra khỏi local store
func (s *AccountService) RemoveAccount(id string) []AccountProfile {
	s.store.Remove(id)
	return s.ListAccounts()
}

// SetDefaultAccount đặt một tài khoản là mặc định (các tài khoản khác tự unset)
func (s *AccountService) SetDefaultAccount(id string) []AccountProfile {
	s.store.SetDefault(id)
	return s.ListAccounts()
}

// GetDefaultAccount trả về tài khoản đang được đặt làm mặc định
func (s *AccountService) GetDefaultAccount() *AccountProfile {
	acc, _ := s.store.GetDefault()
	return acc
}

// AddAccountProfile thêm trực tiếp (dùng cho fallback/manual case khi OAuth chưa có)
func (s *AccountService) AddAccountProfile(payload AccountProfile) []AccountProfile {
	if payload.ID == "" {
		payload.ID = "acc_" + fmt.Sprintf("%d", time.Now().UnixMilli())
	}
	if payload.AttachedAt == "" {
		payload.AttachedAt = time.Now().Format(time.RFC3339)
	}
	if payload.SessionStatus == "" {
		payload.SessionStatus = string(SessionActive)
	}
	if payload.LastCheckedAt == "" {
		payload.LastCheckedAt = time.Now().Format("2006-01-02 15:04:05")
	}
	s.store.Add(payload)
	return s.ListAccounts()
}

// ─────────────────────────────────────────────────
// Attach Flow Commands
// ─────────────────────────────────────────────────

// StartAccountAttachFlow khởi tạo flow gắn tài khoản mới
// Mở cửa sổ trình duyệt thật để đăng nhập
func (s *AccountService) StartAccountAttachFlow() AttachFlowStatusResponse {
	flow, err := s.flowManager.Start()
	if err != nil {
		return AttachFlowStatusResponse{
			State:   string(FlowFailed),
			Message: "Không thể khởi tạo flow: " + err.Error(),
		}
	}

	// Chuyển sang waiting_auth trước khi mở cửa sổ
	s.flowManager.SetWaitingAuth(flow.FlowID)

	// [THẬT] Mở browser hệ thống với URL đăng nhập
	if err := s.windowHandler.OpenLoginWindow(flow.FlowID); err != nil {
		s.flowManager.Cancel(flow.FlowID)
		return AttachFlowStatusResponse{
			FlowID:  flow.FlowID,
			State:   string(FlowFailed),
			Message: "Không thể mở cửa sổ đăng nhập: " + err.Error(),
		}
	}

	// Đếm số account hiện có để đặt tên Clone N+1
	existing, _ := s.store.List()
	cloneIndex := len(existing) + 1

	// [PLACEHOLDER] Tự động tạo preview theo số thứ tự
	preview := BuildMockAccountPreview(flow.FlowID, cloneIndex)
	s.flowManager.SetAuthenticated(flow.FlowID, preview)

	resp, _ := s.flowManager.GetStatus(flow.FlowID)
	if resp == nil {
		return AttachFlowStatusResponse{
			FlowID:  flow.FlowID,
			State:   string(FlowWaitingAuth),
			Message: "Đang chờ xác thực...",
		}
	}
	return *resp
}

// GetAttachFlowStatus kiểm tra trạng thái của một flow đang chạy
func (s *AccountService) GetAttachFlowStatus(flowID string) AttachFlowStatusResponse {
	resp, err := s.flowManager.GetStatus(flowID)
	if err != nil {
		return AttachFlowStatusResponse{
			FlowID:  flowID,
			State:   string(FlowFailed),
			Message: err.Error(),
		}
	}
	return *resp
}

// CompleteAccountAttachFlow hoàn tất flow: lưu account vào store, trả list mới
// customName: tên do user nhập, nếu rỗng dùng tên từ preview
func (s *AccountService) CompleteAccountAttachFlow(flowID string, customName string) []AccountProfile {
	flow, err := s.flowManager.Complete(flowID)
	if err != nil || flow.AccountPreview == nil {
		return s.ListAccounts()
	}

	// Lưu account vào store
	profile := *flow.AccountPreview
	if customName != "" {
		profile.DisplayName = customName
		// Cập nhật avatar initials theo tên mới
		profile.Avatar = fmt.Sprintf("https://ui-avatars.com/api/?name=%s&background=dbeafe&color=1d4ed8", 
			url.QueryEscape(firstLetters(customName)))
	}
	profile.AttachedAt = time.Now().Format(time.RFC3339)
	profile.LastCheckedAt = time.Now().Format("2006-01-02 15:04:05")
	profile.SessionStatus = string(SessionActive)
	s.store.Add(profile)

	// Dọn sạch flow sau khi hoàn tất
	defer s.flowManager.Cleanup(flowID)

	return s.ListAccounts()
}

// CancelAccountAttachFlow hủy flow đang chạy, không tạo account rác
func (s *AccountService) CancelAccountAttachFlow(flowID string) {
	s.flowManager.Cancel(flowID)
	s.flowManager.Cleanup(flowID)
}

// firstLetters lấy các chữ cái đầu của tên để tạo avatar initials
func firstLetters(name string) string {
	parts := strings.Fields(name)
	if len(parts) == 0 {
		return "FB"
	}
	if len(parts) == 1 {
		if len(parts[0]) > 0 {
			return string([]rune(parts[0])[:1])
		}
		return "FB"
	}
	runes0 := []rune(parts[0])
	runesLast := []rune(parts[len(parts)-1])
	if len(runes0) > 0 && len(runesLast) > 0 {
		return string(runes0[:1]) + string(runesLast[:1])
	}
	return "FB"
}

// ─────────────────────────────────────────────────
// Session Check Commands
// ─────────────────────────────────────────────────

// CheckSessionStatus kiểm tra phiên của một tài khoản theo ID, trả list cập nhật
func (s *AccountService) CheckSessionStatus(id string) []AccountProfile {
	s.sessionCheck.CheckOne(id)
	return s.ListAccounts()
}

// CheckAllSessionStatuses kiểm tra tất cả tài khoản, trả list cập nhật
func (s *AccountService) CheckAllSessionStatuses() []AccountProfile {
	s.sessionCheck.CheckAll()
	return s.ListAccounts()
}
