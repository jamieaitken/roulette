package table

import (
	"betting/api"
	"betting/internal/domain"
	"betting/internal/pkg/responses"
	"betting/storage/memory"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func TestHandler_Create_Success(t *testing.T) {
	tests := []struct {
		name            string
		givenController Controller
		givenURL        string
		expectedStatus  int
		expectedBody    api.TableResponse
	}{
		{
			name: "given controller success, expect 201",
			givenController: mockController{
				GivenCreateTable: domain.Table{
					ID:       uuid.MustParse("00812e8f-7fca-49a9-b141-9a52a0d0a82e"),
					Bets:     nil,
					IsClosed: false,
					Outcome:  nil,
				},
			},
			givenURL:       "/v1/tables",
			expectedStatus: http.StatusCreated,
			expectedBody: api.TableResponse{
				ID:       uuid.MustParse("00812e8f-7fca-49a9-b141-9a52a0d0a82e"),
				Bets:     []api.BetResponse{},
				IsClosed: false,
				Outcome:  nil,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := New(test.givenController)

			rr := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodPost, test.givenURL, nil)

			router := new(mux.Router)
			router.HandleFunc(test.givenURL, handler.Create)
			router.ServeHTTP(rr, req)

			resp := rr.Result()

			if !cmp.Equal(resp.StatusCode, test.expectedStatus) {
				t.Fatal(cmp.Diff(resp.StatusCode, test.expectedStatus))
			}

			var res api.TableResponse
			err := json.NewDecoder(resp.Body).Decode(&res)
			if err != nil {
				t.Fatal(err)
			}

			if !cmp.Equal(res, test.expectedBody) {
				t.Fatal(cmp.Diff(res, test.expectedBody))
			}
		})
	}
}

