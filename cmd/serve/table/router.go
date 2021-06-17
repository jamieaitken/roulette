package table

import (
	"betting/internal/bet"
	"betting/internal/pkg/ballplacer"
	"betting/internal/pkg/winnerlocator"
	"betting/internal/table"
	"net/http"

	"github.com/gorilla/mux"
)

func Load(r *mux.Router, tableStorage table.StorageProvider, betStorage bet.StorageProvider) *mux.Router {
	controller := table.NewController(table.ControllerParams{
		RepositoryProvider:    table.NewRepository(tableStorage),
		BallPlacer:            ballplacer.New(),
		WinnerLocator:         winnerlocator.New(),
		BetRepositoryProvider: bet.NewRepository(betStorage),
	})

	handler := New(controller)

	r.HandleFunc("/v1/tables", handler.Create).Methods(http.MethodPost)
	r.HandleFunc("/v1/tables", handler.List).Methods(http.MethodGet)
	r.HandleFunc("/v1/tables/{id}", handler.Get).Methods(http.MethodGet)
	r.HandleFunc("/v1/tables/{id}/spin", handler.Spin).Methods(http.MethodPut)
	r.HandleFunc("/v1/tables/{id}/settle", handler.Settle).Methods(http.MethodPut)

	return r
}
