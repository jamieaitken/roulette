package winnerlocator

import (
	"betting/internal/domain"
	"context"
)

type Locator struct {
}

func New() Locator {
	return Locator{}
}

func (l Locator) Locate(_ context.Context, table domain.Table) domain.Table {
	for i := range table.Bets {
		bet := table.Bets[i]

		if !hasBetWon(bet.SelectedSpaces, table.Outcome) {
			continue
		}

		bet.Win = true

		table.Bets[i] = bet
	}

	return table
}

func hasBetWon(bets []int, outcome *domain.Outcome) bool {
	for i := range bets {
		if bets[i] != outcome.Value {
			continue
		}

		return true
	}

	return false
}
