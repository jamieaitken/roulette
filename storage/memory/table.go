package memory

import (
	"betting/storage"
	"context"
	"errors"
	"sync"

	"github.com/google/uuid"
)

// Errors returned by the Storage.
var (
	ErrNoTables       = errors.New("could not locate table")
	ErrDuplicateTable = errors.New("duplicate key for table")
)

// TableStorage holds the record of all created Tables.
type TableStorage struct {
	tables map[uuid.UUID]storage.Table
	sync.RWMutex
}

// NewTableStorage instantiates TableStorage.
func NewTableStorage() *TableStorage {
	return &TableStorage{
		tables: make(map[uuid.UUID]storage.Table),
	}
}

// Close sets the table to closed so no more bets can be added to it.
func (t *TableStorage) Close(_ context.Context, id uuid.UUID) error {
	t.Lock()
	defer t.Unlock()

	table, ok := t.tables[id]
	if !ok {
		return ErrNoTables
	}

	table.IsClosed = true

	t.tables[id] = table

	return nil
}

// Get returns a Table for a given ID.
func (t *TableStorage) Get(_ context.Context, id uuid.UUID) (storage.Table, error) {
	t.RLock()
	defer t.RUnlock()

	table, ok := t.tables[id]
	if !ok {
		return storage.Table{}, ErrNoTables
	}

	return table, nil
}

// Insert creates a new Table in memory.
func (t *TableStorage) Insert(_ context.Context, table storage.Table) error {
	t.Lock()
	defer t.Unlock()

	_, ok := t.tables[table.ID]
	if ok {
		return ErrDuplicateTable
	}

	t.tables[table.ID] = table

	return nil
}

// List returns all Tables in memory.
func (t *TableStorage) List(_ context.Context) ([]storage.Table, error) {
	t.RLock()
	defer t.RUnlock()

	var list []storage.Table

	for _, table := range t.tables {
		list = append(list, table)
	}

	return list, nil
}

// SetOutcome updates the table with the result.
func (t *TableStorage) SetOutcome(_ context.Context, id uuid.UUID, outcome storage.Outcome) error {
	t.Lock()
	defer t.Unlock()

	table, ok := t.tables[id]
	if !ok {
		return ErrNoTables
	}

	table.Outcome = &outcome

	t.tables[id] = table

	return nil
}
