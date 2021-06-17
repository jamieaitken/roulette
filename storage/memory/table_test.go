package memory

import (
	"betting/storage"
	"context"
	"sync"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
)

func TestTableStorage_Close_Success(t *testing.T) {
	tests := []struct {
		name           string
		givenID        uuid.UUID
		givenTables    map[uuid.UUID]storage.Table
		expectedTables map[uuid.UUID]storage.Table
	}{
		{
			name:    "given a valid ID, expect it to close",
			givenID: uuid.MustParse("86510953-65f4-4b28-a8ec-398a605e5210"),
			givenTables: map[uuid.UUID]storage.Table{
				uuid.MustParse("86510953-65f4-4b28-a8ec-398a605e5210"): {
					ID:       uuid.MustParse("86510953-65f4-4b28-a8ec-398a605e5210"),
					IsClosed: false,
				},
			},
			expectedTables: map[uuid.UUID]storage.Table{
				uuid.MustParse("86510953-65f4-4b28-a8ec-398a605e5210"): {
					ID:       uuid.MustParse("86510953-65f4-4b28-a8ec-398a605e5210"),
					IsClosed: true,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store := TableStorage{
				tables:  test.givenTables,
				RWMutex: sync.RWMutex{},
			}

			err := store.Close(context.Background(), test.givenID)
			if err != nil {
				t.Fatal(err)
			}

			if !cmp.Equal(store.tables, test.expectedTables) {
				t.Fatal(cmp.Diff(store.tables, test.expectedTables))
			}
		})
	}
}

func TestTableStorage_Close_Fail(t *testing.T) {
	tests := []struct {
		name          string
		givenID       uuid.UUID
		givenTables   map[uuid.UUID]storage.Table
		expectedError error
	}{
		{
			name:    "given an invalid ID, expect it to error",
			givenID: uuid.MustParse("99510953-65f4-4b28-a8ec-398a605e5210"),
			givenTables: map[uuid.UUID]storage.Table{
				uuid.MustParse("86510953-65f4-4b28-a8ec-398a605e5210"): {
					ID:       uuid.MustParse("86510953-65f4-4b28-a8ec-398a605e5210"),
					IsClosed: false,
				},
			},
			expectedError: ErrNoTables,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store := TableStorage{
				tables:  test.givenTables,
				RWMutex: sync.RWMutex{},
			}

			err := store.Close(context.Background(), test.givenID)
			if err == nil {
				t.Fatalf("expected %v, got nil", test.expectedError)
			}

			if !cmp.Equal(err, test.expectedError, cmpopts.EquateErrors()) {
				t.Fatal(cmp.Diff(err, test.expectedError, cmpopts.EquateErrors()))
			}
		})
	}
}

func TestTableStorage_Get_Success(t *testing.T) {
	tests := []struct {
		name          string
		givenID       uuid.UUID
		givenTables   map[uuid.UUID]storage.Table
		expectedTable storage.Table
	}{
		{
			name:    "given a valid ID, expect a Table to be returned",
			givenID: uuid.MustParse("86510953-65f4-4b28-a8ec-398a605e5210"),
			givenTables: map[uuid.UUID]storage.Table{
				uuid.MustParse("86510953-65f4-4b28-a8ec-398a605e5210"): {
					ID:       uuid.MustParse("86510953-65f4-4b28-a8ec-398a605e5210"),
					IsClosed: false,
				},
			},
			expectedTable: storage.Table{
				ID:       uuid.MustParse("86510953-65f4-4b28-a8ec-398a605e5210"),
				IsClosed: false,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store := TableStorage{
				tables:  test.givenTables,
				RWMutex: sync.RWMutex{},
			}

			actual, err := store.Get(context.Background(), test.givenID)
			if err != nil {
				t.Fatal(err)
			}

			if !cmp.Equal(actual, test.expectedTable) {
				t.Fatal(cmp.Diff(actual, test.expectedTable))
			}
		})
	}
}

