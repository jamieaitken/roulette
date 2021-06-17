package bet

import (
	"betting/api"
	"betting/internal/domain"
	"betting/internal/pkg/responses"
	"betting/storage/memory"
	"betting/testing/opts"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Rhymond/go-money"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func TestHandler_Create_Success(t *testing.T) {
	tests := []struct {
		name            string
		givenController Controller
		givenURL        string
		givenBody       api.BetRequest
		expectedStatus  int
		expectedBody    api.BetResponse
	}{
		{
			name: "given controller success, expect 201",
			givenController: mockController{
				GivenCreateBet: domain.Bet{
					ID:    uuid.MustParse("00812e8f-7fca-49a9-b141-9a52a0d0a82e"),
					Table: uuid.MustParse("bd88dfac-a3b9-43ee-ac7a-f958de23b26d"),
					Stake: money.New(100, "GBP"),
				},
			},
			givenURL: "/v1/tables/bd88dfac-a3b9-43ee-ac7a-f958de23b26d/bet",
			givenBody: api.BetRequest{
				SelectedSpaces: nil,
				Stake:          money.New(100, "GBP"),
				Table:          uuid.UUID{},
			},
			expectedStatus: http.StatusCreated,
			expectedBody: api.BetResponse{
				ID: uuid.MustParse("00812e8f-7fca-49a9-b141-9a52a0d0a82e"),
				BetRequest: api.BetRequest{
					SelectedSpaces: nil,
					Stake:          money.New(100, "GBP"),
					Table:          uuid.MustParse("bd88dfac-a3b9-43ee-ac7a-f958de23b26d"),
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := New(test.givenController)

			rr := httptest.NewRecorder()

			b, err := json.Marshal(test.givenBody)
			if err != nil {
				t.Fatal(err)
			}

			req := httptest.NewRequest(http.MethodPost, test.givenURL, bytes.NewReader(b))

			router := new(mux.Router)
			router.HandleFunc("/v1/tables/{id}/bet", handler.Create)
			router.ServeHTTP(rr, req)

			resp := rr.Result()

			if !cmp.Equal(resp.StatusCode, test.expectedStatus) {
				t.Fatal(cmp.Diff(resp.StatusCode, test.expectedStatus))
			}

			var res api.BetResponse
			err = json.NewDecoder(resp.Body).Decode(&res)
			if err != nil {
				t.Fatal(err)
			}

			if !cmp.Equal(res, test.expectedBody, opts.MoneyComparer) {
				t.Fatal(cmp.Diff(res, test.expectedBody, opts.MoneyComparer))
			}
		})
	}
}

func TestHandler_Create_Fail(t *testing.T) {
	tests := []struct {
		name            string
		givenController Controller
		givenURL        string
		givenBody       api.BetRequest
		expectedStatus  int
		expectedBody    responses.Error
	}{
		{
			name: "given an invalid id, expect 400",
			givenController: mockController{
				GivenCreateError: memory.ErrDuplicateKey,
			},
			givenURL:       "/v1/tables/test/bet",
			expectedStatus: http.StatusBadRequest,
			expectedBody: responses.Error{
				Status: http.StatusBadRequest,
				Detail: ErrInvalidID.Error(),
			},
		},
		{
			name: "given controller error, expect 400",
			givenController: mockController{
				GivenCreateError: memory.ErrDuplicateKey,
			},
			givenURL:       "/v1/tables/bd88dfac-a3b9-43ee-ac7a-f958de23b26d/bet",
			expectedStatus: http.StatusBadRequest,
			expectedBody: responses.Error{
				Status: http.StatusBadRequest,
				Detail: memory.ErrDuplicateKey.Error(),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := New(test.givenController)

			rr := httptest.NewRecorder()

			b, err := json.Marshal(test.givenBody)
			if err != nil {
				t.Fatal(err)
			}

			req := httptest.NewRequest(http.MethodPost, test.givenURL, bytes.NewReader(b))

			router := new(mux.Router)
			router.HandleFunc("/v1/tables/{id}/bet", handler.Create)
			router.ServeHTTP(rr, req)

			resp := rr.Result()

			if !cmp.Equal(resp.StatusCode, test.expectedStatus) {
				t.Fatal(cmp.Diff(resp.StatusCode, test.expectedStatus))
			}

			var res responses.Error
			err = json.NewDecoder(resp.Body).Decode(&res)
			if err != nil {
				t.Fatal(err)
			}

			if !cmp.Equal(res, test.expectedBody, opts.MoneyComparer) {
				t.Fatal(cmp.Diff(res, test.expectedBody, opts.MoneyComparer))
			}
		})
	}
}

func TestHandler_Get_Success(t *testing.T) {
	tests := []struct {
		name            string
		givenController Controller
		givenURL        string
		expectedStatus  int
		expectedBody    api.BetResponse
	}{
		{
			name: "given controller success, expect 200",
			givenController: mockController{
				GivenGetBet: domain.Bet{
					ID:    uuid.MustParse("00812e8f-7fca-49a9-b141-9a52a0d0a82e"),
					Table: uuid.MustParse("bd88dfac-a3b9-43ee-ac7a-f958de23b26d"),
					Stake: money.New(100, "GBP"),
				},
			},
			givenURL:       "/v1/bets/00812e8f-7fca-49a9-b141-9a52a0d0a82e",
			expectedStatus: http.StatusOK,
			expectedBody: api.BetResponse{
				ID: uuid.MustParse("00812e8f-7fca-49a9-b141-9a52a0d0a82e"),
				BetRequest: api.BetRequest{
					SelectedSpaces: nil,
					Stake:          money.New(100, "GBP"),
					Table:          uuid.MustParse("bd88dfac-a3b9-43ee-ac7a-f958de23b26d"),
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := New(test.givenController)

			rr := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodGet, test.givenURL, nil)

			router := new(mux.Router)
			router.HandleFunc("/v1/bets/{id}", handler.Get)
			router.ServeHTTP(rr, req)

			resp := rr.Result()

			if !cmp.Equal(resp.StatusCode, test.expectedStatus) {
				t.Fatal(cmp.Diff(resp.StatusCode, test.expectedStatus))
			}

			var res api.BetResponse
			err := json.NewDecoder(resp.Body).Decode(&res)
			if err != nil {
				t.Fatal(err)
			}

			if !cmp.Equal(res, test.expectedBody, opts.MoneyComparer) {
				t.Fatal(cmp.Diff(res, test.expectedBody, opts.MoneyComparer))
			}
		})
	}
}

func TestHandler_Get_Fail(t *testing.T) {
	tests := []struct {
		name            string
		givenController Controller
		givenURL        string
		expectedStatus  int
		expectedBody    responses.Error
	}{
		{
			name: "given an invalid id, expect 400",
			givenController: mockController{
				GivenGetError: memory.ErrDuplicateKey,
			},
			givenURL:       "/v1/bets/test",
			expectedStatus: http.StatusBadRequest,
			expectedBody: responses.Error{
				Status: http.StatusBadRequest,
				Detail: ErrInvalidID.Error(),
			},
		},
		{
			name: "given controller error, expect 400",
			givenController: mockController{
				GivenGetError: memory.ErrDuplicateKey,
			},
			givenURL:       "/v1/bets/00812e8f-7fca-49a9-b141-9a52a0d0a82e",
			expectedStatus: http.StatusBadRequest,
			expectedBody: responses.Error{
				Status: http.StatusBadRequest,
				Detail: memory.ErrDuplicateKey.Error(),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := New(test.givenController)

			rr := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodGet, test.givenURL, nil)

			router := new(mux.Router)
			router.HandleFunc("/v1/bets/{id}", handler.Get)
			router.ServeHTTP(rr, req)

			resp := rr.Result()

			if !cmp.Equal(resp.StatusCode, test.expectedStatus) {
				t.Fatal(cmp.Diff(resp.StatusCode, test.expectedStatus))
			}

			var res responses.Error
			err := json.NewDecoder(resp.Body).Decode(&res)
			if err != nil {
				t.Fatal(err)
			}

			if !cmp.Equal(res, test.expectedBody, opts.MoneyComparer) {
				t.Fatal(cmp.Diff(res, test.expectedBody, opts.MoneyComparer))
			}
		})
	}
}

type mockController struct {
	GivenCreateBet   domain.Bet
	GivenCreateError error
	GivenGetBet      domain.Bet
	GivenGetError    error
}

func (m mockController) Create(_ context.Context, _ domain.Bet) (domain.Bet, error) {
	return m.GivenCreateBet, m.GivenCreateError
}

func (m mockController) Get(_ context.Context, _ uuid.UUID) (domain.Bet, error) {
	return m.GivenGetBet, m.GivenGetError
}
