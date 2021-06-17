package bet

import (
	"betting/internal/domain"
	"betting/storage"
	"betting/storage/memory"
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/google/go-cmp/cmp"

	"github.com/google/uuid"
)

func TestRepository_Insert_Success(t *testing.T) {
	tests := []struct {
		name            string
		givenBet        domain.Bet
		givenBetStorage *mockBetStorage
		expectedBet     storage.Bet
	}{
		{
			name: "given a bet, expect it to be stored",
			givenBet: domain.Bet{
				ID:             uuid.MustParse("c4b39dc0-2ff4-4405-b3cb-c4f87a9c82fb"),
				Status:         "",
				SelectedSpaces: nil,
				Stake:          nil,
				PlacedAt:       time.Time{},
				SettledAt:      nil,
				Win:            false,
				Table:          uuid.UUID{},
			},
			givenBetStorage: &mockBetStorage{},
			expectedBet: storage.Bet{
				ID:             uuid.MustParse("c4b39dc0-2ff4-4405-b3cb-c4f87a9c82fb"),
				Status:         "",
				SelectedSpaces: nil,
				Stake:          nil,
				PlacedAt:       time.Time{},
				SettledAt:      nil,
				Win:            false,
				Table:          uuid.UUID{},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := NewRepository(test.givenBetStorage)

			err := repo.Insert(context.Background(), test.givenBet)
			if err != nil {
				t.Fatal(err)
			}

			if !cmp.Equal(test.givenBetStorage.SpyInsertBet, test.expectedBet) {
				t.Fatal(cmp.Diff(test.givenBetStorage.SpyInsertBet, test.expectedBet))
			}
		})
	}
}

func TestRepository_Get_Success(t *testing.T) {
	tests := []struct {
		name            string
		givenID         uuid.UUID
		givenBetStorage *mockBetStorage
		expectedBet     domain.Bet
	}{
		{
			name:    "given an id, expect the bet to be returned",
			givenID: uuid.MustParse("c4b39dc0-2ff4-4405-b3cb-c4f87a9c82fb"),
			givenBetStorage: &mockBetStorage{
				GivenGetBet: storage.Bet{
					ID:             uuid.MustParse("c4b39dc0-2ff4-4405-b3cb-c4f87a9c82fb"),
					Status:         "",
					SelectedSpaces: nil,
					Stake:          nil,
					PlacedAt:       time.Time{},
					SettledAt:      nil,
					Win:            false,
					Table:          uuid.UUID{},
				},
			},
			expectedBet: domain.Bet{
				ID:             uuid.MustParse("c4b39dc0-2ff4-4405-b3cb-c4f87a9c82fb"),
				Status:         "",
				SelectedSpaces: nil,
				Stake:          nil,
				PlacedAt:       time.Time{},
				SettledAt:      nil,
				Win:            false,
				Table:          uuid.UUID{},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := NewRepository(test.givenBetStorage)

			actual, err := repo.Get(context.Background(), test.givenID)
			if err != nil {
				t.Fatal(err)
			}

			if !cmp.Equal(actual, test.expectedBet) {
				t.Fatal(cmp.Diff(actual, test.expectedBet))
			}
		})
	}
}

func TestRepository_Get_Fail(t *testing.T) {
	tests := []struct {
		name            string
		givenID         uuid.UUID
		givenBetStorage *mockBetStorage
		expectedError   error
	}{
		{
			name:    "given an invalid id, expect an error to be returned",
			givenID: uuid.MustParse("c4b39dc0-2ff4-4405-b3cb-c4f87a9c82fb"),
			givenBetStorage: &mockBetStorage{
				GivenGetError: memory.ErrInvalidKey,
			},
			expectedError: memory.ErrInvalidKey,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := NewRepository(test.givenBetStorage)

			_, err := repo.Get(context.Background(), test.givenID)
			if err == nil {
				t.Fatalf("exected %v, got nil", test.expectedError)
			}

			if !cmp.Equal(err, test.expectedError, cmpopts.EquateErrors()) {
				t.Fatal(cmp.Diff(err, test.expectedError, cmpopts.EquateErrors()))
			}
		})
	}
}

