package bet

import (
	"betting/internal/domain"
	"context"

	"github.com/google/uuid"
)

type RepositoryProvider interface {
	RepoReader
	RepoWriter
}

type RepoWriter interface {
	Insert(ctx context.Context, bet domain.Bet) error
}

type RepoReader interface {
	Get(ctx context.Context, id uuid.UUID) (domain.Bet, error)
}

type TableRepoProvider interface {
	TableRepoReader
}

// TableRepoReader provides read operations for tables.
type TableRepoReader interface {
	Get(ctx context.Context, id uuid.UUID) (domain.Table, error)
}

type Controller struct {
	RepositoryProvider RepositoryProvider
	TableRepoProvider  TableRepoProvider
}

func NewController(provider RepositoryProvider, repoProvider TableRepoProvider) Controller {
	return Controller{
		RepositoryProvider: provider,
		TableRepoProvider:  repoProvider,
	}
}

func (c Controller) Create(ctx context.Context, bet domain.Bet) (domain.Bet, error) {
	table, err := c.TableRepoProvider.Get(ctx, bet.Table)
	if err != nil {
		return domain.Bet{}, err
	}

	if table.IsClosed {
		return domain.Bet{}, ErrTableClosed
	}

	err = c.RepositoryProvider.Insert(ctx, bet)
	if err != nil {
		return domain.Bet{}, err
	}

	return bet, nil
}

func (c Controller) Get(ctx context.Context, id uuid.UUID) (domain.Bet, error) {
	bet, err := c.RepositoryProvider.Get(ctx, id)
	if err != nil {
		return domain.Bet{}, err
	}

	return bet, nil
}
