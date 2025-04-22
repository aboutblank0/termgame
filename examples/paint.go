package main

import (
	"aboutblank0/termgame"
)

type Painter struct {
	Painting bool
	Game  *termgame.TerminalGame
	Color termgame.Color
}

func main() {
	game, err := termgame.NewGame()
	if err != nil {
		panic(err)
	}

	paint := new(Painter)
	paint.Color = termgame.RedColor
	paint.Game = game
	paint.Painting = false

	game.AddElement(paint)
	game.Start()
}

func (painter *Painter) Update(screen *termgame.Screen, delta float64, input termgame.Input) {
	painter.Painting = input.GetMouseButtonDown(0)

	if painter.Painting {
		screen.SetCell(input.Mouse.X, input.Mouse.Y, painter.Color)
	}
}
