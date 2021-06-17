package api

import (
	"betting/internal/domain"

	"github.com/google/uuid"
)

// TableResponse is the presentation representation of a domain.Table.
type TableResponse struct {
	ID       uuid.UUID     `json:"id"`
	Bets     []BetResponse `json:"bets"`
	IsClosed bool          `json:"isClosed"`
	Outcome  *Outcome      `json:"outcome"`
}

// Outcome is the result of the Table.
type Outcome struct {
	Position int           `json:"position"`
	Colour   domain.Colour `json:"colour"`
}

func AdaptOutcomeFromDomain(outcome *domain.Outcome) *Outcome {
	if outcome == nil {
		return nil
	}

	return &Outcome{
		Position: outcome.Value,
		Colour:   outcome.Colour,
	}
}

func AdaptTableFromDomain(table domain.Table) TableResponse {
	return TableResponse{
		ID:       table.ID,
		Bets:     AdaptBetsFromDomain(table.Bets),
		IsClosed: table.IsClosed,
		Outcome:  AdaptOutcomeFromDomain(table.Outcome),
	}
}

func AdaptTablesFromDomain(tables []domain.Table) []TableResponse {
	ts := make([]TableResponse, len(tables))

	for i := range tables {
		ts[i] = AdaptTableFromDomain(tables[i])
	}

	return ts
}
