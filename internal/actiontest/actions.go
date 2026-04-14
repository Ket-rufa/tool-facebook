package actiontest

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/yourname/tool-facebook/internal/accounts"
)

type ActionHandler struct {
	sessionHandler *SessionHandler
	accountStore   accounts.Store
}

// NewActionHandler tạo handler không có account store (backward compat)
func NewActionHandler(sessionHandler *SessionHandler) *ActionHandler {
	return &ActionHandler{
		sessionHandler: sessionHandler,
		accountStore:   nil,
	}
}

// NewActionHandlerWithStore tạo handler có account store thật
func NewActionHandlerWithStore(sessionHandler *SessionHandler, store accounts.Store) *ActionHandler {
	return &ActionHandler{
		sessionHandler: sessionHandler,
		accountStore:   store,
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
			DryRun: req.DryRun, Message: "Post ID không được để trống.", ExecutedAt: now,
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
			DryRun: req.DryRun, Message: fmt.Sprintf("Loại cảm xúc không hợp lệ: %s", reactionType), ExecutedAt: now,
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
			DryRun: req.DryRun, Message: "Chưa chọn tài khoản thực thi.", ExecutedAt: now,
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
				DryRun: req.DryRun,
				Message: fmt.Sprintf("Tài khoản '%s' không tìm thấy trong local store.", accountID),
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
			DryRun: req.DryRun,
			Message: fmt.Sprintf("Phiên của tài khoản '%s' không hợp lệ (trạng thái: %s). Vui lòng kiểm tra lại phiên.", displayName, sessionStatusStr),
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
			DryRun: true,
			Message: fmt.Sprintf("[DRY RUN] Mô phỏng gửi '%s' thành công — Tài khoản: %s — Bài viết: %s", reactionType, displayName, postID),
			ExecutedAt: time.Now(),
		}
		AddLog(toLog(resp, resp.Message))
		return resp
	}

	// ── 6. Real Run — chưa bật executor ───────────────────────────────
	resp := LikePostResponse{
		Success: false, Status: StatusNotEnabled,
		Action: reactionType, PostID: postID,
		AccountID: accountID, AccountDisplayName: displayName, ReactionType: reactionType,
		DryRun: false,
		Message: "Chế độ Real Run chưa được kích hoạt. Executor thật chưa được triển khai.",
		ExecutedAt: time.Now(),
	}
	AddLog(toLog(resp, resp.Message))
	return resp
}

// GetActionLogs exposes the logs to Wails frontend
func (h *ActionHandler) GetActionLogs() []ActionLog {
	return GetLogs()
}

// toLog chuyển response thành ActionLog entry
func toLog(resp LikePostResponse, msg string) ActionLog {
	return ActionLog{
		ActionType:         resp.Action,
		AccountID:          resp.AccountID,
		AccountDisplayName: resp.AccountDisplayName,
		PostID:             resp.PostID,
		ReactionType:       resp.ReactionType,
		Status:             resp.Status,
		DryRun:             resp.DryRun,
		Message:            msg,
		ExecutedAt:         resp.ExecutedAt,
	}
}
