package memory

import (
	"betting/internal/domain"
	"betting/storage"
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Errors returned by Storage.
var (
	ErrDuplicateKey = errors.New("duplicate key given")
	ErrInvalidKey   = errors.New("invalid key given")
)

// BetStorage holds the record of all created Bets.
type BetStorage struct {
	bets map[uuid.UUID]storage.Bet
	sync.RWMutex
}

// NewBetStorage instantiates BetStorage.
func NewBetStorage() *BetStorage {
	return &BetStorage{
		bets: make(map[uuid.UUID]storage.Bet),
	}
}

// Get returns a Bet for a given ID.
func (b *BetStorage) Get(_ context.Context, id uuid.UUID) (storage.Bet, error) {
	b.RLock()
	defer b.RUnlock()

	bet, ok := b.bets[id]
	if !ok {
		return storage.Bet{}, ErrInvalidKey
	}

	return bet, nil
}

// Insert creates a new Bet in memory.
func (b *BetStorage) Insert(_ context.Context, bet storage.Bet) error {
	b.Lock()
	defer b.Unlock()
	_, ok := b.bets[bet.ID]
	if ok {
		return ErrDuplicateKey
	}

	b.bets[bet.ID] = bet

	return nil
}

// List returns all the Bets for a given Table ID.
func (b *BetStorage) List(_ context.Context, id uuid.UUID) ([]storage.Bet, error) {
	b.RLock()
	defer b.RUnlock()

	var bets []storage.Bet

	for i := range b.bets {
		if b.bets[i].Table != id {
			continue
		}

		bets = append(bets, b.bets[i])
	}

	return bets, nil
}

// UpdateStateByTableID updates all Bets for a given Table with the given status.
func (b *BetStorage) UpdateStateByTableID(_ context.Context, id uuid.UUID, status domain.BetStatus) error {
	b.Lock()
	defer b.Unlock()

	for i := range b.bets {
		if b.bets[i].Table != id {
			continue
		}

		bet := b.bets[i]

		bet.Status = status.String()

		if status == domain.Settled {
			t := time.Now().UTC()

			bet.SettledAt = &t
		}

		b.bets[i] = bet
	}

	return nil
}

// SetWinners sets Bets to won status if they have the same number as the Outcome.
func (b *BetStorage) SetWinners(_ context.Context, bets []storage.Bet) error {
	b.Lock()
	defer b.Unlock()

	for i := range bets {
		bet := bets[i]

		b.bets[bet.ID] = bet
	}

	return nil
}
