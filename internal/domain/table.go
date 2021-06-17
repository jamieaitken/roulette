package domain

import "github.com/google/uuid"

// Table represents a single play of roulette.
type Table struct {
	ID       uuid.UUID
	Bets     []Bet
	IsClosed bool
	Outcome  *Outcome
}

// Outcome is the result of roulette wheel, it has a value and a Colour.
type Outcome struct {
	Value  int
	Colour Colour
}

// Colour represents the colour of the Outcome.
type Colour string

// String allows Colour to have a string representation.
func (c Colour) String() string {
	return string(c)
}

// Available options for Colour.
var (
	Red   Colour = "red"
	Black Colour = "black"
	Green Colour = "green"
)

// NumbersToColours allows the service to know what colour should be assigned a given number.
var NumbersToColours = map[int]Colour{
	0:  Green,
	32: Red,
	15: Black,
	19: Red,
	4:  Black,
	21: Red,
	2:  Black,
	25: Red,
	17: Black,
	34: Red,
	6:  Black,
	27: Red,
	13: Black,
	36: Red,
	11: Black,
	30: Red,
	8:  Black,
	23: Red,
	10: Black,
	5:  Red,
	24: Black,
	16: Red,
	33: Black,
	1:  Red,
	20: Black,
	14: Red,
	31: Black,
	9:  Red,
	22: Black,
	18: Red,
	29: Black,
	7:  Red,
	28: Black,
	12: Red,
	35: Black,
	3:  Red,
	26: Black,
}
