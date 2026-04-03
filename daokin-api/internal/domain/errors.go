package domain

import "errors"

var (
	ErrNotFound         = errors.New("not found")
	ErrConflict         = errors.New("conflict")
	ErrInvalidArgument  = errors.New("invalid argument")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrForbidden        = errors.New("forbidden")
	ErrRuleViolation    = errors.New("rule violation")
	ErrChallengeInvalid = errors.New("challenge invalid")
)
