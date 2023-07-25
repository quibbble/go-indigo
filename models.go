package go_indigo

// Action types
const (
	ActionRotateTileClockwise = "RotateTileClockwise" // NOTE - this is not tracked by BGN
	ActionPlaceTile           = "PlaceTile"
)

// Indigo Variants
const (
	VariantClassic    = "Classic"    // normal Indigo
	VariantLargeHands = "LargeHands" // players have a hand size of 3 instead of 1
)

var Variants = []string{VariantClassic, VariantLargeHands}

// IndigoMoreOptions are the additional options for creating a game of Indigo
type IndigoMoreOptions struct {
	Seed           int64
	Variant        string
	RoundsUntilEnd int // the number of rounds until the game ends - 0 means infinite
}

type RotateTileClockwiseActionDetails struct {
	Tile string
}

type PlaceTileActionDetails struct {
	Tile        string
	Row, Column int
}

// IndigoSnapshotData is the game data unique to Indigo
type IndigoSnapshotData struct {
	Board [][]*string
}

// the number of times a unique tile is added to a new deck
const NumTileCopies = 6

// list of all the tiles that can be played - 9 total
var UniqueTiles = []string{}
