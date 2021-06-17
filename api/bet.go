package api

import (
	"betting/internal/domain"
	"time"

	"github.com/Rhymond/go-money"
	"github.com/google/uuid"
)

// BetRequest represents the required fields to create a bet.
type BetRequest struct {
	SelectedSpaces []int        `json:"selectedSpaces"`
	Stake          *money.Money `json:"stake"`
	Table          uuid.UUID    `json:"table"`
}

// BetResponse represents a Bet in responses to clients.
type BetResponse struct {
	ID        uuid.UUID        `json:"id"`
	PlacedAt  time.Time        `json:"placedAt"`
	Status    domain.BetStatus `json:"status"`
	SettledAt *time.Time       `json:"settledAt"`
	Win       bool             `json:"win"`
	BetRequest
}

func AdaptBetFromDomain(bet domain.Bet) BetResponse {
	return BetResponse{
		ID:        bet.ID,
		PlacedAt:  bet.PlacedAt,
		Status:    bet.Status,
		SettledAt: bet.SettledAt,
		Win:       bet.Win,
		BetRequest: BetRequest{
			Stake:          bet.Stake,
			Table:          bet.Table,
			SelectedSpaces: bet.SelectedSpaces,
		},
	}
}

func AdaptBetToDomain(bet BetRequest, tableID uuid.UUID) domain.Bet {
	return domain.Bet{
		ID:             uuid.New(),
		Status:         domain.Unsettled,
		Stake:          bet.Stake,
		SelectedSpaces: bet.SelectedSpaces,
		PlacedAt:       time.Now().UTC(),
		SettledAt:      nil,
		Table:          tableID,
	}
}

func AdaptBetsFromDomain(bets []domain.Bet) []BetResponse {
	responseBets := make([]BetResponse, len(bets))

	for i := range bets {
		responseBets[i] = AdaptBetFromDomain(bets[i])
	}

	return responseBets
}
