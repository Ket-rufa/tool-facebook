package actiontest

import "time"

// Enum types for action statuses
const (
	StatusSuccess         = "success"
	StatusValidationError = "validation_error"
	StatusSessionInvalid  = "session_invalid"
	StatusNotEnabled      = "not_enabled"
	StatusFailed          = "failed"
)

// LikePostRequest is the minimal payload
type LikePostRequest struct {
	PostID      string `json:"post_id"`
	DryRun      bool   `json:"dry_run"`
	ActorSource string `json:"actor_source"`
}

// LikePostResponse is the minimal response
type LikePostResponse struct {
	Success    bool      `json:"success"`
	Status     string    `json:"status"`
	Action     string    `json:"action"`
	PostID     string    `json:"post_id"`
	DryRun     bool      `json:"dry_run"`
	Message    string    `json:"message"`
	ExecutedAt time.Time `json:"executed_at"`
}

// SessionStatus represents the current session status
type SessionStatus struct {
	IsActive bool   `json:"is_active"`
	Provider string `json:"provider"`
	Message  string `json:"message"`
}

// ActionLog is the in-memory minimal model for an action log
type ActionLog struct {
	ID         string    `json:"id"`
	ActionType string    `json:"action_type"`
	PostID     string    `json:"post_id"`
	Status     string    `json:"status"`
	DryRun     bool      `json:"dry_run"`
	Message    string    `json:"message"`
	ExecutedAt time.Time `json:"executed_at"`
}
