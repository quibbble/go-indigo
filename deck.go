package go_indigo

import (
	"errors"
	"math/rand"
)

type deck struct {
	tiles  []*tile
	random *rand.Rand
}

func newDeck(random *rand.Rand) *deck {
	d := make([]*tile, 0)
	for _, edges := range UniqueTiles {
		for i := 0; i < NumTileCopies; i++ {
			t, _ := newTile(edges)
			d = append(d, t)
		}
	}
	result := &deck{
		tiles:  d,
		random: random,
	}
	result.Shuffle()
	return result
}

func (d *deck) Remove(tile *tile) error {
	for idx, t := range d.tiles {
		if tile.equals(t) {
			d.tiles = append(d.tiles[:idx], d.tiles[idx+1:]...)
			return nil
		}
	}
	return errors.New("tile not found")
}

func (d *deck) Add(tiles ...*tile) {
	d.tiles = append(d.tiles, tiles...)
	d.Shuffle()
}

func (d *deck) Draw() (*tile, error) {
	size := len(d.tiles)
	if size <= 0 {
		return nil, errors.New("deck is empty so cannot draw")
	}
	tile := d.tiles[size-1]
	d.tiles = d.tiles[:size-1]
	return tile, nil
}

func (d *deck) Shuffle() {
	for i := 0; i < len(d.tiles); i++ {
		r := d.random.Intn(len(d.tiles))
		if i != r {
			d.tiles[r], d.tiles[i] = d.tiles[i], d.tiles[r]
		}
	}
}
