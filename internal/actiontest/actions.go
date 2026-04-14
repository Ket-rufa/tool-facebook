package actiontest

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/yourname/tool-facebook/internal/accounts"
)

type ActionHandler struct {
	sessionHandler *SessionHandler
	accountStore   accounts.Store
	configStore    *ConfigStore
}

// NewActionHandler tạo handler không có account store (backward compat)
func NewActionHandler(sessionHandler *SessionHandler) *ActionHandler {
	return &ActionHandler{
		sessionHandler: sessionHandler,
		accountStore:   nil,
		configStore:    nil,
	}
}

// NewActionHandlerWithStore tạo handler có account store thật
func NewActionHandlerWithStore(sessionHandler *SessionHandler, store accounts.Store) *ActionHandler {
	return &ActionHandler{
		sessionHandler: sessionHandler,
		accountStore:   store,
		configStore:    nil,
	}
}

// NewActionHandlerWithStoreAndConfig tao handler co account store va config store.
func NewActionHandlerWithStoreAndConfig(sessionHandler *SessionHandler, store accounts.Store, cfg *ConfigStore) *ActionHandler {
	return &ActionHandler{
		sessionHandler: sessionHandler,
		accountStore:   store,
		configStore:    cfg,
	}
}

func (h *ActionHandler) LikePost(req LikePostRequest) LikePostResponse {
	now := time.Now()

	// ── 1. Validate post_id ────────────────────────────────────────────
	postID := strings.TrimSpace(req.PostID)
	if postID == "" {
		resp := LikePostResponse{
			Success: false, Status: StatusValidationError,
			Action: req.ReactionType, PostID: "",
			AccountID: req.AccountID, ReactionType: req.ReactionType,
			ReactionID: req.ReactionID,
			DryRun:     req.DryRun, Message: "Post ID không được để trống.", ExecutedAt: now,
		}
		return resp
	}

	// ── 2. Validate reaction_type ──────────────────────────────────────
	reactionType := strings.ToLower(strings.TrimSpace(req.ReactionType))
	if reactionType == "" {
		reactionType = ReactionLike
	}
	if !ValidReactions[reactionType] {
		resp := LikePostResponse{
			Success: false, Status: StatusValidationError,
			Action: reactionType, PostID: postID,
			AccountID: req.AccountID, ReactionType: reactionType,
			ReactionID: req.ReactionID,
			DryRun:     req.DryRun, Message: fmt.Sprintf("Loại cảm xúc không hợp lệ: %s", reactionType), ExecutedAt: now,
		}
		return resp
	}

	// ── 3. Validate & load account từ local store ──────────────────────
	accountID := strings.TrimSpace(req.AccountID)
	displayName := accountID
	if accountID == "" {
		resp := LikePostResponse{
			Success: false, Status: StatusValidationError,
			Action: reactionType, PostID: postID,
			AccountID: "", ReactionType: reactionType,
			ReactionID: req.ReactionID,
			DryRun:     req.DryRun, Message: "Chưa chọn tài khoản thực thi.", ExecutedAt: now,
		}
		return resp
	}

	// Kiểm tra account trong local store nếu có
	sessionValid := true
	sessionStatusStr := "active"
	if h.accountStore != nil {
		acc, err := h.accountStore.Get(accountID)
		if err != nil || acc == nil {
			resp := LikePostResponse{
				Success: false, Status: StatusValidationError,
				Action: reactionType, PostID: postID,
				AccountID: accountID, ReactionType: reactionType,
				ReactionID: req.ReactionID,
				DryRun:     req.DryRun,
				Message:    fmt.Sprintf("Tài khoản '%s' không tìm thấy trong local store.", accountID),
				ExecutedAt: now,
			}
			AddLog(toLog(resp, "Tài khoản không tồn tại"))
			return resp
		}
		displayName = acc.DisplayName
		sessionStatusStr = acc.SessionStatus

		// session invalid/expired -> không thể chạy ngay cả dry run
		if acc.SessionStatus == string(accounts.SessionInvalid) ||
			acc.SessionStatus == string(accounts.SessionExpired) ||
			acc.SessionStatus == "error" {
			sessionValid = false
		}
	} else {
		// Fallback: dùng session handler cũ nếu không có store
		status := h.sessionHandler.GetSessionStatus()
		if !status.IsActive {
			sessionValid = false
		}
	}

	// ── 4. Session phải hợp lệ ngay cả khi dry_run ────────────────────
	if !sessionValid {
		resp := LikePostResponse{
			Success: false, Status: StatusSessionInvalid,
			Action: reactionType, PostID: postID,
			AccountID: accountID, AccountDisplayName: displayName, ReactionType: reactionType,
			ReactionID: req.ReactionID,
			DryRun:     req.DryRun,
			Message:    fmt.Sprintf("Phiên của tài khoản '%s' không hợp lệ (trạng thái: %s). Vui lòng kiểm tra lại phiên.", displayName, sessionStatusStr),
			ExecutedAt: now,
		}
		AddLog(toLog(resp, resp.Message))
		return resp
	}

	// ── 5. Dry Run ─────────────────────────────────────────────────────
	if req.DryRun {
		sleepMs := 500 + rand.Intn(701)
		time.Sleep(time.Duration(sleepMs) * time.Millisecond)

		resp := LikePostResponse{
			Success: true, Status: StatusSuccess,
			Action: reactionType, PostID: postID,
			AccountID: accountID, AccountDisplayName: displayName, ReactionType: reactionType,
			ReactionID: req.ReactionID,
			DryRun:     true,
			Message:    fmt.Sprintf("[DRY RUN] Mô phỏng gửi '%s' (ID: %s) thành công — Tài khoản: %s — Bài viết: %s", reactionType, req.ReactionID, displayName, postID),
			ExecutedAt: time.Now(),
		}
		AddLog(toLog(resp, resp.Message))
		return resp
	}

	// ── 6. Real Run ───────────────────────────────────────────────────
	var cookie, fbDtsg, accUID string
	if h.accountStore != nil {
		acc, _ := h.accountStore.Get(accountID)
		if acc != nil {
			cookie = acc.Cookie
			fbDtsg = acc.FbDtsg
			accUID = acc.AccountID
		}
	}

	// ── 6a. Pre-flight log — strategy is resolved inside executor ────────
	log.Printf("[RealRun] === Pre-flight ===")
	log.Printf("[RealRun]   strategy   = auto (comet/doc_id fallback)")
	log.Printf("[RealRun]   postID     = %q", postID)
	log.Printf("[RealRun]   reactionID = %q", req.ReactionID)
	log.Printf("[RealRun]   actorID    = %q", accUID)
	configuredDocID := h.GetReactDocID()
	if configuredDocID != "" {
		log.Printf("[RealRun]   doc_id     = configured in app settings")
	} else {
		log.Printf("[RealRun]   doc_id     = from env FB_REACT_DOC_ID (if set), else inline query")
	}
	fbResp, errCode, err := ExecuteFacebookReaction(cookie, fbDtsg, accUID, postID, req.ReactionID, configuredDocID)

	if err != nil {
		errorStatus := StatusValidationError
		if errCode != "" {
			errorStatus = string(errCode)
		}

		resp := LikePostResponse{
			Success: false, Status: errorStatus,
			Action: reactionType, PostID: postID,
			AccountID: accountID, AccountDisplayName: displayName, ReactionType: reactionType,
			ReactionID: req.ReactionID,
			DryRun:     false,
			ErrorCode:  errCode,
			Message:    fmt.Sprintf("[%s] %v", errCode, err),
			ExecutedAt: time.Now(),
		}
		AddLog(toLog(resp, resp.Message))
		return resp
	}

	resp := LikePostResponse{
		Success: true, Status: StatusSuccess,
		Action: reactionType, PostID: postID,
		AccountID: accountID, AccountDisplayName: displayName, ReactionType: reactionType,
		ReactionID: req.ReactionID,
		DryRun:     false,
		Message:    "Tha cam xuc that thanh cong. Phan hoi FB: " + fbResp,
		ExecutedAt: time.Now(),
	}
	AddLog(toLog(resp, resp.Message))
	return resp
}

