package storage

import (
	"betting/internal/domain"
	"time"

	"github.com/Rhymond/go-money"
	"github.com/google/uuid"
)

// Bet is storage representation of domain.Bet.
type Bet struct {
	ID             uuid.UUID
	Status         string
	SelectedSpaces []int
	Stake          *money.Money
	PlacedAt       time.Time
	SettledAt      *time.Time
	Win            bool
	Table          uuid.UUID
}

// AdaptBetsToDomain adapts multiple Bet to domain.Bet.
func AdaptBetsToDomain(bets []Bet) []domain.Bet {
	domainBets := make([]domain.Bet, len(bets))

	for i := range bets {
		domainBets[i] = AdaptBetToDomain(bets[i])
	}

	return domainBets
}

func AdaptBetsToStorage(bets []domain.Bet) []Bet {
	b := make([]Bet, len(bets))

	for i := range bets {
		b[i] = AdaptBetToStorage(bets[i])
	}

	return b
}

// AdaptBetToStorage adapts a single domain.Bet to Bet.
func AdaptBetToStorage(bet domain.Bet) Bet {
	return Bet{
		ID:             bet.ID,
		Status:         bet.Status.String(),
		SelectedSpaces: bet.SelectedSpaces,
		Stake:          bet.Stake,
		PlacedAt:       bet.PlacedAt,
		Win:            bet.Win,
		SettledAt:      bet.SettledAt,
		Table:          bet.Table,
	}
}

// AdaptBetToDomain adapts a single Bet to domain.Bet.
func AdaptBetToDomain(bet Bet) domain.Bet {
	return domain.Bet{
		ID:             bet.ID,
		Status:         domain.BetStatus(bet.Status),
		SelectedSpaces: bet.SelectedSpaces,
		Stake:          bet.Stake,
		PlacedAt:       bet.PlacedAt,
		SettledAt:      bet.SettledAt,
		Win:            bet.Win,
		Table:          bet.Table,
	}
}
