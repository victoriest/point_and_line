package games

import (
	"testing"
)

func Test_NewPointAndLineGame(t *testing.T) {
	game := NewPointAndLineGame(2)
	t.Log(game)
}