// GetActionLogs exposes the logs to Wails frontend
func (h *ActionHandler) GetActionLogs() []ActionLog {
	return GetLogs()
}

// GetReactDocID lay doc_id tu local config, fallback ve env.
func (h *ActionHandler) GetReactDocID() string {
	if h.configStore != nil {
		if v := h.configStore.GetReactDocID(); v != "" {
			return v
		}
	}
	return strings.TrimSpace(os.Getenv("FB_REACT_DOC_ID"))
}

// UpdateReactDocID cap nhat doc_id local.
// Truyen chuoi rong de xoa doc_id da luu.
func (h *ActionHandler) UpdateReactDocID(docID string) (string, error) {
	if h.configStore == nil {
		return "", errors.New("config store chua duoc khoi tao")
	}
	if err := h.configStore.SetReactDocID(docID); err != nil {
		return "", err
	}
	return h.configStore.GetReactDocID(), nil
}

// toLog chuyển response thành ActionLog entry
func toLog(resp LikePostResponse, msg string) ActionLog {
	return ActionLog{
		ActionType:         resp.Action,
		AccountID:          resp.AccountID,
		AccountDisplayName: resp.AccountDisplayName,
		PostID:             resp.PostID,
		ReactionType:       resp.ReactionType,
		ReactionID:         resp.ReactionID,
		Status:             resp.Status,
		DryRun:             resp.DryRun,
		Message:            msg,
		ExecutedAt:         resp.ExecutedAt,
	}
}
