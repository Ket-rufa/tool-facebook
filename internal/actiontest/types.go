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

// ActionErrorCode represents detailed logical errors
type ActionErrorCode string

const (
	ErrIDMismatch      ActionErrorCode = "ERR_ID_MISMATCH"
	ErrSessionMissing  ActionErrorCode = "ERR_SESSION_CONTEXT_MISSING"
	ErrActorMismatch   ActionErrorCode = "ERR_ACTOR_MISMATCH"
	ErrStaleContext    ActionErrorCode = "ERR_STALE_TRACKING_CONTEXT"
	ErrGraphqlFB       ActionErrorCode = "ERR_GRAPHQL_FB_ERROR"
	ErrMalformed       ActionErrorCode = "ERR_MALFORMED_REQUEST"
)

// LikePostRequest — payload chuan hoa day du tu frontend
// doc_id KHONG ton tai trong struct nay — da xoa hoan toan khoi runtime path (HUONG A)
type LikePostRequest struct {
	PostID       string `json:"post_id"`
	AccountID    string `json:"account_id"`
	ReactionType string `json:"reaction_type"`
	ReactionID   string `json:"reaction_id"`
	DryRun       bool   `json:"dry_run"`
	ActorSource  string `json:"actor_source"`
}

// LikePostResponse — response đầy đủ trả về frontend
type LikePostResponse struct {
	Success            bool            `json:"success"`
	Status             string          `json:"status"`
	Action             string          `json:"action"`
	PostID             string          `json:"post_id"`
	AccountID          string          `json:"account_id"`
	AccountDisplayName string          `json:"account_display_name"`
	ReactionType       string          `json:"reaction_type"`
	ReactionID         string          `json:"reaction_id"`
	DryRun             bool            `json:"dry_run"`
	Message            string          `json:"message"`
	ErrorCode          ActionErrorCode `json:"error_code,omitempty"`
	ExecutedAt         time.Time       `json:"executed_at"`
}

// CommentPostRequest — payload cho bình luận
type CommentPostRequest struct {
	PostID      string `json:"post_id"`
	AccountID   string `json:"account_id"`
	CommentText string `json:"comment_text"`
	DryRun      bool   `json:"dry_run"`
	ActorSource string `json:"actor_source"`
}

// CommentPostResponse — response từ việc bình luận
type CommentPostResponse struct {
	Success            bool      `json:"success"`
	Status             string    `json:"status"`
	Action             string    `json:"action"`
	PostID             string    `json:"post_id"`
	AccountID          string    `json:"account_id"`
	AccountDisplayName string    `json:"account_display_name"`
	CommentText        string    `json:"comment_text"`
	DryRun             bool      `json:"dry_run"`
	Message            string    `json:"message"`
	ExecutedAt         time.Time `json:"executed_at"`
}

// CreatePostRequest — payload cho đăng bài
type CreatePostRequest struct {
	AccountID   string   `json:"account_id"`
	PostText    string   `json:"post_text"`
	ImagePaths  []string `json:"image_paths,omitempty"` // Dành cho sau này
	DryRun      bool     `json:"dry_run"`
	ActorSource string   `json:"actor_source"`
}

// CreatePostResponse — response từ việc đăng bài
type CreatePostResponse struct {
	Success            bool      `json:"success"`
	Status             string    `json:"status"`
	Action             string    `json:"action"`
	AccountID          string    `json:"account_id"`
	AccountDisplayName string    `json:"account_display_name"`
	PostText           string    `json:"post_text"`
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
	CommentText        string    `json:"comment_text,omitempty"`
	PostText           string    `json:"post_text,omitempty"`
	ReactionType       string    `json:"reaction_type,omitempty"`
	ReactionID         string    `json:"reaction_id,omitempty"`
	Status             string    `json:"status"`
	DryRun             bool      `json:"dry_run"`
	Message            string    `json:"message"`
	ExecutedAt         time.Time `json:"executed_at"`
}
