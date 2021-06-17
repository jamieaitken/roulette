package table

import (
	"betting/internal/domain"
	"betting/storage"
	"betting/storage/memory"
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/google/uuid"
)

func TestRepository_Close_Success(t *testing.T) {
	tests := []struct {
		name              string
		givenID           uuid.UUID
		givenTableStorage StorageProvider
	}{
		{
			name:              "given an id, expect the table to be closed",
			givenID:           uuid.MustParse("c4b39dc0-2ff4-4405-b3cb-c4f87a9c82fb"),
			givenTableStorage: mockTableStorage{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := NewRepository(test.givenTableStorage)

			err := repo.Close(context.Background(), test.givenID)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestRepository_Close_Fail(t *testing.T) {
	tests := []struct {
		name              string
		givenID           uuid.UUID
		givenTableStorage StorageProvider
		expectedError     error
	}{
		{
			name:    "given a storage error, expect it to be returned",
			givenID: uuid.MustParse("c4b39dc0-2ff4-4405-b3cb-c4f87a9c82fb"),
			givenTableStorage: mockTableStorage{
				GivenCloseError: memory.ErrNoTables,
			},
			expectedError: memory.ErrNoTables,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := NewRepository(test.givenTableStorage)

			err := repo.Close(context.Background(), test.givenID)
			if err == nil {
				t.Fatalf("expected %v, got nil", test.expectedError)
			}

			if !cmp.Equal(err, test.expectedError, cmpopts.EquateErrors()) {
				t.Fatal(cmp.Diff(err, test.expectedError, cmpopts.EquateErrors()))
			}
		})
	}
}

func TestRepository_Insert_Success(t *testing.T) {
	tests := []struct {
		name              string
		givenTable        domain.Table
		givenTableStorage StorageProvider
	}{
		{
			name: "given a table, expect the table to be created",
			givenTable: domain.Table{
				ID:       uuid.MustParse("c4b39dc0-2ff4-4405-b3cb-c4f87a9c82fb"),
				Bets:     nil,
				IsClosed: false,
				Outcome:  nil,
			},
			givenTableStorage: mockTableStorage{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := NewRepository(test.givenTableStorage)

			err := repo.Insert(context.Background(), test.givenTable)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestRepository_Get_Success(t *testing.T) {
	tests := []struct {
		name              string
		givenID           uuid.UUID
		givenTableStorage StorageProvider
		expectedTable     domain.Table
	}{
		{
			name:    "given an id, expect the table to be returned",
			givenID: uuid.MustParse("c4b39dc0-2ff4-4405-b3cb-c4f87a9c82fb"),
			givenTableStorage: mockTableStorage{
				GivenGetTable: storage.Table{
					ID:       uuid.MustParse("c4b39dc0-2ff4-4405-b3cb-c4f87a9c82fb"),
					IsClosed: false,
					Outcome:  nil,
				},
			},
			expectedTable: domain.Table{
				ID:       uuid.MustParse("c4b39dc0-2ff4-4405-b3cb-c4f87a9c82fb"),
				Bets:     nil,
				IsClosed: false,
				Outcome:  nil,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := NewRepository(test.givenTableStorage)

			actual, err := repo.Get(context.Background(), test.givenID)
			if err != nil {
				t.Fatal(err)
			}

			if !cmp.Equal(actual, test.expectedTable) {
				t.Fatal(cmp.Diff(actual, test.expectedTable))
			}
		})
	}
}

func TestRepository_Get_Fail(t *testing.T) {
	tests := []struct {
		name              string
		givenID           uuid.UUID
		givenTableStorage StorageProvider
		expectedError     error
	}{
		{
			name:    "given a table storage error, expect it to be returned",
			givenID: uuid.MustParse("c4b39dc0-2ff4-4405-b3cb-c4f87a9c82fb"),
			givenTableStorage: mockTableStorage{
				GivenGetError: memory.ErrNoTables,
			},
			expectedError: memory.ErrNoTables,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := NewRepository(test.givenTableStorage)

			_, err := repo.Get(context.Background(), test.givenID)
			if err == nil {
				t.Fatalf("expected %v, got nil", test.expectedError)
			}

			if !cmp.Equal(err, test.expectedError, cmpopts.EquateErrors()) {
				t.Fatal(cmp.Diff(err, test.expectedError, cmpopts.EquateErrors()))
			}
		})
	}
}

func TestRepository_List_Success(t *testing.T) {
	tests := []struct {
		name              string
		givenTableStorage StorageProvider
		expectedTables    []domain.Table
	}{
		{
			name: "expect all tables to be returned",
			givenTableStorage: mockTableStorage{
				GivenListTables: []storage.Table{
					{
						ID:       uuid.MustParse("c4b39dc0-2ff4-4405-b3cb-c4f87a9c82fb"),
						IsClosed: false,
						Outcome:  nil,
					},
				},
			},
			expectedTables: []domain.Table{
				{
					ID:       uuid.MustParse("c4b39dc0-2ff4-4405-b3cb-c4f87a9c82fb"),
					Bets:     nil,
					IsClosed: false,
					Outcome:  nil,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := NewRepository(test.givenTableStorage)

			actual, err := repo.List(context.Background())
			if err != nil {
				t.Fatal(err)
			}

			if !cmp.Equal(actual, test.expectedTables) {
				t.Fatal(cmp.Diff(actual, test.expectedTables))
			}
		})
	}
}

func TestRepository_List_Fail(t *testing.T) {
	tests := []struct {
		name              string
		givenTableStorage StorageProvider
		expectedError     error
	}{
		{
			name: "given fail to find table, expect error to be returned",
			givenTableStorage: mockTableStorage{
				GivenListError: memory.ErrNoTables,
			},
			expectedError: memory.ErrNoTables,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := NewRepository(test.givenTableStorage)

			_, err := repo.List(context.Background())
			if err == nil {
				t.Fatalf("expected %v, got nil", test.expectedError)
			}

			if !cmp.Equal(err, test.expectedError, cmpopts.EquateErrors()) {
				t.Fatal(cmp.Diff(err, test.expectedError, cmpopts.EquateErrors()))
			}
		})
	}
}

func TestRepository_SetOutcome_Success(t *testing.T) {
	tests := []struct {
		name              string
		givenID           uuid.UUID
		givenOutcome      domain.Outcome
		givenTableStorage StorageProvider
	}{
		{
			name:    "given an id and outcome, expect the table to updated",
			givenID: uuid.MustParse("c4b39dc0-2ff4-4405-b3cb-c4f87a9c82fb"),
			givenOutcome: domain.Outcome{
				Value:  16,
				Colour: "red",
			},
			givenTableStorage: mockTableStorage{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := NewRepository(test.givenTableStorage)

			err := repo.SetOutcome(context.Background(), test.givenID, test.givenOutcome)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

type mockTableStorage struct {
	GivenCloseError      error
	GivenInsertError     error
	GivenSetOutcomeError error
	GivenGetTable        storage.Table
	GivenGetError        error
	GivenListTables      []storage.Table
	GivenListError       error
}

func (m mockTableStorage) Close(_ context.Context, _ uuid.UUID) error {
	return m.GivenCloseError
}

func (m mockTableStorage) Insert(_ context.Context, _ storage.Table) error {
	return m.GivenInsertError
}

func (m mockTableStorage) SetOutcome(_ context.Context, _ uuid.UUID, _ storage.Outcome) error {
	return m.GivenSetOutcomeError
}

func (m mockTableStorage) Get(_ context.Context, _ uuid.UUID) (storage.Table, error) {
	return m.GivenGetTable, m.GivenGetError
}

func (m mockTableStorage) List(_ context.Context) ([]storage.Table, error) {
	return m.GivenListTables, m.GivenListError
}