func TestTableStorage_Get_Fail(t *testing.T) {
	tests := []struct {
		name          string
		givenID       uuid.UUID
		givenTables   map[uuid.UUID]storage.Table
		expectedError error
	}{
		{
			name:    "given an invalid ID, expect ErrNoTables returned",
			givenID: uuid.MustParse("11510953-65f4-4b28-a8ec-398a605e5210"),
			givenTables: map[uuid.UUID]storage.Table{
				uuid.MustParse("86510953-65f4-4b28-a8ec-398a605e5210"): {
					ID:       uuid.MustParse("86510953-65f4-4b28-a8ec-398a605e5210"),
					IsClosed: false,
				},
			},
			expectedError: ErrNoTables,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store := TableStorage{
				tables:  test.givenTables,
				RWMutex: sync.RWMutex{},
			}

			_, err := store.Get(context.Background(), test.givenID)
			if err == nil {
				t.Fatalf("expected %v, got nil", err)
			}

			if !cmp.Equal(err, test.expectedError, cmpopts.EquateErrors()) {
				t.Fatal(cmp.Diff(err, test.expectedError, cmpopts.EquateErrors()))
			}
		})
	}
}

func TestTableStorage_Insert_Success(t *testing.T) {
	tests := []struct {
		name           string
		givenTable     storage.Table
		givenTables    map[uuid.UUID]storage.Table
		expectedTables map[uuid.UUID]storage.Table
	}{
		{
			name: "given a valid ID, expect a Table to be returned",
			givenTable: storage.Table{
				ID:       uuid.MustParse("1e722843-ff0d-4598-8a61-255120b1f3af"),
				IsClosed: false,
				Outcome:  nil,
			},
			givenTables: map[uuid.UUID]storage.Table{
				uuid.MustParse("86510953-65f4-4b28-a8ec-398a605e5210"): {
					ID:       uuid.MustParse("86510953-65f4-4b28-a8ec-398a605e5210"),
					IsClosed: false,
				},
			},
			expectedTables: map[uuid.UUID]storage.Table{
				uuid.MustParse("86510953-65f4-4b28-a8ec-398a605e5210"): {
					ID:       uuid.MustParse("86510953-65f4-4b28-a8ec-398a605e5210"),
					IsClosed: false,
				},
				uuid.MustParse("1e722843-ff0d-4598-8a61-255120b1f3af"): {
					ID:       uuid.MustParse("1e722843-ff0d-4598-8a61-255120b1f3af"),
					IsClosed: false,
					Outcome:  nil,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store := TableStorage{
				tables:  test.givenTables,
				RWMutex: sync.RWMutex{},
			}

			err := store.Insert(context.Background(), test.givenTable)
			if err != nil {
				t.Fatal(err)
			}

			if !cmp.Equal(store.tables, test.expectedTables) {
				t.Fatal(cmp.Diff(store.tables, test.expectedTables))
			}
		})
	}
}

func TestTableStorage_Insert_Fail(t *testing.T) {
	tests := []struct {
		name          string
		givenTable    storage.Table
		givenTables   map[uuid.UUID]storage.Table
		expectedError error
	}{
		{
			name: "given an invalid ID, expect ErrDuplicateTable to be returned",
			givenTable: storage.Table{
				ID:       uuid.MustParse("86510953-65f4-4b28-a8ec-398a605e5210"),
				IsClosed: false,
				Outcome:  nil,
			},
			givenTables: map[uuid.UUID]storage.Table{
				uuid.MustParse("86510953-65f4-4b28-a8ec-398a605e5210"): {
					ID:       uuid.MustParse("86510953-65f4-4b28-a8ec-398a605e5210"),
					IsClosed: false,
				},
			},
			expectedError: ErrDuplicateTable,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store := TableStorage{
				tables:  test.givenTables,
				RWMutex: sync.RWMutex{},
			}

			err := store.Insert(context.Background(), test.givenTable)
			if err == nil {
				t.Fatalf("expected %v, got nil", test.expectedError)
			}

			if !cmp.Equal(err, test.expectedError, cmpopts.EquateErrors()) {
				t.Fatal(cmp.Diff(err, test.expectedError, cmpopts.EquateErrors()))
			}
		})
	}
}

