package ballplacer

import (
	"betting/internal/domain"
	"context"
	"crypto/rand"
	"log"
	"math/big"
	"strings"
	"time"
)

// Placer generates the result of the roulette table.
type Placer struct {
}

// New instantiates a Placer.
func New() Placer {
	return Placer{}
}

// GetPosition returns the domain.Outcome which is a number between 0 and 36 along with the assigned colour.
func (p Placer) GetPosition(_ context.Context) domain.Outcome {
	rnd, err := rand.Int(strings.NewReader(time.Now().UTC().String()), big.NewInt(36))
	if err != nil {
		log.Fatal("random number failed to generate")
		return domain.Outcome{}
	}

	position := int(rnd.Int64())

	colour := domain.NumbersToColours[position]

	return domain.Outcome{
		Value:  position,
		Colour: colour,
	}
}
