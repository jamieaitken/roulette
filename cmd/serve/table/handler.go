package table

import (
	"betting/api"
	"betting/internal/domain"
	"betting/internal/pkg/responses"
	"context"
	"errors"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"

	"github.com/google/uuid"
)

// Errors returned by the Handler.
var (
	ErrFailedToCreateTable = errors.New("failed to create table")
	ErrNoIDPresent         = errors.New("failed locate table")
	ErrInvalidID           = errors.New("id given is not a uuid")
)

// Controller provides business logic capable of both read and writes.
type Controller interface {
	ControllerReader
	ControllerWriter
}

// ControllerWriter provides business logic capable of writes.
type ControllerWriter interface {
	Create(ctx context.Context) (domain.Table, error)
	Spin(ctx context.Context, id uuid.UUID) (domain.Table, error)
	Settle(ctx context.Context, id uuid.UUID) (domain.Table, error)
}

// ControllerReader provides business logic capable of reads.
type ControllerReader interface {
	Get(ctx context.Context, id uuid.UUID) (domain.Table, error)
	List(ctx context.Context) ([]domain.Table, error)
}

// Handler handles requests relating to tables.
type Handler struct {
	Controller Controller
}

// New instantiates a Handler.
func New(controller Controller) Handler {
	return Handler{
		Controller: controller,
	}
}

// Create creates a table in storage where which bets can be placed on it.
func (h Handler) Create(w http.ResponseWriter, r *http.Request) {
	table, err := h.Controller.Create(r.Context())
	if err != nil {
		log.Errorf("failed to create table: %v", err)

		responses.NewJSON(w).Fail(http.StatusBadRequest, ErrFailedToCreateTable)
		return
	}

	log.Infof("created table: %v", table.ID)

	responses.NewJSON(w).Success(http.StatusCreated, api.AdaptTableFromDomain(table))
}

// Get retrieves the table for the given ID.
func (h Handler) Get(w http.ResponseWriter, r *http.Request) {
	path := mux.Vars(r)

	id, ok := path["id"]
	if !ok {
		log.Errorf("invalid id: %v", ErrNoIDPresent)

		responses.NewJSON(w).Fail(http.StatusBadRequest, ErrNoIDPresent)
		return
	}

	tableID, err := uuid.Parse(id)
	if err != nil {
		log.Errorf("invalid id: %v, %v", id, ErrInvalidID)

		responses.NewJSON(w).Fail(http.StatusBadRequest, ErrInvalidID)
		return
	}

	t, err := h.Controller.Get(r.Context(), tableID)
	if err != nil {
		log.Errorf("failed to locate table: %v, %v", tableID, err)

		responses.NewJSON(w).Fail(http.StatusBadRequest, err)
		return
	}

	resBody := api.AdaptTableFromDomain(t)

	responses.NewJSON(w).Success(http.StatusOK, resBody)
}

// Spin locks the table and returns the outcome.
func (h Handler) Spin(w http.ResponseWriter, r *http.Request) {
	path := mux.Vars(r)

	id, ok := path["id"]
	if !ok {
		log.Errorf("invalid id: %v", ErrNoIDPresent)

		responses.NewJSON(w).Fail(http.StatusBadRequest, ErrNoIDPresent)
		return
	}

	tableID, err := uuid.Parse(id)
	if err != nil {
		log.Errorf("invalid id: %v, %v", id, ErrInvalidID)

		responses.NewJSON(w).Fail(http.StatusBadRequest, ErrInvalidID)
		return
	}

	t, err := h.Controller.Spin(r.Context(), tableID)
	if err != nil {
		log.Errorf("failed to spin table: %v, %v", tableID, err)

		responses.NewJSON(w).Fail(http.StatusBadRequest, err)
		return
	}

	resBody := api.AdaptTableFromDomain(t)

	log.Infof("table has been spun: %v", t.ID)

	responses.NewJSON(w).Success(http.StatusOK, resBody)
}

// Settle moves all bets to settled and finds any winners.
func (h Handler) Settle(w http.ResponseWriter, r *http.Request) {
	path := mux.Vars(r)

	id, ok := path["id"]
	if !ok {
		log.Errorf("invalid id: %v", ErrNoIDPresent)

		responses.NewJSON(w).Fail(http.StatusBadRequest, ErrNoIDPresent)
		return
	}

	tableID, err := uuid.Parse(id)
	if err != nil {
		log.Errorf("invalid id: %v, %v", id, ErrInvalidID)

		responses.NewJSON(w).Fail(http.StatusBadRequest, ErrInvalidID)
		return
	}

	t, err := h.Controller.Settle(r.Context(), tableID)
	if err != nil {
		log.Errorf("failed to settle table: %v, %v", tableID, err)

		responses.NewJSON(w).Fail(http.StatusBadRequest, err)
		return
	}

	resBody := api.AdaptTableFromDomain(t)

	log.Infof("settled table: %v", t.ID)

	responses.NewJSON(w).Success(http.StatusOK, resBody)
}

// List returns all tables from storage.
func (h Handler) List(w http.ResponseWriter, r *http.Request) {
	t, err := h.Controller.List(r.Context())
	if err != nil {
		log.Errorf("failed to locate bets: %v", err)

		responses.NewJSON(w).Fail(http.StatusBadRequest, err)
		return
	}

	resBody := api.AdaptTablesFromDomain(t)

	log.Infof("located tables")

	responses.NewJSON(w).Success(http.StatusOK, resBody)
}
