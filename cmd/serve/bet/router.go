package bet

import (
	"betting/internal/bet"
	"betting/internal/table"
	"net/http"

	"github.com/gorilla/mux"
)

func Load(r *mux.Router, tableStorage table.StorageProvider, betStorage bet.StorageProvider) *mux.Router {
	controller := bet.NewController(
		bet.NewRepository(betStorage),
		table.NewRepository(tableStorage),
	)

	handler := New(controller)

	r.HandleFunc("/v1/tables/{id}/bet", handler.Create).Methods(http.MethodPost)
	r.HandleFunc("/v1/bets/{id}", handler.Get).Methods(http.MethodGet)

	return r
}
