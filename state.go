package go_indigo

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/quibbble/go-boardgame/pkg/bgerr"
)

type state struct {
	turn                  string
	teams                 []string
	winners               []string
	board                 *board
	deck                  *deck
	hands                 map[string]*hand
	variant               string
	points                map[string]int
	round, roundsUntilEnd int
}

func newState(teams []string, random *rand.Rand, variant string, roundsUntilEnd int) (*state, error) {
	if random == nil {
		return nil, fmt.Errorf("random seed is null")
	}

	hands := make(map[string]*hand)
	points := make(map[string]int)
	deck := newDeck(random)

	switch variant {
	case VariantClassic:
		for _, team := range teams {
			hand := newHand()
			for i := 0; i < 1; i++ {
				tile, err := deck.Draw()
				if err != nil {
					return nil, err
				}
				hand.Add(tile)
			}
			points[team] = 0
			hands[team] = hand
		}
	case VariantLargeHands:
		for _, team := range teams {
			hand := newHand()
			for i := 0; i < 3; i++ {
				tile, err := deck.Draw()
				if err != nil {
					return nil, err
				}
				hand.Add(tile)
			}
			points[team] = 0
			hands[team] = hand
		}
	}

	return &state{
		turn:           teams[0],
		teams:          teams,
		winners:        make([]string, 0),
		board:          newBoard(),
		deck:           deck,
		hands:          hands,
		variant:        variant,
		points:         points,
		round:          0,
		roundsUntilEnd: roundsUntilEnd,
	}, nil
}

func (s *state) rotateTileClockwise(team, tile string) error

func (s *state) placeTile(team, tile string, row, col int) error

func (s *state) setWinners(winners []string) error {
	for _, winner := range winners {
		if !contains(s.teams, winner) {
			return &bgerr.Error{
				Err:    fmt.Errorf("winner not in teams"),
				Status: bgerr.StatusInvalidActionDetails,
			}
		}
	}
	s.winners = winners
	return nil
}

func (s *state) message() string {
	message := fmt.Sprintf("%s must place a tile", s.turn)
	if len(s.winners) > 0 {
		message = fmt.Sprintf("%s tie", strings.Join(s.winners, ", "))
		if len(s.winners) == 1 {
			message = fmt.Sprintf("%s wins", s.winners[0])
		}
	}
	return message
}
