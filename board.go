package go_indigo

import (
	"fmt"
	"strings"
)

const (
	rows       = 9
	minColumns = 5
	maxColumns = 9
)

type board struct {
	Tiles    [][]*tile // 0,0 is upper left most tile
	Gateways []*gateway
	Gems     []*gem
}

func newBoard(teams []string) *board {
	var b = make([][]*tile, rows)
	columns := minColumns
	for i := 0; i < rows; i++ {
		b[i] = make([]*tile, columns)
		if i < 4 {
			columns++
		} else {
			columns--
		}
	}

	// place treasure tiles
	for edges, location := range initTreasureTiles {
		b[location[0]][location[1]] = newTreasureTile(edges)
	}

	// create gateways
	gateways := make([]*gateway, 0)
	for edges, teamsIdxs := range numTeamsToGatewayOwnership[len(teams)] {
		owners := make([]string, 0)
		for _, idx := range teamsIdxs {
			owners = append(owners, teams[idx])
		}
		gateways = append(gateways, newGateway(initGateways[edges], edges, owners...))
	}

	return &board{
		Tiles:    b,
		Gems:     initGems,
		Gateways: gateways,
	}
}

func (b *board) place(tile *tile, row, col int) error {
	if row < 0 || col < 0 || row >= rows || col >= len(b.Tiles[row]) {
		return fmt.Errorf("index out of bounds")
	}
	if b.Tiles[row][col] != nil {
		return fmt.Errorf("tile already exists at (%d, %d)", row, col)
	}
	b.Tiles[row][col] = tile
	return nil
}

func (b *board) moveGems(placedRow, placedCol int) ([]*gem, error) {
	moved := []*gem{}
	centerGemMoved := false

gemsOuter:
	for _, gem := range b.Gems {
		if gem.collided || gem.gateway != nil {
			continue
		}

		var (
			adjRow, adjCol int
			adjEdge        string
		)

		// case where tile placed adj to middle treasure tile and one gem must be moved
		if gem.Edge == Special && !centerGemMoved && placedRow >= 0 && placedCol >= 0 {
			edgeToRowCol := map[string][2]int{A: {-1, -1}, B: {-1, 0}, C: {0, 1}, D: {1, 1}, E: {1, 0}, F: {0, -1}}
			edgeToEdge := map[string]string{A: D, B: E, C: F, D: A, E: B, F: C}
			for edge, loc := range edgeToRowCol {
				if gem.Row+loc[0] == placedRow &&
					gem.Column+loc[1] == placedCol {
					adjRow = placedRow
					adjCol = placedCol
					adjEdge = edgeToEdge[edge]
					centerGemMoved = true
					break
				}
			}
			if !centerGemMoved {
				continue
			}
		}

		// base case where gem has a adj tile and must be moved
		if adjEdge == "" {
			adjRow, adjCol, adjEdge = b.getAdjacent(gem.Row, gem.Column, gem.Edge)
			if adjRow < 0 || adjRow >= len(b.Tiles) ||
				adjCol < 0 || adjCol >= len(b.Tiles[adjRow]) ||
				b.Tiles[adjRow][adjCol] == nil {
				continue
			}
		}

		// check for collision
		for _, g := range b.Gems {
			if g.Row == adjRow && g.Column == adjCol && g.Edge == adjEdge {
				gem.collided = true
				g.collided = true
				continue gemsOuter
			}
		}

		movedEdge, err := b.Tiles[adjRow][adjCol].GetDestination(adjEdge)
		if err != nil {
			return nil, err
		}

		gem.Row = adjRow
		gem.Column = adjCol
		gem.Edge = movedEdge

		moved = append(moved, gem)

		// check for gateway reached
	gatewayOuter:
		for _, gateway := range b.Gateways {
			for _, loc := range gateway.Locations {
				if loc[0] == gem.Row && loc[1] == gem.Column && strings.Contains(gateway.Edges, gem.Edge) {
					gem.gateway = gateway
					break gatewayOuter
				}
			}
		}
	}

	if len(moved) > 0 {
		// NOTE only gems moved the first iteration could be moved again so do not need to concat returned gems on future iterations
		_, err := b.moveGems(-1, -1)
		if err != nil {
			return nil, err
		}
	}

	return moved, nil
}

// getAdjacent returns the adjacent row, col, and edge
func (b *board) getAdjacent(row, col int, edge string) (adjRow, adjCol int, adjEdge string) {
	edgeToRowCol := map[string][2]int{A: {-1, -1}, B: {-1, 0}, C: {0, 1}, D: {1, 1}, E: {1, 0}, F: {0, -1}}
	edgeToEdge := map[string]string{A: D, B: E, C: F, D: A, E: B, F: C}
	return row + edgeToRowCol[edge][0], col + edgeToRowCol[edge][1], edgeToEdge[edge]
}

func (b *board) gemsInPlay() int {
	count := 0
	for _, gem := range b.Gems {
		if !gem.collided && gem.gateway == nil {
			count++
		}
	}
	return count
}