func TestHandler_Create_Fail(t *testing.T) {
	tests := []struct {
		name            string
		givenController Controller
		givenURL        string
		expectedStatus  int
		expectedBody    responses.Error
	}{
		{
			name: "given controller error, expect 400",
			givenController: mockController{
				GivenCreateError: memory.ErrDuplicateKey,
			},
			givenURL:       "/v1/tables",
			expectedStatus: http.StatusBadRequest,
			expectedBody: responses.Error{
				Status: http.StatusBadRequest,
				Detail: ErrFailedToCreateTable.Error(),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := New(test.givenController)

			rr := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodPost, test.givenURL, nil)

			router := new(mux.Router)
			router.HandleFunc(test.givenURL, handler.Create)
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

			if !cmp.Equal(res, test.expectedBody) {
				t.Fatal(cmp.Diff(res, test.expectedBody))
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
		expectedBody    api.TableResponse
	}{
		{
			name: "given controller success, expect 200",
			givenController: mockController{
				GivenGetTable: domain.Table{
					ID:       uuid.MustParse("00812e8f-7fca-49a9-b141-9a52a0d0a82e"),
					Bets:     nil,
					IsClosed: false,
					Outcome:  nil,
				},
			},
			givenURL:       "/v1/tables/00812e8f-7fca-49a9-b141-9a52a0d0a82e",
			expectedStatus: http.StatusOK,
			expectedBody: api.TableResponse{
				ID:       uuid.MustParse("00812e8f-7fca-49a9-b141-9a52a0d0a82e"),
				Bets:     []api.BetResponse{},
				IsClosed: false,
				Outcome:  nil,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := New(test.givenController)

			rr := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodGet, test.givenURL, nil)

			router := new(mux.Router)
			router.HandleFunc("/v1/tables/{id}", handler.Get)
			router.ServeHTTP(rr, req)

			resp := rr.Result()

			if !cmp.Equal(resp.StatusCode, test.expectedStatus) {
				t.Fatal(cmp.Diff(resp.StatusCode, test.expectedStatus))
			}

			var res api.TableResponse
			err := json.NewDecoder(resp.Body).Decode(&res)
			if err != nil {
				t.Fatal(err)
			}

			if !cmp.Equal(res, test.expectedBody) {
				t.Fatal(cmp.Diff(res, test.expectedBody))
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
			name:            "given invalid ID, expect 400",
			givenController: mockController{},
			givenURL:        "/v1/tables/test",
			expectedStatus:  http.StatusBadRequest,
			expectedBody: responses.Error{
				Status: http.StatusBadRequest,
				Detail: ErrInvalidID.Error(),
			},
		},
		{
			name: "given controller error, expect 400",
			givenController: mockController{
				GivenGetError: memory.ErrInvalidKey,
			},
			givenURL:       "/v1/tables/70ee9bba-87ac-4155-8ec7-f83c8663315e",
			expectedStatus: http.StatusBadRequest,
			expectedBody: responses.Error{
				Status: http.StatusBadRequest,
				Detail: memory.ErrInvalidKey.Error(),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := New(test.givenController)

			rr := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodGet, test.givenURL, nil)

			router := new(mux.Router)
			router.HandleFunc("/v1/tables/{id}", handler.Get)
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

			if !cmp.Equal(res, test.expectedBody) {
				t.Fatal(cmp.Diff(res, test.expectedBody))
			}
		})
	}
}

func TestHandler_Spin_Success(t *testing.T) {
	tests := []struct {
		name            string
		givenController Controller
		givenURL        string
		expectedStatus  int
		expectedBody    api.TableResponse
	}{
		{
			name: "given controller success, expect 200",
			givenController: mockController{
				GivenSpinTable: domain.Table{
					ID:       uuid.MustParse("00812e8f-7fca-49a9-b141-9a52a0d0a82e"),
					Bets:     nil,
					IsClosed: true,
					Outcome: &domain.Outcome{
						Value:  16,
						Colour: domain.Red,
					},
				},
			},
			givenURL:       "/v1/tables/00812e8f-7fca-49a9-b141-9a52a0d0a82e/spin",
			expectedStatus: http.StatusOK,
			expectedBody: api.TableResponse{
				ID:       uuid.MustParse("00812e8f-7fca-49a9-b141-9a52a0d0a82e"),
				Bets:     []api.BetResponse{},
				IsClosed: true,
				Outcome: &api.Outcome{
					Position: 16,
					Colour:   domain.Red,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := New(test.givenController)

			rr := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodPut, test.givenURL, nil)

			router := new(mux.Router)
			router.HandleFunc("/v1/tables/{id}/spin", handler.Spin)
			router.ServeHTTP(rr, req)

			resp := rr.Result()

			if !cmp.Equal(resp.StatusCode, test.expectedStatus) {
				t.Fatal(cmp.Diff(resp.StatusCode, test.expectedStatus))
			}

			var res api.TableResponse
			err := json.NewDecoder(resp.Body).Decode(&res)
			if err != nil {
				t.Fatal(err)
			}

			if !cmp.Equal(res, test.expectedBody) {
				t.Fatal(cmp.Diff(res, test.expectedBody))
			}
		})
	}
}

func TestHandler_Spin_Fail(t *testing.T) {
	tests := []struct {
		name            string
		givenController Controller
		givenURL        string
		expectedStatus  int
		expectedBody    responses.Error
	}{
		{
			name:            "given invalid ID, expect 400",
			givenController: mockController{},
			givenURL:        "/v1/tables/test/spin",
			expectedStatus:  http.StatusBadRequest,
			expectedBody: responses.Error{
				Status: http.StatusBadRequest,
				Detail: ErrInvalidID.Error(),
			},
		},
		{
			name: "given controller error, expect 400",
			givenController: mockController{
				GivenSpinError: memory.ErrInvalidKey,
			},
			givenURL:       "/v1/tables/70ee9bba-87ac-4155-8ec7-f83c8663315e/spin",
			expectedStatus: http.StatusBadRequest,
			expectedBody: responses.Error{
				Status: http.StatusBadRequest,
				Detail: memory.ErrInvalidKey.Error(),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := New(test.givenController)

			rr := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodPut, test.givenURL, nil)

			router := new(mux.Router)
			router.HandleFunc("/v1/tables/{id}/spin", handler.Spin)
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

			if !cmp.Equal(res, test.expectedBody) {
				t.Fatal(cmp.Diff(res, test.expectedBody))
			}
		})
	}
}

func TestHandler_Settle_Success(t *testing.T) {
	tests := []struct {
		name            string
		givenController Controller
		givenURL        string
		expectedStatus  int
		expectedBody    api.TableResponse
	}{
		{
			name: "given controller success, expect 200",
			givenController: mockController{
				GivenSettleTable: domain.Table{
					ID: uuid.MustParse("00812e8f-7fca-49a9-b141-9a52a0d0a82e"),
					Bets: []domain.Bet{
						{
							Status: domain.Settled,
						},
					},
					IsClosed: true,
					Outcome: &domain.Outcome{
						Value:  16,
						Colour: domain.Red,
					},
				},
			},
			givenURL:       "/v1/tables/00812e8f-7fca-49a9-b141-9a52a0d0a82e/settle",
			expectedStatus: http.StatusOK,
			expectedBody: api.TableResponse{
				ID: uuid.MustParse("00812e8f-7fca-49a9-b141-9a52a0d0a82e"),
				Bets: []api.BetResponse{
					{
						Status: domain.Settled,
					},
				},
				IsClosed: true,
				Outcome: &api.Outcome{
					Position: 16,
					Colour:   domain.Red,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := New(test.givenController)

			rr := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodPut, test.givenURL, nil)

			router := new(mux.Router)
			router.HandleFunc("/v1/tables/{id}/settle", handler.Settle)
			router.ServeHTTP(rr, req)

			resp := rr.Result()

			if !cmp.Equal(resp.StatusCode, test.expectedStatus) {
				t.Fatal(cmp.Diff(resp.StatusCode, test.expectedStatus))
			}

			var res api.TableResponse
			err := json.NewDecoder(resp.Body).Decode(&res)
			if err != nil {
				t.Fatal(err)
			}

			if !cmp.Equal(res, test.expectedBody) {
				t.Fatal(cmp.Diff(res, test.expectedBody))
			}
		})
	}
}

func TestHandler_Settle_Fail(t *testing.T) {
	tests := []struct {
		name            string
		givenController Controller
		givenURL        string
		expectedStatus  int
		expectedBody    responses.Error
	}{
		{
			name:            "given invalid ID, expect 400",
			givenController: mockController{},
			givenURL:        "/v1/tables/test/settle",
			expectedStatus:  http.StatusBadRequest,
			expectedBody: responses.Error{
				Status: http.StatusBadRequest,
				Detail: ErrInvalidID.Error(),
			},
		},
		{
			name: "given controller error, expect 400",
			givenController: mockController{
				GivenSettleError: memory.ErrInvalidKey,
			},
			givenURL:       "/v1/tables/70ee9bba-87ac-4155-8ec7-f83c8663315e/settle",
			expectedStatus: http.StatusBadRequest,
			expectedBody: responses.Error{
				Status: http.StatusBadRequest,
				Detail: memory.ErrInvalidKey.Error(),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := New(test.givenController)

			rr := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodPut, test.givenURL, nil)

			router := new(mux.Router)
			router.HandleFunc("/v1/tables/{id}/settle", handler.Settle)
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

			if !cmp.Equal(res, test.expectedBody) {
				t.Fatal(cmp.Diff(res, test.expectedBody))
			}
		})
	}
}

