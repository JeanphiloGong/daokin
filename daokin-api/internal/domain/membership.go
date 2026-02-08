package domain

import (
	"strings"
	"time"
)

type DAOMembership struct {
	DAOID       string     `json:"dao_id"`
	Wallet      string     `json:"wallet"`
	JoinedAt    time.Time  `json:"joined_at"`
	LeftAt      *time.Time `json:"left_at"`
	LeaveReason string     `json:"leave_reason,omitempty"`
}

func NewDAOMembership(daoID, wallet string, joinedAt time.Time) (DAOMembership, error) {
	daoID = strings.TrimSpace(daoID)
	wallet = strings.TrimSpace(wallet)
	if daoID == "" || wallet == "" {
		return DAOMembership{}, ErrInvalidArgument
	}

	return DAOMembership{
		DAOID:    daoID,
		Wallet:   wallet,
		JoinedAt: joinedAt.UTC(),
	}, nil
}
