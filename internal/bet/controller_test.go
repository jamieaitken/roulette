package bet

import (
	"betting/internal/domain"
	"betting/storage/memory"
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/google/go-cmp/cmp"

	"github.com/google/uuid"
)

func TestController_Create_Success(t *testing.T) {
	tests := []struct {
		name           string
		givenBet       domain.Bet
		givenTableRepo TableRepoProvider
		givenBetRepo   RepositoryProvider
		expectedBet    domain.Bet
	}{
		{
			name: "given a bet, expect it to be inserted",
			givenBet: domain.Bet{
				ID:             uuid.MustParse("49cffe67-9798-4327-9760-c4b81562f928"),
				Status:         "",
				SelectedSpaces: []int{5},
				Stake:          nil,
				PlacedAt:       time.Time{},
				SettledAt:      nil,
				Win:            false,
				Table:          uuid.MustParse("0173b64f-e07e-4fa0-bcb3-231856390dce"),
			},
			givenTableRepo: mockTableRepo{
				GivenGetTable: domain.Table{
					ID:       uuid.MustParse("0173b64f-e07e-4fa0-bcb3-231856390dce"),
					Bets:     nil,
					IsClosed: false,
					Outcome:  nil,
				},
			},
			givenBetRepo: mockBetRepo{},
			expectedBet: domain.Bet{
				ID:             uuid.MustParse("49cffe67-9798-4327-9760-c4b81562f928"),
				Status:         "",
				SelectedSpaces: []int{5},
				Stake:          nil,
				PlacedAt:       time.Time{},
				SettledAt:      nil,
				Win:            false,
				Table:          uuid.MustParse("0173b64f-e07e-4fa0-bcb3-231856390dce"),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := NewController(test.givenBetRepo, test.givenTableRepo)

			actual, err := c.Create(context.Background(), test.givenBet)
			if err != nil {
				t.Fatal(err)
			}

			if !cmp.Equal(actual, test.expectedBet) {
				t.Fatal(cmp.Diff(actual, test.expectedBet))
			}
		})
	}
}

func TestController_Create_Fail(t *testing.T) {
	tests := []struct {
		name           string
		givenBet       domain.Bet
		givenTableRepo TableRepoProvider
		givenBetRepo   RepositoryProvider
		expectedError  error
	}{
		{
			name: "given a table repo error, expect it to be returned",
			givenBet: domain.Bet{
				ID:             uuid.MustParse("49cffe67-9798-4327-9760-c4b81562f928"),
				Status:         "",
				SelectedSpaces: []int{5},
				Stake:          nil,
				PlacedAt:       time.Time{},
				SettledAt:      nil,
				Win:            false,
				Table:          uuid.MustParse("0173b64f-e07e-4fa0-bcb3-231856390dce"),
			},
			givenTableRepo: mockTableRepo{
				GivenGetError: memory.ErrInvalidKey,
			},
			givenBetRepo:  mockBetRepo{},
			expectedError: memory.ErrInvalidKey,
		},
		{
			name: "given a table that is closed, expect error to be returned",
			givenBet: domain.Bet{
				ID:             uuid.MustParse("49cffe67-9798-4327-9760-c4b81562f928"),
				Status:         "",
				SelectedSpaces: []int{5},
				Stake:          nil,
				PlacedAt:       time.Time{},
				SettledAt:      nil,
				Win:            false,
				Table:          uuid.MustParse("0173b64f-e07e-4fa0-bcb3-231856390dce"),
			},
			givenTableRepo: mockTableRepo{
				GivenGetTable: domain.Table{
					ID:       uuid.MustParse("0173b64f-e07e-4fa0-bcb3-231856390dce"),
					Bets:     nil,
					IsClosed: true,
					Outcome:  nil,
				},
			},
			givenBetRepo:  mockBetRepo{},
			expectedError: ErrTableClosed,
		},
		{
			name: "given an insert repo error, expect it to be returned",
			givenBet: domain.Bet{
				ID:             uuid.MustParse("49cffe67-9798-4327-9760-c4b81562f928"),
				Status:         "",
				SelectedSpaces: []int{5},
				Stake:          nil,
				PlacedAt:       time.Time{},
				SettledAt:      nil,
				Win:            false,
				Table:          uuid.MustParse("0173b64f-e07e-4fa0-bcb3-231856390dce"),
			},
			givenTableRepo: mockTableRepo{},
			givenBetRepo: mockBetRepo{
				GivenInsertError: memory.ErrDuplicateKey,
			},
			expectedError: memory.ErrDuplicateKey,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := NewController(test.givenBetRepo, test.givenTableRepo)

			_, err := c.Create(context.Background(), test.givenBet)
			if err == nil {
				t.Fatalf("expected %v, got nil", test.expectedError)
			}

			if !cmp.Equal(err, test.expectedError, cmpopts.EquateErrors()) {
				t.Fatal(cmp.Diff(err, test.expectedError, cmpopts.EquateErrors()))
			}
		})
	}
}

