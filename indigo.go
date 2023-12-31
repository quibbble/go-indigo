package go_indigo

import (
	"fmt"
	"strings"

	"github.com/mitchellh/mapstructure"
	bg "github.com/quibbble/go-boardgame"
	"github.com/quibbble/go-boardgame/pkg/bgerr"
	"github.com/quibbble/go-boardgame/pkg/bgn"
)

const (
	minTeams = 2
	maxTeams = 4
)

type Indigo struct {
	state   *state
	actions []*bg.BoardGameAction
	options *IndigoMoreOptions
}

func NewIndigo(options *bg.BoardGameOptions) (*Indigo, error) {
	if len(options.Teams) < minTeams {
		return nil, &bgerr.Error{
			Err:    fmt.Errorf("at least %d teams required to create a game of %s", minTeams, key),
			Status: bgerr.StatusTooFewTeams,
		}
	} else if len(options.Teams) > maxTeams {
		return nil, &bgerr.Error{
			Err:    fmt.Errorf("at most %d teams allowed to create a game of %s", maxTeams, key),
			Status: bgerr.StatusTooManyTeams,
		}
	} else if duplicates(options.Teams) {
		return nil, &bgerr.Error{
			Err:    fmt.Errorf("duplicate teams found"),
			Status: bgerr.StatusInvalidOption,
		}
	}
	var details IndigoMoreOptions
	if err := mapstructure.Decode(options.MoreOptions, &details); err != nil {
		return nil, &bgerr.Error{
			Err:    err,
			Status: bgerr.StatusInvalidOption,
		}
	}
	if details.Variant == "" {
		details.Variant = VariantClassic
	} else if !contains(variants, details.Variant) {
		return nil, &bgerr.Error{
			Err:    fmt.Errorf("invalid Indigo variant"),
			Status: bgerr.StatusInvalidOption,
		}
	}
	if details.RoundsUntilEnd == 0 {
		details.RoundsUntilEnd = 999
	}
	state, err := newState(options.Teams, details.Seed, details.Variant, details.RoundsUntilEnd)
	if err != nil {
		return nil, &bgerr.Error{
			Err:    err,
			Status: bgerr.StatusInvalidOption,
		}
	}
	return &Indigo{
		state:   state,
		actions: make([]*bg.BoardGameAction, 0),
		options: &details,
	}, nil
}

func (i *Indigo) Do(action *bg.BoardGameAction) error {
	if len(i.state.winners) > 0 {
		return &bgerr.Error{
			Err:    fmt.Errorf("game already over"),
			Status: bgerr.StatusGameOver,
		}
	}
	switch action.ActionType {
	case ActionRotateTileClockwise:
		var details RotateTileActionDetails
		if err := mapstructure.Decode(action.MoreDetails, &details); err != nil {
			return &bgerr.Error{
				Err:    err,
				Status: bgerr.StatusInvalidActionDetails,
			}
		}
		if err := i.state.rotateTileClockwise(action.Team, details.Tile); err != nil {
			return err
		}
	case ActionPlaceTile:
		var details PlaceTileActionDetails
		if err := mapstructure.Decode(action.MoreDetails, &details); err != nil {
			return &bgerr.Error{
				Err:    err,
				Status: bgerr.StatusInvalidActionDetails,
			}
		}
		if err := i.state.placeTile(action.Team, details.Tile, details.Row, details.Column); err != nil {
			return err
		}
		i.actions = append(i.actions, action)
	case bg.ActionSetWinners:
		var details bg.SetWinnersActionDetails
		if err := mapstructure.Decode(action.MoreDetails, &details); err != nil {
			return &bgerr.Error{
				Err:    err,
				Status: bgerr.StatusInvalidActionDetails,
			}
		}
		if err := i.state.setWinners(details.Winners); err != nil {
			return err
		}
		i.actions = append(i.actions, action)
	default:
		return &bgerr.Error{
			Err:    fmt.Errorf("cannot process action type %s", action.ActionType),
			Status: bgerr.StatusUnknownActionType,
		}
	}
	return nil
}

func (i *Indigo) GetSnapshot(team ...string) (*bg.BoardGameSnapshot, error) {
	if len(team) > 1 {
		return nil, &bgerr.Error{
			Err:    fmt.Errorf("get snapshot requires zero or one team"),
			Status: bgerr.StatusTooManyTeams,
		}
	}

	hands := make(map[string][]tile)
	for t, hand := range i.state.hands {
		if len(team) == 0 || (t == team[0]) {
			hands[t] = hand.GetItems()
		}
	}

	return &bg.BoardGameSnapshot{
		Turn:    i.state.turn,
		Teams:   i.state.teams,
		Winners: i.state.winners,
		MoreData: IndigoSnapshotData{
			Board:          i.state.board,
			Hands:          hands,
			Points:         i.state.points,
			Round:          i.state.round,
			RoundsUntilEnd: i.state.roundsUntilEnd,
			Variant:        i.state.variant,
		},
		Targets: i.state.targets(),
		Actions: i.actions,
		Message: i.state.message(),
	}, nil
}

func (i *Indigo) GetBGN() *bgn.Game {
	tags := map[string]string{
		"Game":  key,
		"Teams": strings.Join(i.state.teams, ", "),
		"Seed":  fmt.Sprintf("%d", i.options.Seed),
	}
	actions := make([]bgn.Action, 0)
	for _, action := range i.actions {
		bgnAction := bgn.Action{
			TeamIndex: indexOf(i.state.teams, action.Team),
			ActionKey: rune(actionToNotation[action.ActionType][0]),
		}
		switch action.ActionType {
		case ActionPlaceTile:
			var details PlaceTileActionDetails
			_ = mapstructure.Decode(action.MoreDetails, &details)
			bgnAction.Details = details.encodeBGN()
		case bg.ActionSetWinners:
			var details bg.SetWinnersActionDetails
			_ = mapstructure.Decode(action.MoreDetails, &details)
			bgnAction.Details, _ = details.EncodeBGN(i.state.teams)
		}
		actions = append(actions, bgnAction)
	}
	return &bgn.Game{
		Tags:    tags,
		Actions: actions,
	}
}
