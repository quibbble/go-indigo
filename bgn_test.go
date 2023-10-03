package go_indigo

import (
	"fmt"
	"testing"

	"github.com/quibbble/go-boardgame/pkg/bgn"
)

func Test_(t *testing.T) {
	raw := `
	[Game "Indigo"]
	[Teams "red, blue"]
	[Seed "1696338136223"]

	0p&8.3.ADBFCE 1p&8.1.FEABCD 0p&8.2.BACDEF 1p&7.0.CFDBEA 0p&7.1.CFDBEA`

	builder := Builder{}
	game, err := bgn.Parse(raw)
	if err != nil {
		fmt.Println("Parse: ", err)
		t.FailNow()
	}
	_, err = builder.Load(game)
	if err != nil {
		fmt.Println("Load: ", err)
		t.FailNow()
	}
}
