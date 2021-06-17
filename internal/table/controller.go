package table

import (
	"betting/internal/domain"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

// Domain errors.
var (
	ErrFailedToCreateTable  = errors.New("failed to create table")
	ErrFailedToCloseTable   = errors.New("failed to close table")
	ErrFailedToSpinTable    = errors.New("failed to spin table")
	ErrFailedToSetOutcome   = errors.New("failed to set outcome on table")
	ErrFailedToFetchTable   = errors.New("failed to locate table")
	ErrFailedToFetchBets    = errors.New("failed to locate bets")
	ErrFailedFailedToSettle = errors.New("failed to settle bets")
	ErrFailedToSetWinners   = errors.New("failed to set winners")
)

// RepositoryProvider provides both read and write operations for Tables.
type RepositoryProvider interface {
	Reader
	Writer
}

// Reader provides read operations for Tables.
type Reader interface {
	Get(ctx context.Context, id uuid.UUID) (domain.Table, error)
	List(ctx context.Context) ([]domain.Table, error)
}

// Writer provides write operations for Tables.
type Writer interface {
	Insert(ctx context.Context, table domain.Table) error
	Close(ctx context.Context, id uuid.UUID) error
	SetOutcome(ctx context.Context, id uuid.UUID, outcome domain.Outcome) error
}

// BallPlacer generates the landing position of the ball.
type BallPlacer interface {
	GetPosition(ctx context.Context) domain.Outcome
}

// WinnerLocator finds all winning bets.
type WinnerLocator interface {
	Locate(ctx context.Context, table domain.Table) domain.Table
}

// BetRepositoryProvider provides read and write operations for Bet storage.
type BetRepositoryProvider interface {
	BetRepositoryWriter
	BetRepositoryReader
}

// BetRepositoryWriter provides write operations for Bet storage.
type BetRepositoryWriter interface {
	SetWinners(ctx context.Context, bets []domain.Bet) error
	Spin(ctx context.Context, id uuid.UUID) error
	Settle(ctx context.Context, id uuid.UUID) error
}

// BetRepositoryReader provides read operations for Bet storage.
type BetRepositoryReader interface {
	List(ctx context.Context, id uuid.UUID) ([]domain.Bet, error)
}

// Controller is responsible for doing business logic for a Table.
type Controller struct {
	RepositoryProvider    RepositoryProvider
	BallPlacer            BallPlacer
	WinnerLocator         WinnerLocator
	BetRepositoryProvider BetRepositoryProvider
}

// ControllerParams hold the dependencies required for a Controller.
type ControllerParams struct {
	RepositoryProvider    RepositoryProvider
	BallPlacer            BallPlacer
	WinnerLocator         WinnerLocator
	BetRepositoryProvider BetRepositoryProvider
}

// NewController instantiates Controller.
func NewController(p ControllerParams) Controller {
	return Controller{
		RepositoryProvider:    p.RepositoryProvider,
		BallPlacer:            p.BallPlacer,
		WinnerLocator:         p.WinnerLocator,
		BetRepositoryProvider: p.BetRepositoryProvider,
	}
}

// Create generates a new Table in memory.
func (c Controller) Create(ctx context.Context) (domain.Table, error) {
	table := domain.Table{
		ID:       uuid.New(),
		Bets:     nil,
		IsClosed: false,
	}

	err := c.RepositoryProvider.Insert(ctx, table)
	if err != nil {
		return domain.Table{}, fmt.Errorf("%v: %w", err, ErrFailedToCreateTable)
	}

	return table, nil
}

// Spin closes the Table, sets all Bets to live, generates the outcome, updates the Table with outcome and returns
// updated resource.
func (c Controller) Spin(ctx context.Context, id uuid.UUID) (domain.Table, error) {
	err := c.RepositoryProvider.Close(ctx, id)
	if err != nil {
		return domain.Table{}, fmt.Errorf("%v: %w", err, ErrFailedToCloseTable)
	}

	err = c.BetRepositoryProvider.Spin(ctx, id)
	if err != nil {
		return domain.Table{}, fmt.Errorf("%v: %w", err, ErrFailedToSpinTable)
	}

	position := c.BallPlacer.GetPosition(ctx)

	err = c.RepositoryProvider.SetOutcome(ctx, id, position)
	if err != nil {
		return domain.Table{}, fmt.Errorf("%v: %w", err, ErrFailedToSetOutcome)
	}

	table, err := c.RepositoryProvider.Get(ctx, id)
	if err != nil {
		return domain.Table{}, fmt.Errorf("%v: %w", err, ErrFailedToFetchTable)
	}

	bets, err := c.BetRepositoryProvider.List(ctx, table.ID)
	if err != nil {
		return domain.Table{}, fmt.Errorf("%v: %w", err, ErrFailedToFetchBets)
	}

	table.Bets = bets

	return table, nil
}

// Settle updates all Bets to settled, finds all Winners (if any) and returns the updated Table.
func (c Controller) Settle(ctx context.Context, id uuid.UUID) (domain.Table, error) {
	err := c.BetRepositoryProvider.Settle(ctx, id)
	if err != nil {
		return domain.Table{}, fmt.Errorf("%v: %w", err, ErrFailedFailedToSettle)
	}

	table, err := c.RepositoryProvider.Get(ctx, id)
	if err != nil {
		return domain.Table{}, fmt.Errorf("%v: %w", err, ErrFailedToFetchTable)
	}

	table = c.WinnerLocator.Locate(ctx, table)

	err = c.BetRepositoryProvider.SetWinners(ctx, table.Bets)
	if err != nil {
		return domain.Table{}, fmt.Errorf("%v: %w", err, ErrFailedToSetWinners)
	}

	bets, err := c.BetRepositoryProvider.List(ctx, table.ID)
	if err != nil {
		return domain.Table{}, fmt.Errorf("%v: %w", err, ErrFailedToFetchBets)
	}

	table.Bets = bets

	return table, nil
}

// Get retrieves a Table for a given ID.
func (c Controller) Get(ctx context.Context, id uuid.UUID) (domain.Table, error) {
	table, err := c.RepositoryProvider.Get(ctx, id)
	if err != nil {
		return domain.Table{}, fmt.Errorf("%v: %w", err, ErrFailedToFetchTable)
	}

	bets, err := c.BetRepositoryProvider.List(ctx, table.ID)
	if err != nil {
		return domain.Table{}, fmt.Errorf("%v: %w", err, ErrFailedToFetchBets)
	}

	table.Bets = bets

	return table, nil
}

// List returns all Tables and their associated Bets.
func (c Controller) List(ctx context.Context) ([]domain.Table, error) {
	tables, err := c.RepositoryProvider.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, ErrFailedToFetchTable)
	}

	for i := range tables {
		bets, err := c.BetRepositoryProvider.List(ctx, tables[i].ID)
		if err != nil {
			continue
		}

		tables[i].Bets = bets
	}

	return tables, nil
}
