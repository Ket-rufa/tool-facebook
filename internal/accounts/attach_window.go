package accounts

import (
	"fmt"
	"time"

	"github.com/pkg/browser"
)

// AttachWindowHandler xử lý việc mở cửa sổ đăng nhập riêng
// Wails v2 trên Windows chưa hỗ trợ thêm WebviewWindow độc lập,
// nên giải pháp là mở URL bằng trình duyệt hệ thống có thể cùng domain redirect
// để backend nhận callback (local HTTP callback approach).
// Phase này: mở browser ở URL gốc, flow manager xử lý state.

type AttachWindowHandler struct {
	// baseURL là URL callback local để nhận kết quả sau khi đăng nhập
	// Khi có backend HTTP server thật, đây sẽ là endpoint nhận redirect.
	baseCallbackURL string
}

func NewAttachWindowHandler(callbackURL string) *AttachWindowHandler {
	if callbackURL == "" {
		callbackURL = "http://localhost:37891/auth/callback" // Cổng dành riêng cho auth callback
	}
	return &AttachWindowHandler{
		baseCallbackURL: callbackURL,
	}
}

// OpenLoginWindow mở trình duyệt hệ thống để người dùng đăng nhập
// flowID được truyền vào URL để backend nhận diện được callback thuộc về flow nào
// [THẬT]: Gọi browser.OpenURL thật để mở trình duyệt
// [PLACEHOLDER]: URL đích hiện vẫn chưa phải Facebook auth endpoint thật
func (h *AttachWindowHandler) OpenLoginWindow(flowID string) error {
	// Build URL với flow_id để sau này callback service biết hoàn tất flow nào
	loginURL := fmt.Sprintf(
		"https://www.facebook.com/login?flow_id=%s&redirect_uri=%s",
		flowID,
		h.baseCallbackURL,
	)
	return browser.OpenURL(loginURL)
}

// BuildMockAccountPreview tạo account preview với tên tự động theo số thứ tự
// [PLACEHOLDER]: Sẽ thay bằng dữ liệu từ OAuth callback thật
func BuildMockAccountPreview(flowID string, cloneIndex int) *AccountProfile {
	uid := "100" + fmt.Sprintf("%08d", time.Now().UnixMilli()%100000000)
	displayName := fmt.Sprintf("Clone %d", cloneIndex)

	return &AccountProfile{
		ID:            "acc_" + uid,
		AccountID:     uid,
		DisplayName:   displayName,
		Avatar:        fmt.Sprintf("https://ui-avatars.com/api/?name=C%d&background=dbeafe&color=1d4ed8", cloneIndex),
		AccountType:   string(AccountTypeProfile),
		Provider:      "Facebook (Mock)",
		SessionID:     "sess_" + uid,
		SessionStatus: string(SessionActive),
		AttachedAt:    time.Now().Format(time.RFC3339),
		LastCheckedAt: time.Now().Format("2006-01-02 15:04:05"),
		IsDefault:     false,
		Note:          "Phiên được gắn qua cửa sổ đăng nhập riêng. Chưa có xác thực OAuth thật.",
	}
}