func TestRepository_List_Success(t *testing.T) {
	tests := []struct {
		name            string
		givenTableID    uuid.UUID
		givenBetStorage *mockBetStorage
		expectedBets    []domain.Bet
	}{
		{
			name:         "given an id, expect the bet to be returned",
			givenTableID: uuid.MustParse("c4b39dc0-2ff4-4405-b3cb-c4f87a9c82fb"),
			givenBetStorage: &mockBetStorage{
				GivenListBets: []storage.Bet{
					{
						ID:             uuid.MustParse("68c094d3-ae37-4e4d-a533-8701dc1d7e5c"),
						Status:         "",
						SelectedSpaces: nil,
						Stake:          nil,
						PlacedAt:       time.Time{},
						SettledAt:      nil,
						Win:            false,
						Table:          uuid.MustParse("c4b39dc0-2ff4-4405-b3cb-c4f87a9c82fb"),
					},
				},
			},
			expectedBets: []domain.Bet{
				{
					ID:             uuid.MustParse("68c094d3-ae37-4e4d-a533-8701dc1d7e5c"),
					Status:         "",
					SelectedSpaces: nil,
					Stake:          nil,
					PlacedAt:       time.Time{},
					SettledAt:      nil,
					Win:            false,
					Table:          uuid.MustParse("c4b39dc0-2ff4-4405-b3cb-c4f87a9c82fb"),
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := NewRepository(test.givenBetStorage)

			actual, err := repo.List(context.Background(), test.givenTableID)
			if err != nil {
				t.Fatal(err)
			}

			if !cmp.Equal(actual, test.expectedBets) {
				t.Fatal(cmp.Diff(actual, test.expectedBets))
			}
		})
	}
}

func TestRepository_List_Fail(t *testing.T) {
	tests := []struct {
		name            string
		givenID         uuid.UUID
		givenBetStorage *mockBetStorage
		expectedError   error
	}{
		{
			name:    "given an invalid id, expect an error to be returned",
			givenID: uuid.MustParse("c4b39dc0-2ff4-4405-b3cb-c4f87a9c82fb"),
			givenBetStorage: &mockBetStorage{
				GivenListError: memory.ErrInvalidKey,
			},
			expectedError: memory.ErrInvalidKey,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := NewRepository(test.givenBetStorage)

			_, err := repo.List(context.Background(), test.givenID)
			if err == nil {
				t.Fatalf("exected %v, got nil", test.expectedError)
			}

			if !cmp.Equal(err, test.expectedError, cmpopts.EquateErrors()) {
				t.Fatal(cmp.Diff(err, test.expectedError, cmpopts.EquateErrors()))
			}
		})
	}
}

func TestRepository_Spin_Success(t *testing.T) {
	tests := []struct {
		name            string
		givenID         uuid.UUID
		givenBetStorage StorageProvider
	}{
		{
			name:            "given an id, expect the table's bets to be live",
			givenID:         uuid.MustParse("c4b39dc0-2ff4-4405-b3cb-c4f87a9c82fb"),
			givenBetStorage: &mockBetStorage{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := NewRepository(test.givenBetStorage)

			err := repo.Spin(context.Background(), test.givenID)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestRepository_Settle_Success(t *testing.T) {
	tests := []struct {
		name            string
		givenID         uuid.UUID
		givenBetStorage StorageProvider
	}{
		{
			name:            "given an id, expect the bets to be settled",
			givenID:         uuid.MustParse("c4b39dc0-2ff4-4405-b3cb-c4f87a9c82fb"),
			givenBetStorage: &mockBetStorage{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := NewRepository(test.givenBetStorage)

			err := repo.Settle(context.Background(), test.givenID)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestRepository_SetWinners_Success(t *testing.T) {
	tests := []struct {
		name            string
		givenBets       []domain.Bet
		givenBetStorage StorageProvider
	}{
		{
			name:            "given a table, expect any winners to be found",
			givenBets:       []domain.Bet{},
			givenBetStorage: &mockBetStorage{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := NewRepository(test.givenBetStorage)

			err := repo.SetWinners(context.Background(), test.givenBets)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

type mockBetStorage struct {
	GivenGetBet           storage.Bet
	GivenGetError         error
	GivenListBets         []storage.Bet
	GivenListError        error
	GivenInsertError      error
	GivenSetWinnersError  error
	GivenUpdateStateError error
	SpyInsertBet          storage.Bet
}

func (m *mockBetStorage) Get(_ context.Context, _ uuid.UUID) (storage.Bet, error) {
	return m.GivenGetBet, m.GivenGetError
}

func (m *mockBetStorage) List(_ context.Context, _ uuid.UUID) ([]storage.Bet, error) {
	return m.GivenListBets, m.GivenListError
}

func (m *mockBetStorage) Insert(_ context.Context, bet storage.Bet) error {
	m.SpyInsertBet = bet

	return m.GivenInsertError
}

func (m *mockBetStorage) SetWinners(_ context.Context, _ []storage.Bet) error {
	return m.GivenSetWinnersError
}

func (m *mockBetStorage) UpdateStateByTableID(_ context.Context, _ uuid.UUID, _ domain.BetStatus) error {
	return m.GivenUpdateStateError
}
