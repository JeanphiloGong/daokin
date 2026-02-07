package domain

import "time"

type DAOMembership struct {
	DAOID       string     `json:"dao_id"`
	Wallet      string     `json:"wallet"`
	JoinedAt    time.Time  `json:"joined_at"`
	LeftAt      *time.Time `json:"left_at"`
	LeaveReason string     `json:"leave_reason,omitempty"`
}
