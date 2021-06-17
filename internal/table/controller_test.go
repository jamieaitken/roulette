package table

import (
	"betting/internal/domain"
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/google/uuid"
)

func TestController_Create_Success(t *testing.T) {
	tests := []struct {
		name               string
		givenRepository    RepositoryProvider
		givenBetRepository BetRepositoryProvider
		givenBallPlacer    BallPlacer
		givenLocator       WinnerLocator
		expectedTable      domain.Table
	}{
		{
			name:               "expected the created table to be returned",
			givenRepository:    mockTableRepositoryProvider{},
			givenBetRepository: mockBetRepository{},
			givenBallPlacer:    mockBallPlacer{},
			givenLocator:       mockLocator{},
			expectedTable: domain.Table{
				Bets:     nil,
				IsClosed: false,
				Outcome:  nil,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			controller := NewController(ControllerParams{
				RepositoryProvider:    test.givenRepository,
				BallPlacer:            test.givenBallPlacer,
				WinnerLocator:         test.givenLocator,
				BetRepositoryProvider: test.givenBetRepository,
			})

			actual, err := controller.Create(context.Background())
			if err != nil {
				t.Fatal(err)
			}

			if !cmp.Equal(actual, test.expectedTable, cmpopts.IgnoreTypes(uuid.UUID{})) {
				t.Fatal(cmp.Diff(actual, test.expectedTable, cmpopts.IgnoreTypes(uuid.UUID{})))
			}
		})
	}
}

func TestController_Create_Fail(t *testing.T) {
	tests := []struct {
		name               string
		givenRepository    RepositoryProvider
		givenBetRepository BetRepositoryProvider
		givenBallPlacer    BallPlacer
		givenLocator       WinnerLocator
		expectedError      error
	}{
		{
			name: "given a repo error, expect error to be returned",
			givenRepository: mockTableRepositoryProvider{
				GivenInsertError: ErrFailedToCreateTable,
			},
			givenBallPlacer: mockBallPlacer{},
			expectedError:   ErrFailedToCreateTable,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			controller := NewController(ControllerParams{
				RepositoryProvider:    test.givenRepository,
				BallPlacer:            test.givenBallPlacer,
				WinnerLocator:         test.givenLocator,
				BetRepositoryProvider: test.givenBetRepository,
			})

			_, err := controller.Create(context.Background())
			if err == nil {
				t.Fatalf("expected %v, got nil", test.expectedError)
			}

			if !cmp.Equal(err, test.expectedError, cmpopts.EquateErrors()) {
				t.Fatal(cmp.Diff(err, test.expectedError, cmpopts.EquateErrors()))
			}
		})
	}
}

