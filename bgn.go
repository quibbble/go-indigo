package go_indigo

import (
	"fmt"
	"strconv"

	bg "github.com/quibbble/go-boardgame"
	"github.com/quibbble/go-boardgame/pkg/bgerr"
)

var (
	actionToNotation = map[string]string{ActionPlaceTile: "p", bg.ActionSetWinners: "w"}
	notationToAction = reverseMap(actionToNotation)
)

func (p *PlaceTileActionDetails) encodeBGN() []string {
	return []string{strconv.Itoa(p.Row), strconv.Itoa(p.Column), p.Tile}
}

func decodePlaceTileActionDetailsBGN(notation []string) (*PlaceTileActionDetails, error) {
	if len(notation) != 3 {
		return nil, errDecoding(fmt.Errorf("invalid place tile notation"))
	}
	row, err := strconv.Atoi(notation[0])
	if err != nil {
		return nil, errDecoding(err)
	}
	column, err := strconv.Atoi(notation[1])
	if err != nil {
		return nil, errDecoding(err)
	}
	tile := notation[2]
	return &PlaceTileActionDetails{
		Row:    row,
		Column: column,
		Tile:   tile,
	}, nil
}

func errDecoding(err error) error {
	return &bgerr.Error{
		Err:    err,
		Status: bgerr.StatusBGNDecodingFailure,
	}
}
