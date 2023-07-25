package go_indigo

import "fmt"

const (
	rows       = 9
	minColumns = 5
	maxColumns = 9
)

type board struct {
	tiles [][]*tile // 0,0 is upper left most tile
	gems  []*Gem
}

func newBoard() *board {
	var b = make([][]*tile, rows)
	columns := minColumns
	for i := 0; i < rows; i++ {
		b[i] = make([]*tile, columns)
		if columns < maxColumns {
			columns++
		} else {
			columns--
		}
	}

	// treasure tile data of paths:(row, col)
	treasureTiles := map[string][]int{
		C + E + D: {0, 0},
		D + F + E: {0, 4},
		A + E + F: {4, 8},
		B + F + A: {8, 8},
		A + C + B: {8, 4},
		B + D + C: {4, 0},
		Special:   {4, 4},
	}
	for edges, location := range treasureTiles {
		b[location[0]][location[1]] = newTreasureTile(edges)
	}

	// create gems on treasure tiles
	gems := []*Gem{
		newGem(Amber, D, 0, 0),
		newGem(Amber, E, 0, 4),
		newGem(Amber, F, 4, 8),
		newGem(Amber, A, 8, 4),
		newGem(Amber, B, 0, 0),
		newGem(Amber, C, 4, 0),
		newGem(Emerald, Special, 4, 4),
		newGem(Emerald, Special, 4, 4),
		newGem(Emerald, Special, 4, 4),
		newGem(Emerald, Special, 4, 4),
		newGem(Emerald, Special, 4, 4),
		newGem(Sapphire, Special, 4, 4),
	}

	return &board{
		tiles: b,
		gems:  gems,
	}
}

func (b *board) Place(tile *tile, row, col int) error {
	if row < 0 || col < 0 || row >= rows || col >= len(b.tiles[rows]) {
		return fmt.Errorf("index out of bounds")
	}
	if b.tiles[row][col] != nil {
		return fmt.Errorf("tile already exists at (%d, %d)", row, col)
	}
	b.tiles[row][col] = tile
	return nil
}

func (b *board) getTileCount() int {
	counter := 0
	for _, row := range b.tiles {
		for _, tile := range row {
			if tile != nil {
				counter++
			}
		}
	}
	return counter
}
