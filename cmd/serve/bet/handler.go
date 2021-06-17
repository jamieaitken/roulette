package bet

import (
	"betting/api"
	"betting/internal/domain"
	"betting/internal/pkg/responses"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/google/uuid"

	"github.com/gorilla/mux"
)

// Errors returned by the Handler.
var (
	ErrNoIDPresent = errors.New("failed to associate bet with table")
	ErrInvalidID   = errors.New("id given is not a uuid")
)

// Controller provides business logic capable of both read and writes.
type Controller interface {
	ControllerWriter
	ControllerReader
}

// ControllerWriter provides business logic capable of writes.
type ControllerWriter interface {
	Create(ctx context.Context, bet domain.Bet) (domain.Bet, error)
}

// ControllerReader provides business logic capable of reads.
type ControllerReader interface {
	Get(ctx context.Context, id uuid.UUID) (domain.Bet, error)
}

// Handler handles requests relating to bets.
type Handler struct {
	Controller Controller
}

// New instantiates a Handler.
func New(controller Controller) Handler {
	return Handler{
		Controller: controller,
	}
}

// Create inserts the given bet into storage.
func (h Handler) Create(w http.ResponseWriter, r *http.Request) {
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

	var betRequest api.BetRequest
	decoder := json.NewDecoder(r.Body)

	err = decoder.Decode(&betRequest)
	if err != nil {
		log.Errorf("could not decode request body: %v", err)
		responses.NewJSON(w).Fail(http.StatusBadRequest, err)
		return
	}

	domainBet := api.AdaptBetToDomain(betRequest, tableID)

	bet, err := h.Controller.Create(r.Context(), domainBet)
	if err != nil {
		log.Errorf("failed to create bet: %v, %v", tableID, err)
		responses.NewJSON(w).Fail(http.StatusBadRequest, err)
		return
	}

	resBody := api.AdaptBetFromDomain(bet)

	log.Infof("created bet: %v", bet.ID)

	responses.NewJSON(w).Success(http.StatusCreated, resBody)
}

// Get retrieves the bet for the given ID.
func (h Handler) Get(w http.ResponseWriter, r *http.Request) {
	path := mux.Vars(r)

	id, ok := path["id"]
	if !ok {
		log.Errorf("invalid id: %v", ErrNoIDPresent)
		responses.NewJSON(w).Fail(http.StatusBadRequest, ErrNoIDPresent)
		return
	}

	betID, err := uuid.Parse(id)
	if err != nil {
		log.Errorf("invalid id: %v, %v", id, ErrInvalidID)

		responses.NewJSON(w).Fail(http.StatusBadRequest, ErrInvalidID)
		return
	}

	bet, err := h.Controller.Get(r.Context(), betID)
	if err != nil {
		log.Errorf("failed to locate bet: %v, %v", betID, err)
		responses.NewJSON(w).Fail(http.StatusBadRequest, err)
		return
	}

	resBody := api.AdaptBetFromDomain(bet)

	log.Infof("located bet: %v", bet.ID)

	responses.NewJSON(w).Success(http.StatusOK, resBody)
}
