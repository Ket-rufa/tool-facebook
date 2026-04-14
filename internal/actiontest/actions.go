package actiontest

import (
	"math/rand"
	"strings"
	"time"
)

type ActionHandler struct {
	sessionHandler *SessionHandler
}

func NewActionHandler(sessionHandler *SessionHandler) *ActionHandler {
	return &ActionHandler{
		sessionHandler: sessionHandler,
	}
}

func (h *ActionHandler) LikePost(req LikePostRequest) LikePostResponse {
	// 1. Check Session Status
	status := h.sessionHandler.GetSessionStatus()
	if !status.IsActive {
		resp := LikePostResponse{
			Success:    false,
			Status:     StatusSessionInvalid,
			Action:     "like",
			PostID:     req.PostID,
			DryRun:     req.DryRun,
			Message:    "Phiên không hợp lệ",
			ExecutedAt: time.Now(),
		}
		AddLog(ActionLog{
			ActionType: "like",
			PostID:     req.PostID,
			Status:     resp.Status,
			DryRun:     req.DryRun,
			Message:    resp.Message,
			ExecutedAt: resp.ExecutedAt,
		})
		return resp
	}

	// 2. Validate PostID
	postID := strings.TrimSpace(req.PostID)
	if postID == "" {
		resp := LikePostResponse{
			Success:    false,
			Status:     StatusValidationError,
			Action:     "like",
			PostID:     "",
			DryRun:     req.DryRun,
			Message:    "Post ID không được để trống",
			ExecutedAt: time.Now(),
		}
		return resp // Usually we might not log validation errors for empty fields, but optionally can.
	}

	// 3. Handle by mode
	if req.DryRun {
		// Simulate network latency (500 - 1200ms)
		sleepMs := 500 + rand.Intn(701)
		time.Sleep(time.Duration(sleepMs) * time.Millisecond)

		resp := LikePostResponse{
			Success:    true,
			Status:     StatusSuccess,
			Action:     "like",
			PostID:     postID,
			DryRun:     true,
			Message:    "Đã mô phỏng gửi like thành công",
			ExecutedAt: time.Now(),
		}
		
		AddLog(ActionLog{
			ActionType: "like",
			PostID:     postID,
			Status:     resp.Status,
			DryRun:     true,
			Message:    resp.Message,
			ExecutedAt: resp.ExecutedAt,
		})
		return resp
	}

	// Real run (not enabled)
	resp := LikePostResponse{
		Success:    false,
		Status:     StatusNotEnabled,
		Action:     "like",
		PostID:     postID,
		DryRun:     false,
		Message:    "Chế độ Real Run chưa được kích hoạt",
		ExecutedAt: time.Now(),
	}

	AddLog(ActionLog{
		ActionType: "like",
		PostID:     postID,
		Status:     resp.Status,
		DryRun:     false,
		Message:    resp.Message,
		ExecutedAt: resp.ExecutedAt,
	})

	return resp
}

// GetActionLogs exposes the logs to Wails frontend
func (h *ActionHandler) GetActionLogs() []ActionLog {
	return GetLogs()
}
