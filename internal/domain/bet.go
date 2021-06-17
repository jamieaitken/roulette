package domain

import (
	"time"

	"github.com/Rhymond/go-money"
	"github.com/google/uuid"
)

// Bet represents an individuals single pot for a single Table.
type Bet struct {
	ID             uuid.UUID
	Status         BetStatus
	SelectedSpaces []int
	Stake          *money.Money
	PlacedAt       time.Time
	SettledAt      *time.Time
	Win            bool
	Table          uuid.UUID
}

type BetStatus string

func (b BetStatus) String() string {
	return string(b)
}

var (
	Unsettled BetStatus = "unsettled"
	Live      BetStatus = "live"
	Settled   BetStatus = "settled"
)
