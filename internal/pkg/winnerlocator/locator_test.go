package winnerlocator

import (
	"betting/internal/domain"
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/google/uuid"
)

func TestLocator_Locate(t *testing.T) {
	tests := []struct {
		name          string
		givenTable    domain.Table
		expectedTable domain.Table
	}{
		{
			name: "given a table that has a winner, mark the bet as won",
			givenTable: domain.Table{
				ID: uuid.MustParse("160998da-2d89-4f06-a690-fd189213958d"),
				Bets: []domain.Bet{
					{
						ID:             uuid.MustParse("e49779f6-3507-4063-bed8-18d50174868d"),
						SelectedSpaces: []int{14, 16},
						Table:          uuid.MustParse("160998da-2d89-4f06-a690-fd189213958d"),
					},
					{
						ID:             uuid.MustParse("e49779f6-3507-4063-bed8-18d50174868d"),
						SelectedSpaces: []int{14},
						Table:          uuid.MustParse("160998da-2d89-4f06-a690-fd189213958d"),
					},
				},
				IsClosed: false,
				Outcome: &domain.Outcome{
					Value:  16,
					Colour: domain.Red,
				},
			},
			expectedTable: domain.Table{
				ID: uuid.MustParse("160998da-2d89-4f06-a690-fd189213958d"),
				Bets: []domain.Bet{
					{
						ID:             uuid.MustParse("e49779f6-3507-4063-bed8-18d50174868d"),
						SelectedSpaces: []int{14, 16},
						Win:            true,
						Table:          uuid.MustParse("160998da-2d89-4f06-a690-fd189213958d"),
					},
					{
						ID:             uuid.MustParse("e49779f6-3507-4063-bed8-18d50174868d"),
						SelectedSpaces: []int{14},
						Table:          uuid.MustParse("160998da-2d89-4f06-a690-fd189213958d"),
					},
				},
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
			l := New()

			actual := l.Locate(context.Background(), test.givenTable)

			if !cmp.Equal(actual, test.expectedTable) {
				t.Fatal(cmp.Diff(actual, test.expectedTable))
			}
		})
	}
}
