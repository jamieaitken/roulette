package memory

import (
	"betting/internal/domain"
	"betting/storage"
	"betting/testing/opts"
	"context"
	"sync"
	"testing"

	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/google/go-cmp/cmp"

	"github.com/Rhymond/go-money"
	"github.com/google/uuid"
)

func TestBetStorage_Get_Success(t *testing.T) {
	tests := []struct {
		name        string
		givenBets   map[uuid.UUID]storage.Bet
		givenID     uuid.UUID
		expectedBet storage.Bet
	}{
		{
			name: "given a valid ID, expect associated bet returned",
			givenBets: map[uuid.UUID]storage.Bet{
				uuid.MustParse("22ee17b5-fae7-4c13-80cc-4354820df3d4"): {
					ID:    uuid.MustParse("22ee17b5-fae7-4c13-80cc-4354820df3d4"),
					Stake: money.New(100, "GBP"),
				},
			},
			givenID: uuid.MustParse("22ee17b5-fae7-4c13-80cc-4354820df3d4"),
			expectedBet: storage.Bet{
				ID:    uuid.MustParse("22ee17b5-fae7-4c13-80cc-4354820df3d4"),
				Stake: money.New(100, "GBP"),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store := BetStorage{
				bets:    test.givenBets,
				RWMutex: sync.RWMutex{},
			}

			actual, err := store.Get(context.Background(), test.givenID)
			if err != nil {
				t.Fatal(err)
			}

			if !cmp.Equal(actual, test.expectedBet, opts.MoneyComparer) {
				t.Fatal(cmp.Diff(actual, test.expectedBet, opts.MoneyComparer))
			}
		})
	}
}

func TestBetStorage_Get_Fail(t *testing.T) {
	tests := []struct {
		name          string
		givenBets     map[uuid.UUID]storage.Bet
		givenID       uuid.UUID
		expectedError error
	}{
		{
			name: "given an invalid ID, expect ErrInvalidKey",
			givenBets: map[uuid.UUID]storage.Bet{
				uuid.MustParse("22ee17b5-fae7-4c13-80cc-4354820df3d4"): {
					ID:    uuid.MustParse("22ee17b5-fae7-4c13-80cc-4354820df3d4"),
					Stake: money.New(100, "GBP"),
				},
			},
			givenID:       uuid.MustParse("33ee17b5-fae7-4c13-80cc-4354820df3d4"),
			expectedError: ErrInvalidKey,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store := BetStorage{
				bets:    test.givenBets,
				RWMutex: sync.RWMutex{},
			}

			_, err := store.Get(context.Background(), test.givenID)
			if err == nil {
				t.Fatalf("expected %v, got nil", test.expectedError)
			}

			if !cmp.Equal(err, test.expectedError, cmpopts.EquateErrors()) {
				t.Fatal(cmp.Diff(err, test.expectedError, cmpopts.EquateErrors()))
			}
		})
	}
}