func TestController_Spin_Success(t *testing.T) {
	tests := []struct {
		name               string
		givenRepository    RepositoryProvider
		givenBetRepository BetRepositoryProvider
		givenBallPlacer    BallPlacer
		givenLocator       WinnerLocator
		givenID            uuid.UUID
		expectedTable      domain.Table
	}{
		{
			name: "expected the created table to be returned",
			givenRepository: mockTableRepositoryProvider{
				GivenGetTable: domain.Table{
					ID:       uuid.MustParse("160998da-2d89-4f06-a690-fd189213958d"),
					Bets:     nil,
					IsClosed: false,
					Outcome: &domain.Outcome{
						Value:  16,
						Colour: domain.Red,
					},
				},
			},
			givenBetRepository: mockBetRepository{},
			givenBallPlacer: mockBallPlacer{
				GivenOutcome: domain.Outcome{
					Value:  16,
					Colour: domain.Red,
				},
			},
			givenLocator: mockLocator{
				GivenTable: domain.Table{
					ID:       uuid.MustParse("160998da-2d89-4f06-a690-fd189213958d"),
					Bets:     nil,
					IsClosed: false,
					Outcome: &domain.Outcome{
						Value:  16,
						Colour: domain.Red,
					},
				},
			},
			givenID: uuid.MustParse("160998da-2d89-4f06-a690-fd189213958d"),
			expectedTable: domain.Table{
				ID:       uuid.MustParse("160998da-2d89-4f06-a690-fd189213958d"),
				Bets:     nil,
				IsClosed: false,
				Outcome: &domain.Outcome{
					Value:  16,
					Colour: domain.Red,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			controller := NewController(ControllerParams{
				RepositoryProvider:    test.givenRepository,
				BallPlacer:            test.givenBallPlacer,
				WinnerLocator:         test.givenLocator,
				BetRepositoryProvider: test.givenBetRepository,
			})

			actual, err := controller.Spin(context.Background(), test.givenID)
			if err != nil {
				t.Fatal(err)
			}

			if !cmp.Equal(actual, test.expectedTable) {
				t.Fatal(cmp.Diff(actual, test.expectedTable))
			}
		})
	}
}

func TestController_Spin_Fail(t *testing.T) {
	tests := []struct {
		name               string
		givenRepository    RepositoryProvider
		givenBetRepository BetRepositoryProvider
		givenBallPlacer    BallPlacer
		givenLocator       WinnerLocator
		givenID            uuid.UUID
		expectedError      error
	}{
		{
			name: "given a repo close error, expect error to be returned",
			givenRepository: mockTableRepositoryProvider{
				GivenCloseError: ErrFailedToCreateTable,
			},
			givenBallPlacer:    mockBallPlacer{},
			givenBetRepository: mockBetRepository{},
			expectedError:      ErrFailedToCloseTable,
		},
		{
			name:            "given a repo spin error, expect error to be returned",
			givenRepository: mockTableRepositoryProvider{},
			givenBetRepository: mockBetRepository{
				GivenSpinError: ErrFailedToSpinTable,
			},
			givenBallPlacer: mockBallPlacer{},
			expectedError:   ErrFailedToSpinTable,
		},
		{
			name: "given a repo set outcome error, expect error to be returned",
			givenRepository: mockTableRepositoryProvider{
				GivenSetOutcomeError: ErrFailedToCreateTable,
			},
			givenBetRepository: mockBetRepository{},
			givenBallPlacer:    mockBallPlacer{},
			expectedError:      ErrFailedToSetOutcome,
		},
		{
			name: "given a repo get error, expect error to be returned",
			givenRepository: mockTableRepositoryProvider{
				GivenGetError: ErrFailedToCreateTable,
			},
			givenBetRepository: mockBetRepository{},
			givenBallPlacer:    mockBallPlacer{},
			expectedError:      ErrFailedToFetchTable,
		},
		{
			name:            "given a bet repo list error, expect error to be returned",
			givenRepository: mockTableRepositoryProvider{},
			givenBetRepository: mockBetRepository{
				GivenListError: ErrFailedToFetchBets,
			},
			givenBallPlacer: mockBallPlacer{},
			expectedError:   ErrFailedToFetchBets,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			controller := NewController(ControllerParams{
				RepositoryProvider:    test.givenRepository,
				BallPlacer:            test.givenBallPlacer,
				WinnerLocator:         test.givenLocator,
				BetRepositoryProvider: test.givenBetRepository,
			})

			_, err := controller.Spin(context.Background(), test.givenID)
			if err == nil {
				t.Fatalf("expected %v, got nil", test.expectedError)
			}

			if !cmp.Equal(err, test.expectedError, cmpopts.EquateErrors()) {
				t.Fatal(cmp.Diff(err, test.expectedError, cmpopts.EquateErrors()))
			}
		})
	}
}

