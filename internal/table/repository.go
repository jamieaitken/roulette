package table

import (
	"betting/internal/domain"
	"betting/storage"
	"context"

	"github.com/google/uuid"
)

// StorageProvider provides both read and write operations for Tables.
type StorageProvider interface {
	StorageWriter
	StorageReader
}

// StorageWriter provides write operations for Tables.
type StorageWriter interface {
	Close(ctx context.Context, id uuid.UUID) error
	Insert(ctx context.Context, table storage.Table) error
	SetOutcome(ctx context.Context, id uuid.UUID, outcome storage.Outcome) error
}

// StorageReader provides read operations for Tables.
type StorageReader interface {
	Get(ctx context.Context, id uuid.UUID) (storage.Table, error)
	List(ctx context.Context) ([]storage.Table, error)
}

// Repository allows for storing Tables and Bets in memory.
type Repository struct {
	StorageProvider StorageProvider
}

// NewRepository instantiates Repository.
func NewRepository(storageProvider StorageProvider) Repository {
	return Repository{
		StorageProvider: storageProvider,
	}
}

// Close sets the Table to closed so no more Bets can be added to it.
func (r Repository) Close(ctx context.Context, id uuid.UUID) error {
	return r.StorageProvider.Close(ctx, id)
}

// Insert adapts from domain to storage and stores it in memory.
func (r Repository) Insert(ctx context.Context, table domain.Table) error {
	adaptedTable := storage.AdaptTableFromDomain(table)

	return r.StorageProvider.Insert(ctx, adaptedTable)
}

// Get retrieves a Table for a given ID.
func (r Repository) Get(ctx context.Context, id uuid.UUID) (domain.Table, error) {
	table, err := r.StorageProvider.Get(ctx, id)
	if err != nil {
		return domain.Table{}, err
	}

	return storage.AdaptTableToDomain(table), nil
}

// List retrieves all Tables.
func (r Repository) List(ctx context.Context) ([]domain.Table, error) {
	tables, err := r.StorageProvider.List(ctx)
	if err != nil {
		return nil, err
	}

	ts := make([]domain.Table, len(tables))

	for i := range tables {
		ts[i] = storage.AdaptTableToDomain(tables[i])
	}

	return ts, nil
}

// SetOutcome adapts from domain to storage and updates the given Table in memory.
func (r Repository) SetOutcome(ctx context.Context, id uuid.UUID, outcome domain.Outcome) error {
	storageOutcome := storage.AdaptOutcomeFromDomain(&outcome)

	return r.StorageProvider.SetOutcome(ctx, id, *storageOutcome)
}
