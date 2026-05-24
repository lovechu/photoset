package domain

import "errors"

// Domain errors for community module
var (
	// Post errors
	ErrTitleRequired     = errors.New("title is required")
	ErrTitleTooLong      = errors.New("title cannot exceed 200 characters")
	ErrContentRequired    = errors.New("content is required")
	ErrContentTooLong     = errors.New("content cannot exceed 5000 characters")
	ErrInvalidCategory    = errors.New("invalid category")
	ErrInvalidVisibility  = errors.New("invalid visibility")
	ErrPostNotFound       = errors.New("post not found")
	ErrPostDeleted        = errors.New("post has been deleted")

	// Reply errors
	ErrReplyContentRequired  = errors.New("reply content is required")
	ErrReplyContentTooLong   = errors.New("reply content cannot exceed 2000 characters")
	ErrReplyNotFound         = errors.New("reply not found")

	// Like errors
	ErrAlreadyLiked          = errors.New("already liked")
	ErrNotLiked             = errors.New("not liked yet")

	// Report errors
	ErrReportReasonRequired  = errors.New("report reason is required")
	ErrReportNotFound        = errors.New("report not found")

	// Point errors
	ErrDailyLimitReached    = errors.New("daily point limit reached")
	ErrInvalidPoints         = errors.New("invalid points value")

	// Sensitive word errors
	ErrSensitiveWordExists   = errors.New("sensitive word already exists")
	ErrSensitiveWordNotFound = errors.New("sensitive word not found")

	// Permission errors
	ErrPermissionDenied      = errors.New("permission denied")
	ErrLoginRequired         = errors.New("login required")
)