func TestHandler_List_Success(t *testing.T) {
	tests := []struct {
		name            string
		givenController Controller
		givenURL        string
		expectedStatus  int
		expectedBody    []api.TableResponse
	}{
		{
			name: "given controller success, expect 200",
			givenController: mockController{
				GivenListTables: []domain.Table{
					{
						ID: uuid.MustParse("00812e8f-7fca-49a9-b141-9a52a0d0a82e"),
						Bets: []domain.Bet{
							{
								Status: domain.Settled,
							},
						},
						IsClosed: true,
						Outcome: &domain.Outcome{
							Value:  16,
							Colour: domain.Red,
						},
					},
					{
						ID: uuid.MustParse("161ffdc3-564a-42b1-8340-0cf18b3cbfef"),
						Bets: []domain.Bet{
							{
								Status: domain.Unsettled,
							},
						},
						IsClosed: false,
						Outcome:  nil,
					},
				},
			},
			givenURL:       "/v1/tables",
			expectedStatus: http.StatusOK,
			expectedBody: []api.TableResponse{
				{
					ID: uuid.MustParse("00812e8f-7fca-49a9-b141-9a52a0d0a82e"),
					Bets: []api.BetResponse{
						{
							Status: domain.Settled,
						},
					},
					IsClosed: true,
					Outcome: &api.Outcome{
						Position: 16,
						Colour:   domain.Red,
					},
				},
				{
					ID: uuid.MustParse("161ffdc3-564a-42b1-8340-0cf18b3cbfef"),
					Bets: []api.BetResponse{
						{
							Status: domain.Unsettled,
						},
					},
					IsClosed: false,
					Outcome:  nil,
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
			router.HandleFunc("/v1/tables", handler.List)
			router.ServeHTTP(rr, req)

			resp := rr.Result()

			if !cmp.Equal(resp.StatusCode, test.expectedStatus) {
				t.Fatal(cmp.Diff(resp.StatusCode, test.expectedStatus))
			}

			var res []api.TableResponse
			err := json.NewDecoder(resp.Body).Decode(&res)
			if err != nil {
				t.Fatal(err)
			}

			if !cmp.Equal(res, test.expectedBody) {
				t.Fatal(cmp.Diff(res, test.expectedBody))
			}
		})
	}
}

func TestHandler_List_Fail(t *testing.T) {
	tests := []struct {
		name            string
		givenController Controller
		givenURL        string
		expectedStatus  int
		expectedBody    responses.Error
	}{
		{
			name: "given controller error, expect 400",
			givenController: mockController{
				GivenListError: memory.ErrInvalidKey,
			},
			givenURL:       "/v1/tables",
			expectedStatus: http.StatusBadRequest,
			expectedBody: responses.Error{
				Status: http.StatusBadRequest,
				Detail: memory.ErrInvalidKey.Error(),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := New(test.givenController)

			rr := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodGet, test.givenURL, nil)

			router := new(mux.Router)
			router.HandleFunc("/v1/tables", handler.List)
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

			if !cmp.Equal(res, test.expectedBody) {
				t.Fatal(cmp.Diff(res, test.expectedBody))
			}
		})
	}
}

type mockController struct {
	GivenGetTable    domain.Table
	GivenGetError    error
	GivenListTables  []domain.Table
	GivenListError   error
	GivenCreateTable domain.Table
	GivenCreateError error
	GivenSpinTable   domain.Table
	GivenSpinError   error
	GivenSettleTable domain.Table
	GivenSettleError error
}

func (m mockController) Get(_ context.Context, _ uuid.UUID) (domain.Table, error) {
	return m.GivenGetTable, m.GivenGetError
}

func (m mockController) List(_ context.Context) ([]domain.Table, error) {
	return m.GivenListTables, m.GivenListError
}

func (m mockController) Create(_ context.Context) (domain.Table, error) {
	return m.GivenCreateTable, m.GivenCreateError
}

func (m mockController) Spin(_ context.Context, _ uuid.UUID) (domain.Table, error) {
	return m.GivenSpinTable, m.GivenSpinError
}

func (m mockController) Settle(_ context.Context, _ uuid.UUID) (domain.Table, error) {
	return m.GivenSettleTable, m.GivenSettleError
}