func TestController_Settle_Success(t *testing.T) {
	tests := []struct {
		name               string
		givenRepository    RepositoryProvider
		givenBetRepository BetRepositoryProvider
		givenBallPlacer    BallPlacer
		givenLocator       WinnerLocator
		givenID            uuid.UUID
		expectedTable      domain.Table
	}{
		{
			name: "expect all bets to be settled for the given table",
			givenRepository: mockTableRepositoryProvider{
				GivenGetTable: domain.Table{
					ID:       uuid.MustParse("160998da-2d89-4f06-a690-fd189213958d"),
					Bets:     nil,
					IsClosed: false,
					Outcome: &domain.Outcome{
						Value:  16,
						Colour: domain.Red,
					},
				},
			},
			givenLocator: mockLocator{
				GivenTable: domain.Table{
					ID:       uuid.MustParse("160998da-2d89-4f06-a690-fd189213958d"),
					Bets:     nil,
					IsClosed: false,
					Outcome: &domain.Outcome{
						Value:  16,
						Colour: domain.Red,
					},
				},
			},
			givenBetRepository: mockBetRepository{},
			givenID:            uuid.MustParse("160998da-2d89-4f06-a690-fd189213958d"),
			expectedTable: domain.Table{
				ID:       uuid.MustParse("160998da-2d89-4f06-a690-fd189213958d"),
				Bets:     nil,
				IsClosed: false,
				Outcome: &domain.Outcome{
					Value:  16,
					Colour: domain.Red,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			controller := NewController(ControllerParams{
				RepositoryProvider:    test.givenRepository,
				BallPlacer:            test.givenBallPlacer,
				WinnerLocator:         test.givenLocator,
				BetRepositoryProvider: test.givenBetRepository,
			})

			actual, err := controller.Settle(context.Background(), test.givenID)
			if err != nil {
				t.Fatal(err)
			}

			if !cmp.Equal(actual, test.expectedTable) {
				t.Fatal(cmp.Diff(actual, test.expectedTable))
			}
		})
	}
}

func TestController_Settle_Fail(t *testing.T) {
	tests := []struct {
		name               string
		givenRepository    RepositoryProvider
		givenBetRepository BetRepositoryProvider
		givenBallPlacer    BallPlacer
		givenLocator       WinnerLocator
		givenID            uuid.UUID
		expectedError      error
	}{
		{
			name:            "given a repo settle error, expect error to be returned",
			givenRepository: mockTableRepositoryProvider{},
			givenBetRepository: mockBetRepository{
				GivenSettleError: ErrFailedFailedToSettle,
			},
			givenBallPlacer: mockBallPlacer{},
			expectedError:   ErrFailedFailedToSettle,
		},
		{
			name: "given a repo get error, expect error to be returned",
			givenRepository: mockTableRepositoryProvider{
				GivenGetError: ErrFailedToFetchTable,
			},
			givenBetRepository: mockBetRepository{},
			givenBallPlacer:    mockBallPlacer{},
			expectedError:      ErrFailedToFetchTable,
		},
		{
			name:            "given a repo set winners error, expect error to be returned",
			givenRepository: mockTableRepositoryProvider{},
			givenBetRepository: mockBetRepository{
				GivenSetWinnersError: ErrFailedToSetWinners,
			},
			givenLocator:    mockLocator{},
			givenBallPlacer: mockBallPlacer{},
			expectedError:   ErrFailedToSetWinners,
		},
		{
			name:            "given a bet repo list error, expect error to be returned",
			givenRepository: mockTableRepositoryProvider{},
			givenBetRepository: mockBetRepository{
				GivenListError: ErrFailedToFetchBets,
			},
			givenLocator:    mockLocator{},
			givenBallPlacer: mockBallPlacer{},
			expectedError:   ErrFailedToFetchBets,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			controller := NewController(ControllerParams{
				RepositoryProvider:    test.givenRepository,
				BallPlacer:            test.givenBallPlacer,
				WinnerLocator:         test.givenLocator,
				BetRepositoryProvider: test.givenBetRepository,
			})

			_, err := controller.Settle(context.Background(), test.givenID)
			if err == nil {
				t.Fatalf("expected %v, got nil", test.expectedError)
			}

			if !cmp.Equal(err, test.expectedError, cmpopts.EquateErrors()) {
				t.Fatal(cmp.Diff(err, test.expectedError, cmpopts.EquateErrors()))
			}
		})
	}
}

