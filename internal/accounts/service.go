package accounts

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/yourname/tool-facebook/internal/fbdata"
)

// AccountService là service chính - được Wails bind vào frontend
// Đây là điểm duy nhất frontend gọi vào backend accounts module.
type AccountService struct {
	store            Store
	fbStore          fbdata.Store
	flowManager      *AttachFlowManager
	windowHandler    *AttachWindowHandler
	sessionCheck     *SessionChecker
	automator        *LoginAutomator
	requestAutomator *RequestAutomator
}

func (s *AccountService) syncToLegacyStore(profile AccountProfile) {
	if s.fbStore == nil {
		return
	}

	fbData := fbdata.FBAccountData{
		Info: fbdata.FBInfo{
			Name:     profile.DisplayName,
			Cookie:   profile.Cookie,
			FbDtsg:   profile.FbDtsg,
		},
	}
	// Lưu dưới cả 2 khóa để đảm bảo các module cũ (dùng UID) và mới (dùng acc_UID) đều tìm thấy
	s.fbStore.Save(profile.AccountID, fbData)
	s.fbStore.Save(profile.ID, fbData)
}

// NewAccountService tạo service với data dir, tự tạo store
func NewAccountService(dataDir string) (*AccountService, error) {
	store, err := NewJSONStore(dataDir)
	if err != nil {
		return nil, fmt.Errorf("lỗi khởi tạo store: %w", err)
	}
	// Khởi tạo cả FB Store để đồng bộ dữ liệu
	fbStore, _ := fbdata.NewJSONStore(dataDir)
	return NewAccountServiceWithStore(store, fbStore), nil
}

// NewAccountServiceWithStore tạo service với store đã khởi tạo sẵn (để share store)
func NewAccountServiceWithStore(store *JSONStore, fbStore fbdata.Store) *AccountService {
	return &AccountService{
		store:            store,
		fbStore:          fbStore,
		flowManager:      NewAttachFlowManager(),
		windowHandler:    NewAttachWindowHandler(""),
		sessionCheck:     NewSessionChecker(store),
		automator:        NewLoginAutomator(),
		requestAutomator: NewRequestAutomator(),
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
	s.syncToLegacyStore(payload)
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

	// Chuyển sang waiting_auth
	s.flowManager.SetWaitingAuth(flow.FlowID)

	// Chạy automation trong goroutine để không block Wails
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()

		result, err := s.automator.StartLogin(ctx)
		if err != nil {
			fmt.Printf("[LOGIN] Lỗi automation: %v\n", err)
			s.flowManager.SetFailed(flow.FlowID, "Lỗi đăng nhập: "+err.Error())
			return
		}

		// Tạo preview từ dữ liệu thật
		preview := &AccountProfile{
			ID:            "acc_" + result.UID, // Dùng UID làm prefix ID luôn
			AccountID:     result.UID,
			DisplayName:   result.Name,
			Avatar:        fmt.Sprintf("https://ui-avatars.com/api/?name=%s&background=dbeafe&color=1d4ed8", url.QueryEscape(result.Name)),
			AccountType:   string(AccountTypeProfile),
			Provider:      "Facebook",
			SessionID:     "sess_" + result.UID,
			SessionStatus: string(SessionActive),
			Cookie:        result.Cookie,
			FbDtsg:        result.Dtsg,
			AttachedAt:    time.Now().Format(time.RFC3339),
			LastCheckedAt: time.Now().Format("2006-01-02 15:04:05"),
			Note:          "Được thêm tự động qua trình duyệt.",
		}

		// Cập nhật flow status sang Authenticated để frontend hiển thị Confirm
		s.flowManager.SetAuthenticated(flow.FlowID, preview)
		fmt.Printf("[LOGIN] Thành công: %s (%s)\n", result.Name, result.UID)
	}()

	resp, _ := s.flowManager.GetStatus(flow.FlowID)
	return *resp
}

// GetCookieFromCredentials thực hiện đăng nhập ngầm và chỉ trả về chuỗi Cookie (không tạo flow)
func (s *AccountService) GetCookieFromCredentials(email, password, twoFactorKey string) (string, error) {
	fmt.Printf("[SERVICE] Đang lấy Cookie cho: %s\n", email)
	result, err := s.requestAutomator.Login(email, password, twoFactorKey)
	if err != nil {
		return "", err
	}
	if !result.Success {
		return "", fmt.Errorf(result.Message)
	}
	return result.Cookie, nil
}

// LoginByRequest thực hiện đăng nhập qua trình duyệt tự động (Chrome)
func (s *AccountService) LoginByRequest(email, password, twoFactorKey string) AttachFlowStatusResponse {
	// Khởi tạo flow mới
	flow, err := s.flowManager.Start()
	if err != nil {
		return AttachFlowStatusResponse{State: string(FlowFailed), Message: err.Error()}
	}

	s.flowManager.SetWaitingAuth(flow.FlowID)

	// Chạy automation trong goroutine (hiện trình duyệt)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()

		result, err := s.automator.AutomatedLogin(ctx, email, password, twoFactorKey)
		if err != nil {
			s.flowManager.SetFailed(flow.FlowID, "Lỗi đăng nhập: "+err.Error())
			return
		}

		// Tạo profile
		preview := &AccountProfile{
			ID:            "acc_" + result.UID,
			AccountID:     result.UID,
			DisplayName:   result.Name,
			Avatar:        fmt.Sprintf("https://ui-avatars.com/api/?name=%s&background=dbeafe&color=1d4ed8", url.QueryEscape(result.Name)),
			AccountType:   string(AccountTypeProfile),
			Provider:      "Facebook (Auto)",
			SessionID:     "sess_" + result.UID,
			SessionStatus: string(SessionActive),
			Cookie:        result.Cookie,
			FbDtsg:        result.Dtsg,
			AttachedAt:    time.Now().Format(time.RFC3339),
			LastCheckedAt: time.Now().Format("2006-01-02 15:04:05"),
			Note:          "Được thêm tự động qua trình duyệt.",
		}

		s.flowManager.SetAuthenticated(flow.FlowID, preview)
	}()

	resp, _ := s.flowManager.GetStatus(flow.FlowID)
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

	// Đồng bộ sang facebook_data.json
	s.syncToLegacyStore(profile)

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
