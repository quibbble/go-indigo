package go_indigo

const (
	Amber    = "Amber"
	Emerald  = "Emerald"
	Sapphire = "Sapphire"
)

type Gem struct {
	Color          string
	Edge           string
	Row, Column    int
	GatewayReached bool
}

func newGem(color, edge string, row, column int) *Gem {
	return &Gem{
		Color:          color,
		Edge:           edge,
		Row:            row,
		Column:         column,
		GatewayReached: false,
	}
}
