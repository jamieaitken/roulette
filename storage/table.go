package storage

import (
	"betting/internal/domain"

	"github.com/google/uuid"
)

// Table is the storage representation of domain.Table.
type Table struct {
	ID       uuid.UUID
	IsClosed bool
	Outcome  *Outcome
}

// Outcome is the storage representation of domain.Outcome.
type Outcome struct {
	Colour string
	Value  int
}

// AdaptTableToDomain returns a domain.Table for a given storage.Table.
func AdaptTableToDomain(table Table) domain.Table {
	return domain.Table{
		ID:       table.ID,
		IsClosed: table.IsClosed,
		Outcome:  AdaptOutcomeToDomain(table.Outcome),
	}
}

// AdaptOutcomeToDomain adapts a storage Outcome to a domain.Outcome.
func AdaptOutcomeToDomain(outcome *Outcome) *domain.Outcome {
	if outcome == nil {
		return nil
	}

	return &domain.Outcome{
		Value:  outcome.Value,
		Colour: domain.Colour(outcome.Colour),
	}
}

// AdaptOutcomeFromDomain adapts a domain.Outcome to a Outcome.
func AdaptOutcomeFromDomain(outcome *domain.Outcome) *Outcome {
	if outcome == nil {
		return nil
	}

	return &Outcome{
		Value:  outcome.Value,
		Colour: outcome.Colour.String(),
	}
}

// AdaptTableFromDomain adapts a domain.Table to a Table.
func AdaptTableFromDomain(table domain.Table) Table {
	return Table{
		ID:       table.ID,
		IsClosed: table.IsClosed,
		Outcome:  AdaptOutcomeFromDomain(table.Outcome),
	}
}
