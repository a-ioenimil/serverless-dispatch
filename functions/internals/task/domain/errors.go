package domain

import "errors"

var (
	ErrTaskTitleRequired  = errors.New("task title is required")
	ErrInvalidStatus      = errors.New("invalid task status")
	ErrDescriptionTooLong = errors.New("description cannot exceed 2048 characters")
)
