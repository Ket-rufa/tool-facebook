package actiontest

import "time"

// Reaction type enum
const (
	ReactionLike  = "like"
	ReactionLove  = "love"
	ReactionHaha  = "haha"
	ReactionWow   = "wow"
	ReactionSad   = "sad"
	ReactionCare  = "care"
	ReactionAngry = "angry"
)

var ValidReactions = map[string]bool{
	ReactionLike: true, ReactionLove: true, ReactionHaha: true,
	ReactionWow: true, ReactionSad: true, ReactionCare: true, ReactionAngry: true,
}

// Status enum
const (
	StatusSuccess         = "success"
	StatusValidationError = "validation_error"
	StatusSessionInvalid  = "session_invalid"
	StatusNotEnabled      = "not_enabled"
	StatusFailed          = "failed"
)

// LikePostRequest — payload chuẩn hóa đầy đủ từ frontend
type LikePostRequest struct {
	PostID       string `json:"post_id"`
	AccountID    string `json:"account_id"`
	ReactionType string `json:"reaction_type"`
	DryRun       bool   `json:"dry_run"`
	ActorSource  string `json:"actor_source"`
}

// LikePostResponse — response đầy đủ trả về frontend
type LikePostResponse struct {
	Success            bool      `json:"success"`
	Status             string    `json:"status"`
	Action             string    `json:"action"`
	PostID             string    `json:"post_id"`
	AccountID          string    `json:"account_id"`
	AccountDisplayName string    `json:"account_display_name"`
	ReactionType       string    `json:"reaction_type"`
	DryRun             bool      `json:"dry_run"`
	Message            string    `json:"message"`
	ExecutedAt         time.Time `json:"executed_at"`
}

// SessionStatus represents the current session status
type SessionStatus struct {
	IsActive bool   `json:"is_active"`
	Provider string `json:"provider"`
	Message  string `json:"message"`
}

// ActionLog — model đầy đủ cho log với account và reaction
type ActionLog struct {
	ID                 string    `json:"id"`
	ActionType         string    `json:"action_type"`
	AccountID          string    `json:"account_id"`
	AccountDisplayName string    `json:"account_display_name"`
	PostID             string    `json:"post_id"`
	ReactionType       string    `json:"reaction_type"`
	Status             string    `json:"status"`
	DryRun             bool      `json:"dry_run"`
	Message            string    `json:"message"`
	ExecutedAt         time.Time `json:"executed_at"`
}