func TestTableStorage_List(t *testing.T) {
	tests := []struct {
		name           string
		givenTables    map[uuid.UUID]storage.Table
		expectedTables []storage.Table
	}{
		{
			name: "expect slice of Tables to be returned",
			givenTables: map[uuid.UUID]storage.Table{
				uuid.MustParse("86510953-65f4-4b28-a8ec-398a605e5210"): {
					ID:       uuid.MustParse("86510953-65f4-4b28-a8ec-398a605e5210"),
					IsClosed: false,
				},
			},
			expectedTables: []storage.Table{
				{
					ID:       uuid.MustParse("86510953-65f4-4b28-a8ec-398a605e5210"),
					IsClosed: false,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store := TableStorage{
				tables:  test.givenTables,
				RWMutex: sync.RWMutex{},
			}

			actual, err := store.List(context.Background())
			if err != nil {
				t.Fatal(err)
			}

			if !cmp.Equal(actual, test.expectedTables) {
				t.Fatal(cmp.Diff(actual, test.expectedTables))
			}
		})
	}
}

func TestTableStorage_SetOutcome_Success(t *testing.T) {
	tests := []struct {
		name           string
		givenID        uuid.UUID
		givenOutcome   storage.Outcome
		givenTables    map[uuid.UUID]storage.Table
		expectedTables map[uuid.UUID]storage.Table
	}{
		{
			name:    "given a valid ID, expect it the outcome to be set",
			givenID: uuid.MustParse("86510953-65f4-4b28-a8ec-398a605e5210"),
			givenOutcome: storage.Outcome{
				Colour: "red",
				Value:  16,
			},
			givenTables: map[uuid.UUID]storage.Table{
				uuid.MustParse("86510953-65f4-4b28-a8ec-398a605e5210"): {
					ID:       uuid.MustParse("86510953-65f4-4b28-a8ec-398a605e5210"),
					IsClosed: false,
				},
			},
			expectedTables: map[uuid.UUID]storage.Table{
				uuid.MustParse("86510953-65f4-4b28-a8ec-398a605e5210"): {
					ID:       uuid.MustParse("86510953-65f4-4b28-a8ec-398a605e5210"),
					IsClosed: false,
					Outcome: &storage.Outcome{
						Colour: "red",
						Value:  16,
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store := TableStorage{
				tables:  test.givenTables,
				RWMutex: sync.RWMutex{},
			}

			err := store.SetOutcome(context.Background(), test.givenID, test.givenOutcome)
			if err != nil {
				t.Fatal(err)
			}

			if !cmp.Equal(store.tables, test.expectedTables) {
				t.Fatal(cmp.Diff(store.tables, test.expectedTables))
			}
		})
	}
}

func TestTableStorage_SetOutcome_Fail(t *testing.T) {
	tests := []struct {
		name          string
		givenID       uuid.UUID
		givenOutcome  storage.Outcome
		givenTables   map[uuid.UUID]storage.Table
		expectedError error
	}{
		{
			name:    "given an invalid ID, expect it to error",
			givenID: uuid.MustParse("99510953-65f4-4b28-a8ec-398a605e5210"),
			givenOutcome: storage.Outcome{
				Colour: "red",
				Value:  16,
			},
			givenTables: map[uuid.UUID]storage.Table{
				uuid.MustParse("86510953-65f4-4b28-a8ec-398a605e5210"): {
					ID:       uuid.MustParse("86510953-65f4-4b28-a8ec-398a605e5210"),
					IsClosed: false,
				},
			},
			expectedError: ErrNoTables,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store := TableStorage{
				tables:  test.givenTables,
				RWMutex: sync.RWMutex{},
			}

			err := store.SetOutcome(context.Background(), test.givenID, test.givenOutcome)
			if err == nil {
				t.Fatalf("expected %v, got nil", test.expectedError)
			}

			if !cmp.Equal(err, test.expectedError, cmpopts.EquateErrors()) {
				t.Fatal(cmp.Diff(err, test.expectedError, cmpopts.EquateErrors()))
			}
		})
	}
}