func TestBetStorage_Insert_Success(t *testing.T) {
	tests := []struct {
		name      string
		givenBets map[uuid.UUID]storage.Bet
		givenBet  storage.Bet
	}{
		{
			name:      "given a bet which has not previously been inserted, expect success",
			givenBets: map[uuid.UUID]storage.Bet{},
			givenBet: storage.Bet{
				ID:    uuid.MustParse("22ee17b5-fae7-4c13-80cc-4354820df3d4"),
				Stake: money.New(100, "GBP"),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store := NewBetStorage()

			err := store.Insert(context.Background(), test.givenBet)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestBetStorage_Insert_Fail(t *testing.T) {
	tests := []struct {
		name          string
		givenBets     map[uuid.UUID]storage.Bet
		givenBet      storage.Bet
		expectedError error
	}{
		{
			name: "given a bet which has already been inserted, expect ErrDuplicateKey",
			givenBets: map[uuid.UUID]storage.Bet{
				uuid.MustParse("22ee17b5-fae7-4c13-80cc-4354820df3d4"): {
					ID:    uuid.MustParse("22ee17b5-fae7-4c13-80cc-4354820df3d4"),
					Stake: money.New(100, "GBP"),
				},
			},
			givenBet: storage.Bet{
				ID:    uuid.MustParse("22ee17b5-fae7-4c13-80cc-4354820df3d4"),
				Stake: money.New(100, "GBP"),
			},
			expectedError: ErrDuplicateKey,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store := BetStorage{
				bets:    test.givenBets,
				RWMutex: sync.RWMutex{},
			}

			err := store.Insert(context.Background(), test.givenBet)
			if err == nil {
				t.Fatalf("expected %v, got nil", test.expectedError)
			}

			if !cmp.Equal(err, test.expectedError, cmpopts.EquateErrors()) {
				t.Fatal(cmp.Diff(err, test.expectedError, cmpopts.EquateErrors()))
			}
		})
	}
}

func TestBetStorage_List(t *testing.T) {
	tests := []struct {
		name         string
		givenBets    map[uuid.UUID]storage.Bet
		givenTableID uuid.UUID
		expectedBets []storage.Bet
	}{
		{
			name: "given a table ID, expected associated bets returned",
			givenBets: map[uuid.UUID]storage.Bet{
				uuid.MustParse("22ee17b5-fae7-4c13-80cc-4354820df3d4"): {
					ID:    uuid.MustParse("22ee17b5-fae7-4c13-80cc-4354820df3d4"),
					Stake: money.New(100, "GBP"),
					Table: uuid.MustParse("1a31c7a1-6577-44c6-b3be-829674bf5175"),
				},
				uuid.MustParse("31b128ac-37de-4b24-99ca-1f3646798f41"): {
					ID:    uuid.MustParse("31b128ac-37de-4b24-99ca-1f3646798f41"),
					Stake: money.New(200, "GBP"),
					Table: uuid.MustParse("78e7a130-761e-4188-a204-715e3ab747a3"),
				},
			},
			givenTableID: uuid.MustParse("1a31c7a1-6577-44c6-b3be-829674bf5175"),
			expectedBets: []storage.Bet{
				{
					ID:    uuid.MustParse("22ee17b5-fae7-4c13-80cc-4354820df3d4"),
					Stake: money.New(100, "GBP"),
					Table: uuid.MustParse("1a31c7a1-6577-44c6-b3be-829674bf5175"),
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store := BetStorage{
				bets:    test.givenBets,
				RWMutex: sync.RWMutex{},
			}

			actual, err := store.List(context.Background(), test.givenTableID)
			if err != nil {
				t.Fatal(err)
			}

			if !cmp.Equal(actual, test.expectedBets, opts.MoneyComparer) {
				t.Fatal(cmp.Diff(actual, test.expectedBets, opts.MoneyComparer))
			}
		})
	}
}

func TestBetStorage_FindWinners(t *testing.T) {
	tests := []struct {
		name         string
		givenBets    map[uuid.UUID]storage.Bet
		bets         []storage.Bet
		givenTable   storage.Table
		expectedBets map[uuid.UUID]storage.Bet
	}{
		{
			name: "update one bet to signify they're a winner",
			givenBets: map[uuid.UUID]storage.Bet{
				uuid.MustParse("22ee17b5-fae7-4c13-80cc-4354820df3d4"): {
					ID:             uuid.MustParse("22ee17b5-fae7-4c13-80cc-4354820df3d4"),
					Stake:          money.New(100, "GBP"),
					Table:          uuid.MustParse("1a31c7a1-6577-44c6-b3be-829674bf5175"),
					SelectedSpaces: []int{12, 14},
				},
				uuid.MustParse("0438312a-cd6c-44b2-9c98-966b975e11d2"): {
					ID:             uuid.MustParse("0438312a-cd6c-44b2-9c98-966b975e11d2"),
					Stake:          money.New(100, "GBP"),
					Table:          uuid.MustParse("1a31c7a1-6577-44c6-b3be-829674bf5175"),
					SelectedSpaces: []int{12, 14, 16},
				},
				uuid.MustParse("31b128ac-37de-4b24-99ca-1f3646798f41"): {
					ID:    uuid.MustParse("31b128ac-37de-4b24-99ca-1f3646798f41"),
					Stake: money.New(200, "GBP"),
					Table: uuid.MustParse("78e7a130-761e-4188-a204-715e3ab747a3"),
				},
			},
			bets: []storage.Bet{
				{
					ID:             uuid.MustParse("22ee17b5-fae7-4c13-80cc-4354820df3d4"),
					Stake:          money.New(100, "GBP"),
					Table:          uuid.MustParse("1a31c7a1-6577-44c6-b3be-829674bf5175"),
					SelectedSpaces: []int{12, 14},
				},
				{
					ID:             uuid.MustParse("0438312a-cd6c-44b2-9c98-966b975e11d2"),
					Stake:          money.New(100, "GBP"),
					Table:          uuid.MustParse("1a31c7a1-6577-44c6-b3be-829674bf5175"),
					SelectedSpaces: []int{12, 14, 16},
					Win:            true,
				},
				{
					ID:    uuid.MustParse("31b128ac-37de-4b24-99ca-1f3646798f41"),
					Stake: money.New(200, "GBP"),
					Table: uuid.MustParse("78e7a130-761e-4188-a204-715e3ab747a3"),
				},
			},
			expectedBets: map[uuid.UUID]storage.Bet{
				uuid.MustParse("22ee17b5-fae7-4c13-80cc-4354820df3d4"): {
					ID:             uuid.MustParse("22ee17b5-fae7-4c13-80cc-4354820df3d4"),
					Stake:          money.New(100, "GBP"),
					Table:          uuid.MustParse("1a31c7a1-6577-44c6-b3be-829674bf5175"),
					SelectedSpaces: []int{12, 14},
				},
				uuid.MustParse("0438312a-cd6c-44b2-9c98-966b975e11d2"): {
					ID:             uuid.MustParse("0438312a-cd6c-44b2-9c98-966b975e11d2"),
					Stake:          money.New(100, "GBP"),
					Table:          uuid.MustParse("1a31c7a1-6577-44c6-b3be-829674bf5175"),
					SelectedSpaces: []int{12, 14, 16},
					Win:            true,
				},
				uuid.MustParse("31b128ac-37de-4b24-99ca-1f3646798f41"): {
					ID:    uuid.MustParse("31b128ac-37de-4b24-99ca-1f3646798f41"),
					Stake: money.New(200, "GBP"),
					Table: uuid.MustParse("78e7a130-761e-4188-a204-715e3ab747a3"),
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store := BetStorage{
				bets:    test.givenBets,
				RWMutex: sync.RWMutex{},
			}

			err := store.SetWinners(context.Background(), test.bets)
			if err != nil {
				t.Fatal(err)
			}

			if !cmp.Equal(store.bets, test.expectedBets, opts.MoneyComparer) {
				t.Fatal(cmp.Diff(store.bets, test.expectedBets, opts.MoneyComparer))
			}
		})
	}
}

func TestBetStorage_UpdateStateByTableID(t *testing.T) {
	tests := []struct {
		name         string
		givenBets    map[uuid.UUID]storage.Bet
		givenTableID uuid.UUID
		givenStatus  domain.BetStatus
		expectedBets map[uuid.UUID]storage.Bet
	}{
		{
			name: "update two bets to live",
			givenBets: map[uuid.UUID]storage.Bet{
				uuid.MustParse("22ee17b5-fae7-4c13-80cc-4354820df3d4"): {
					ID:             uuid.MustParse("22ee17b5-fae7-4c13-80cc-4354820df3d4"),
					Stake:          money.New(100, "GBP"),
					Table:          uuid.MustParse("1a31c7a1-6577-44c6-b3be-829674bf5175"),
					SelectedSpaces: []int{12, 14},
				},
				uuid.MustParse("0438312a-cd6c-44b2-9c98-966b975e11d2"): {
					ID:             uuid.MustParse("0438312a-cd6c-44b2-9c98-966b975e11d2"),
					Stake:          money.New(100, "GBP"),
					Table:          uuid.MustParse("1a31c7a1-6577-44c6-b3be-829674bf5175"),
					SelectedSpaces: []int{12, 14, 16},
				},
				uuid.MustParse("31b128ac-37de-4b24-99ca-1f3646798f41"): {
					ID:    uuid.MustParse("31b128ac-37de-4b24-99ca-1f3646798f41"),
					Stake: money.New(200, "GBP"),
					Table: uuid.MustParse("78e7a130-761e-4188-a204-715e3ab747a3"),
				},
			},
			givenTableID: uuid.MustParse("1a31c7a1-6577-44c6-b3be-829674bf5175"),
			givenStatus:  domain.Live,
			expectedBets: map[uuid.UUID]storage.Bet{
				uuid.MustParse("22ee17b5-fae7-4c13-80cc-4354820df3d4"): {
					ID:             uuid.MustParse("22ee17b5-fae7-4c13-80cc-4354820df3d4"),
					Stake:          money.New(100, "GBP"),
					Table:          uuid.MustParse("1a31c7a1-6577-44c6-b3be-829674bf5175"),
					SelectedSpaces: []int{12, 14},
					Status:         "live",
				},
				uuid.MustParse("0438312a-cd6c-44b2-9c98-966b975e11d2"): {
					ID:             uuid.MustParse("0438312a-cd6c-44b2-9c98-966b975e11d2"),
					Stake:          money.New(100, "GBP"),
					Table:          uuid.MustParse("1a31c7a1-6577-44c6-b3be-829674bf5175"),
					SelectedSpaces: []int{12, 14, 16},
					Status:         "live",
				},
				uuid.MustParse("31b128ac-37de-4b24-99ca-1f3646798f41"): {
					ID:    uuid.MustParse("31b128ac-37de-4b24-99ca-1f3646798f41"),
					Stake: money.New(200, "GBP"),
					Table: uuid.MustParse("78e7a130-761e-4188-a204-715e3ab747a3"),
				},
			},
		},
		{
			name: "update two bets to settled",
			givenBets: map[uuid.UUID]storage.Bet{
				uuid.MustParse("22ee17b5-fae7-4c13-80cc-4354820df3d4"): {
					ID:             uuid.MustParse("22ee17b5-fae7-4c13-80cc-4354820df3d4"),
					Stake:          money.New(100, "GBP"),
					Table:          uuid.MustParse("1a31c7a1-6577-44c6-b3be-829674bf5175"),
					SelectedSpaces: []int{12, 14},
				},
				uuid.MustParse("0438312a-cd6c-44b2-9c98-966b975e11d2"): {
					ID:             uuid.MustParse("0438312a-cd6c-44b2-9c98-966b975e11d2"),
					Stake:          money.New(100, "GBP"),
					Table:          uuid.MustParse("1a31c7a1-6577-44c6-b3be-829674bf5175"),
					SelectedSpaces: []int{12, 14, 16},
				},
				uuid.MustParse("31b128ac-37de-4b24-99ca-1f3646798f41"): {
					ID:    uuid.MustParse("31b128ac-37de-4b24-99ca-1f3646798f41"),
					Stake: money.New(200, "GBP"),
					Table: uuid.MustParse("78e7a130-761e-4188-a204-715e3ab747a3"),
				},
			},
			givenTableID: uuid.MustParse("1a31c7a1-6577-44c6-b3be-829674bf5175"),
			givenStatus:  domain.Settled,
			expectedBets: map[uuid.UUID]storage.Bet{
				uuid.MustParse("22ee17b5-fae7-4c13-80cc-4354820df3d4"): {
					ID:             uuid.MustParse("22ee17b5-fae7-4c13-80cc-4354820df3d4"),
					Stake:          money.New(100, "GBP"),
					Table:          uuid.MustParse("1a31c7a1-6577-44c6-b3be-829674bf5175"),
					SelectedSpaces: []int{12, 14},
					Status:         "settled",
				},
				uuid.MustParse("0438312a-cd6c-44b2-9c98-966b975e11d2"): {
					ID:             uuid.MustParse("0438312a-cd6c-44b2-9c98-966b975e11d2"),
					Stake:          money.New(100, "GBP"),
					Table:          uuid.MustParse("1a31c7a1-6577-44c6-b3be-829674bf5175"),
					SelectedSpaces: []int{12, 14, 16},
					Status:         "settled",
				},
				uuid.MustParse("31b128ac-37de-4b24-99ca-1f3646798f41"): {
					ID:    uuid.MustParse("31b128ac-37de-4b24-99ca-1f3646798f41"),
					Stake: money.New(200, "GBP"),
					Table: uuid.MustParse("78e7a130-761e-4188-a204-715e3ab747a3"),
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store := BetStorage{
				bets:    test.givenBets,
				RWMutex: sync.RWMutex{},
			}

			err := store.UpdateStateByTableID(context.Background(), test.givenTableID, test.givenStatus)
			if err != nil {
				t.Fatal(err)
			}

			if !cmp.Equal(store.bets, test.expectedBets, opts.MoneyComparer, cmpopts.IgnoreFields(storage.Bet{}, "SettledAt")) {
				t.Fatal(cmp.Diff(store.bets, test.expectedBets, opts.MoneyComparer, cmpopts.IgnoreFields(storage.Bet{}, "SettledAt")))
			}
		})
	}
}
