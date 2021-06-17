package bet

import (
	"betting/internal/domain"
	"betting/storage"
	"context"
	"errors"

	"github.com/google/uuid"
)

// ErrTableClosed returned if a Bet is made against a closed Table.
var (
	ErrTableClosed = errors.New("table is not accepting anymore bets")
)

// StorageProvider provides both read and write operations for Bets.
type StorageProvider interface {
	StorageReader
	StorageWriter
}

// StorageReader provides read operations for Bets.
type StorageReader interface {
	Get(ctx context.Context, id uuid.UUID) (storage.Bet, error)
	List(ctx context.Context, id uuid.UUID) ([]storage.Bet, error)
}

// StorageWriter provides write operations for Bets.
type StorageWriter interface {
	Insert(context.Context, storage.Bet) error
	SetWinners(ctx context.Context, bets []storage.Bet) error
	UpdateStateByTableID(ctx context.Context, id uuid.UUID, status domain.BetStatus) error
}

// Repository allows for Bets to be stored in memory.
type Repository struct {
	StorageProvider StorageProvider
}

// NewRepository instantiates a Repository.
func NewRepository(provider StorageProvider) Repository {
	return Repository{
		StorageProvider: provider,
	}
}

// Insert creates a Bet in memory.
func (r Repository) Insert(ctx context.Context, bet domain.Bet) error {
	adaptedBet := storage.AdaptBetToStorage(bet)

	return r.StorageProvider.Insert(ctx, adaptedBet)
}

// Get retrieves a Bet for a given ID.
func (r Repository) Get(ctx context.Context, id uuid.UUID) (domain.Bet, error) {
	bet, err := r.StorageProvider.Get(ctx, id)
	if err != nil {
		return domain.Bet{}, err
	}

	return storage.AdaptBetToDomain(bet), nil
}

func (r Repository) List(ctx context.Context, id uuid.UUID) ([]domain.Bet, error) {
	bets, err := r.StorageProvider.List(ctx, id)
	if err != nil {
		return nil, err
	}

	return storage.AdaptBetsToDomain(bets), nil
}

// SetWinners adapts from domain to storage and sets the winners of a given Table if any.
func (r Repository) SetWinners(ctx context.Context, bets []domain.Bet) error {
	b := storage.AdaptBetsToStorage(bets)

	return r.StorageProvider.SetWinners(ctx, b)
}

// Spin sets all Bets for a given Table to live.
func (r Repository) Spin(ctx context.Context, id uuid.UUID) error {
	return r.StorageProvider.UpdateStateByTableID(ctx, id, domain.Live)
}

// Settle sets all Bets for a given Table to settled.
func (r Repository) Settle(ctx context.Context, id uuid.UUID) error {
	return r.StorageProvider.UpdateStateByTableID(ctx, id, domain.Settled)
}
