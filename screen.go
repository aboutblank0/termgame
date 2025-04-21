package termgame

import (
	"fmt"
	"strings"
)

type Screen struct {
	Enabled    bool
	Width      int
	Height     int
	cells      [][]Cell
	dirtyCells []Cell
}

type Cell struct {
	Color BackgroundColor
	x, y  int
}

const EMPTY = " "

func newScreen(width, height int) *Screen {
	screen := new(Screen)
	screen.Width = width
	screen.Height = height

	cells := make([][]Cell, height)
	dirtyCells := make([]Cell, 0, width * height)
	for y := range height {
		cells[y] = make([]Cell, width)
		for x := range cells[y] {
			cells[y][x] = Cell{x: x, y: y, Color: BlackColor}
			dirtyCells = append(dirtyCells, cells[y][x])
		}
	}
	screen.cells = cells
	screen.dirtyCells = dirtyCells

	return screen
}

func (s *Screen) SetCell(x, y int, color BackgroundColor) {
	cell := s.cells[y][x]
	cell.Color = color
	s.cells[y][x] = cell

	s.dirtyCells = append(s.dirtyCells, cell)
}

func (s *Screen) render() {
	var builder strings.Builder
	
	// Render only dirty cells
	for _, cell := range s.dirtyCells {
		builder.WriteString(GetMoveCursorCode(cell.x, cell.y))
		builder.WriteString(GetSetBackgroundColorCode(cell.Color))
		builder.WriteString(EMPTY)
	}

	fmt.Print(builder.String())

	s.dirtyCells = s.dirtyCells[:0]

	ResetAttributes()
}

func (s *Screen) enable() {
	SaveCursor()
	SaveScreen()
	SetCursorInvisible()
	EnableMouseTracking()
}

func (s *Screen) disable() {
	RestoreScreen()
	RestoreCursor()
	SetCursorVisible()
	DisableMouseTracking()
}
