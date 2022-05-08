package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/pkg/errors"
)

func main() {
	ebiten.SetWindowTitle("Snake")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetRunnableOnUnfocused(false)

	if err := ebiten.RunGame(new(game)); err != nil {
		panic(errors.Wrap(err, "error on run game"))
	}
}
