package go_indigo

import (
	"fmt"
	"strings"

	bg "github.com/quibbble/go-boardgame"
	"github.com/quibbble/go-boardgame/pkg/bgerr"
	cl "github.com/quibbble/go-boardgame/pkg/collection"
)

type state struct {
	turn                  string
	teams                 []string
	winners               []string
	board                 *board
	deck                  *cl.Collection[tile]
	hands                 map[string]*cl.Collection[tile]
	variant               string
	points                map[string]int
	gemsCount             map[string]int
	round, roundsUntilEnd int
}

func newState(teams []string, random int64, variant string, roundsUntilEnd int) (*state, error) {

	hands := make(map[string]*cl.Collection[tile])
	points := make(map[string]int)
	gemsCount := make(map[string]int)
	deck := cl.NewCollection[tile](random)
	for idx, numCopies := range numCopiesByUniquePathsIndex {
		for i := 0; i < numCopies; i++ {
			deck.Add(tile{Paths: uniquePaths[idx]})
		}
	}
	deck.Shuffle()

	switch variant {
	case VariantClassic:
		for _, team := range teams {
			hand := cl.NewCollection[tile](0)
			for i := 0; i < 1; i++ {
				tile, err := deck.Draw()
				if err != nil {
					return nil, err
				}
				hand.Add(*tile)
			}
			points[team] = 0
			gemsCount[team] = 0
			hands[team] = hand
		}
	case VariantLargeHands:
		for _, team := range teams {
			hand := cl.NewCollection[tile](0)
			for i := 0; i < 2; i++ {
				tile, err := deck.Draw()
				if err != nil {
					return nil, err
				}
				hand.Add(*tile)
			}
			points[team] = 0
			gemsCount[team] = 0
			hands[team] = hand
		}
	}

	return &state{
		turn:           teams[0],
		teams:          teams,
		winners:        make([]string, 0),
		board:          newBoard(teams),
		deck:           deck,
		hands:          hands,
		variant:        variant,
		points:         points,
		gemsCount:      gemsCount,
		round:          0,
		roundsUntilEnd: roundsUntilEnd,
	}, nil
}

func (s *state) rotateTileClockwise(team, paths string) error {
	if !contains(s.teams, team) {
		return &bgerr.Error{
			Err:    fmt.Errorf("%s not a valid team", team),
			Status: bgerr.StatusUnknownTeam,
		}
	}
	t, err := newTile(paths)
	if err != nil {
		return &bgerr.Error{
			Err:    err,
			Status: bgerr.StatusInvalidActionDetails,
		}
	}
	idx := s.hands[team].IndexOf(*t, func(a, b tile) bool { return a.equals(&b) })
	if idx < 0 {
		return &bgerr.Error{
			Err:    fmt.Errorf("%s's hand does not contain %s", team, paths),
			Status: bgerr.StatusInvalidActionDetails,
		}
	}
	tile, _ := s.hands[team].GetItem(idx)
	tile.RotateClockwise()
	return nil
}

func (s *state) placeTile(team, paths string, row, col int) error {
	if team != s.turn {
		return &bgerr.Error{
			Err:    fmt.Errorf("%s cannot play on %s turn", team, s.turn),
			Status: bgerr.StatusWrongTurn,
		}
	}
	t, err := newTile(paths)
	if err != nil {
		return &bgerr.Error{
			Err:    err,
			Status: bgerr.StatusInvalidActionDetails,
		}
	}

	// place tile and remove it from your hand
	tileIdx := s.hands[team].IndexOf(*t, func(a, b tile) bool { return a.equals(&b) })
	if tileIdx < 0 {
		return &bgerr.Error{
			Err:    fmt.Errorf("%s's hand does not contain %s", team, paths),
			Status: bgerr.StatusInvalidAction,
		}
	}
	if err := s.board.place(t, row, col); err != nil {
		return &bgerr.Error{
			Err:    err,
			Status: bgerr.StatusInvalidActionDetails,
		}
	}
	_ = s.hands[team].Remove(tileIdx)

	// update gem locations
	movedGems, err := s.board.moveGems(row, col)
	if err != nil {
		return &bgerr.Error{
			Err:    err,
			Status: bgerr.StatusInvalidActionDetails,
		}
	}

	// update scores based on new gem locations
	for _, gem := range movedGems {
		if gem.gateway != nil {
			for _, team := range gem.gateway.Teams {
				s.points[team] += colorToPoints[gem.Color]
				s.gemsCount[team] += 1
			}
		}
	}

	// draw tile and add to hand if there tiles left in the deck
	if t, err = s.deck.Draw(); err == nil {
		s.hands[team].Add(*t)
	}

	// change turn
	for idx, team := range s.teams {
		if team == s.turn {
			s.turn = s.teams[(idx+1)%len(s.teams)]
			break
		}
	}

	// inc round counter
	if s.turn == s.teams[0] {
		s.round++
	}

	// check if the game is over and set winners if so
	if s.round >= s.roundsUntilEnd || s.board.gemsInPlay() <= 0 {
		winners := make([]string, 0)
		maxPoints := 0
		for team, points := range s.points {
			if points == maxPoints {
				winners = append(winners, team)
			} else if points > maxPoints {
				winners = []string{team}
				maxPoints = points
			}
		}
		// if tied the player with most points AND gems wins
		if len(winners) > 1 {
			possibleWinners := winners
			winners = make([]string, 0)
			maxGemCount := 0
			for _, team := range possibleWinners {
				gemCount := s.gemsCount[team]
				if gemCount == maxGemCount {
					winners = append(winners, team)
				} else if gemCount > maxGemCount {
					winners = []string{team}
					maxGemCount = gemCount
				}
			}
		}
		s.winners = winners
	}

	return nil
}

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

func (s *state) targets(team ...string) []*bg.BoardGameAction {
	targets := make([]*bg.BoardGameAction, 0)
	if len(s.winners) > 0 {
		return targets
	}
	// rotate tile actions
	if len(team) == 0 || len(team) == 1 {
		for t, hand := range s.hands {
			for _, tile := range hand.GetItems() {
				targets = append(targets, &bg.BoardGameAction{
					Team:       t,
					ActionType: ActionRotateTileClockwise,
					MoreDetails: RotateTileActionDetails{
						Tile: tile.Paths,
					},
				})
			}
		}
	}
	// place tile actions
	if len(team) == 0 || (len(team) == 1 && team[0] == s.turn) {
		for r, row := range s.board.Tiles {
			for c, t := range row {
				if t == nil {
					for _, tile := range s.hands[s.turn].GetItems() {
						targets = append(targets, &bg.BoardGameAction{
							Team:       s.turn,
							ActionType: ActionPlaceTile,
							MoreDetails: PlaceTileActionDetails{
								Tile:   tile.Paths,
								Row:    r,
								Column: c,
							},
						})
					}
				}
			}
		}
	}
	return targets
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