func TestController_Get_Success(t *testing.T) {
	tests := []struct {
		name           string
		givenID        uuid.UUID
		givenTableRepo TableRepoProvider
		givenBetRepo   RepositoryProvider
		expectedBet    domain.Bet
	}{
		{
			name:           "given an ID, expect the associated bet to be returned",
			givenID:        uuid.MustParse("49cffe67-9798-4327-9760-c4b81562f928"),
			givenTableRepo: mockTableRepo{},
			givenBetRepo: mockBetRepo{
				GivenGetBet: domain.Bet{
					ID:             uuid.MustParse("49cffe67-9798-4327-9760-c4b81562f928"),
					Status:         "",
					SelectedSpaces: []int{5},
					Stake:          nil,
					PlacedAt:       time.Time{},
					SettledAt:      nil,
					Win:            false,
					Table:          uuid.MustParse("0173b64f-e07e-4fa0-bcb3-231856390dce"),
				},
			},
			expectedBet: domain.Bet{
				ID:             uuid.MustParse("49cffe67-9798-4327-9760-c4b81562f928"),
				Status:         "",
				SelectedSpaces: []int{5},
				Stake:          nil,
				PlacedAt:       time.Time{},
				SettledAt:      nil,
				Win:            false,
				Table:          uuid.MustParse("0173b64f-e07e-4fa0-bcb3-231856390dce"),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := NewController(test.givenBetRepo, test.givenTableRepo)

			actual, err := c.Get(context.Background(), test.givenID)
			if err != nil {
				t.Fatal(err)
			}

			if !cmp.Equal(actual, test.expectedBet) {
				t.Fatal(cmp.Diff(actual, test.expectedBet))
			}
		})
	}
}

func TestController_Get_Fail(t *testing.T) {
	tests := []struct {
		name           string
		givenID        uuid.UUID
		givenTableRepo TableRepoProvider
		givenBetRepo   RepositoryProvider
		expectedError  error
	}{
		{
			name:           "given a bet repo error, expect it to be returned",
			givenID:        uuid.MustParse("49cffe67-9798-4327-9760-c4b81562f928"),
			givenTableRepo: mockTableRepo{},
			givenBetRepo: mockBetRepo{
				GivenGetError: memory.ErrInvalidKey,
			},
			expectedError: memory.ErrInvalidKey,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := NewController(test.givenBetRepo, test.givenTableRepo)

			_, err := c.Get(context.Background(), test.givenID)
			if err == nil {
				t.Fatalf("expected %v, got nil", test.expectedError)
			}

			if !cmp.Equal(err, test.expectedError, cmpopts.EquateErrors()) {
				t.Fatal(cmp.Diff(err, test.expectedError, cmpopts.EquateErrors()))
			}
		})
	}
}

type mockBetRepo struct {
	GivenGetBet      domain.Bet
	GivenGetError    error
	GivenInsertError error
}

func (m mockBetRepo) Get(_ context.Context, _ uuid.UUID) (domain.Bet, error) {
	return m.GivenGetBet, m.GivenGetError
}

func (m mockBetRepo) Insert(_ context.Context, _ domain.Bet) error {
	return m.GivenInsertError
}

type mockTableRepo struct {
	GivenGetTable domain.Table
	GivenGetError error
}

func (m mockTableRepo) Get(_ context.Context, _ uuid.UUID) (domain.Table, error) {
	return m.GivenGetTable, m.GivenGetError
}