func TestController_List_Success(t *testing.T) {
	tests := []struct {
		name               string
		givenRepository    RepositoryProvider
		givenBetRepository BetRepositoryProvider
		givenBallPlacer    BallPlacer
		givenLocator       WinnerLocator
		givenID            uuid.UUID
		expectedTables     []domain.Table
	}{
		{
			name: "expect all tables to be returned",
			givenRepository: mockTableRepositoryProvider{
				GivenListTables: []domain.Table{
					{
						ID:       uuid.MustParse("160998da-2d89-4f06-a690-fd189213958d"),
						Bets:     nil,
						IsClosed: false,
						Outcome: &domain.Outcome{
							Value:  16,
							Colour: domain.Red,
						},
					},
				},
			},
			givenLocator: mockLocator{},
			givenBetRepository: mockBetRepository{
				GivenListBets: []domain.Bet{
					{
						ID:    uuid.MustParse("e49779f6-3507-4063-bed8-18d50174868d"),
						Table: uuid.MustParse("160998da-2d89-4f06-a690-fd189213958d"),
					},
				},
			},
			givenID: uuid.MustParse("160998da-2d89-4f06-a690-fd189213958d"),
			expectedTables: []domain.Table{
				{
					ID: uuid.MustParse("160998da-2d89-4f06-a690-fd189213958d"),
					Bets: []domain.Bet{
						{
							ID:    uuid.MustParse("e49779f6-3507-4063-bed8-18d50174868d"),
							Table: uuid.MustParse("160998da-2d89-4f06-a690-fd189213958d"),
						},
					},
					IsClosed: false,
					Outcome: &domain.Outcome{
						Value:  16,
						Colour: domain.Red,
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			controller := NewController(ControllerParams{
				RepositoryProvider:    test.givenRepository,
				BallPlacer:            test.givenBallPlacer,
				WinnerLocator:         test.givenLocator,
				BetRepositoryProvider: test.givenBetRepository,
			})

			actual, err := controller.List(context.Background())
			if err != nil {
				t.Fatal(err)
			}

			if !cmp.Equal(actual, test.expectedTables) {
				t.Fatal(cmp.Diff(actual, test.expectedTables))
			}
		})
	}
}

func TestController_List_Fail(t *testing.T) {
	tests := []struct {
		name               string
		givenRepository    RepositoryProvider
		givenBetRepository BetRepositoryProvider
		givenBallPlacer    BallPlacer
		givenLocator       WinnerLocator
		givenID            uuid.UUID
		expectedError      error
	}{
		{
			name: "given repo list error, expect it to be returned",
			givenRepository: mockTableRepositoryProvider{
				GivenListError: ErrFailedToFetchTable,
			},
			givenLocator:       mockLocator{},
			givenBetRepository: mockBetRepository{},
			givenID:            uuid.MustParse("160998da-2d89-4f06-a690-fd189213958d"),
			expectedError:      ErrFailedToFetchTable,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			controller := NewController(ControllerParams{
				RepositoryProvider:    test.givenRepository,
				BallPlacer:            test.givenBallPlacer,
				WinnerLocator:         test.givenLocator,
				BetRepositoryProvider: test.givenBetRepository,
			})

			_, err := controller.List(context.Background())
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
		name               string
		givenRepository    RepositoryProvider
		givenBetRepository BetRepositoryProvider
		givenBallPlacer    BallPlacer
		givenLocator       WinnerLocator
		givenID            uuid.UUID
		expectedTable      domain.Table
	}{
		{
			name: "expect table to be returned for given ID",
			givenRepository: mockTableRepositoryProvider{
				GivenGetTable: domain.Table{
					ID:       uuid.MustParse("160998da-2d89-4f06-a690-fd189213958d"),
					Bets:     nil,
					IsClosed: false,
					Outcome: &domain.Outcome{
						Value:  16,
						Colour: domain.Red,
					},
				},
			},
			givenLocator:       mockLocator{},
			givenBetRepository: mockBetRepository{},
			givenID:            uuid.MustParse("160998da-2d89-4f06-a690-fd189213958d"),
			expectedTable: domain.Table{
				ID:       uuid.MustParse("160998da-2d89-4f06-a690-fd189213958d"),
				Bets:     nil,
				IsClosed: false,
				Outcome: &domain.Outcome{
					Value:  16,
					Colour: domain.Red,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			controller := NewController(ControllerParams{
				RepositoryProvider:    test.givenRepository,
				BallPlacer:            test.givenBallPlacer,
				WinnerLocator:         test.givenLocator,
				BetRepositoryProvider: test.givenBetRepository,
			})

			actual, err := controller.Get(context.Background(), test.givenID)
			if err != nil {
				t.Fatal(err)
			}

			if !cmp.Equal(actual, test.expectedTable) {
				t.Fatal(cmp.Diff(actual, test.expectedTable))
			}
		})
	}
}

func TestController_Get_Fail(t *testing.T) {
	tests := []struct {
		name               string
		givenRepository    RepositoryProvider
		givenBetRepository BetRepositoryProvider
		givenBallPlacer    BallPlacer
		givenLocator       WinnerLocator
		givenID            uuid.UUID
		expectedError      error
	}{
		{
			name: "given a repo get error, expect error to be returned",
			givenRepository: mockTableRepositoryProvider{
				GivenGetError: ErrFailedToFetchTable,
			},
			givenBetRepository: mockBetRepository{},
			givenBallPlacer:    mockBallPlacer{},
			expectedError:      ErrFailedToFetchTable,
		},
		{
			name:            "given a bet repo list error, expect error to be returned",
			givenRepository: mockTableRepositoryProvider{},
			givenBetRepository: mockBetRepository{
				GivenListError: ErrFailedToFetchBets,
			},
			givenBallPlacer: mockBallPlacer{},
			expectedError:   ErrFailedToFetchBets,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			controller := NewController(ControllerParams{
				RepositoryProvider:    test.givenRepository,
				BallPlacer:            test.givenBallPlacer,
				WinnerLocator:         test.givenLocator,
				BetRepositoryProvider: test.givenBetRepository,
			})

			_, err := controller.Get(context.Background(), test.givenID)
			if err == nil {
				t.Fatalf("expected %v, got nil", test.expectedError)
			}

			if !cmp.Equal(err, test.expectedError, cmpopts.EquateErrors()) {
				t.Fatal(cmp.Diff(err, test.expectedError, cmpopts.EquateErrors()))
			}
		})
	}
}

type mockTableRepositoryProvider struct {
	GivenGetTable        domain.Table
	GivenGetError        error
	GivenListTables      []domain.Table
	GivenListError       error
	GivenInsertError     error
	GivenCloseError      error
	GivenSetOutcomeError error
}

type mockBetRepository struct {
	GivenSetWinnersError error
	GivenSettleError     error
	GivenSpinError       error
	GivenListBets        []domain.Bet
	GivenListError       error
}

func (m mockBetRepository) SetWinners(_ context.Context, _ []domain.Bet) error {
	return m.GivenSetWinnersError
}

func (m mockBetRepository) Spin(_ context.Context, _ uuid.UUID) error {
	return m.GivenSpinError
}

func (m mockBetRepository) Settle(_ context.Context, _ uuid.UUID) error {
	return m.GivenSettleError
}

func (m mockBetRepository) List(_ context.Context, _ uuid.UUID) ([]domain.Bet, error) {
	return m.GivenListBets, m.GivenListError
}

func (m mockTableRepositoryProvider) Get(_ context.Context, _ uuid.UUID) (domain.Table, error) {
	return m.GivenGetTable, m.GivenGetError
}

func (m mockTableRepositoryProvider) List(_ context.Context) ([]domain.Table, error) {
	return m.GivenListTables, m.GivenListError
}

func (m mockTableRepositoryProvider) Insert(_ context.Context, _ domain.Table) error {
	return m.GivenInsertError
}

func (m mockTableRepositoryProvider) Close(_ context.Context, _ uuid.UUID) error {
	return m.GivenCloseError
}

func (m mockTableRepositoryProvider) SetOutcome(_ context.Context, _ uuid.UUID, _ domain.Outcome) error {
	return m.GivenSetOutcomeError
}

type mockBallPlacer struct {
	GivenOutcome domain.Outcome
}

func (m mockBallPlacer) GetPosition(_ context.Context) domain.Outcome {
	return m.GivenOutcome
}

type mockLocator struct {
	GivenTable domain.Table
}

func (m mockLocator) Locate(_ context.Context, _ domain.Table) domain.Table {
	return m.GivenTable
}
