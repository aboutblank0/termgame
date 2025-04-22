package termgame

import (
	"fmt"
	"os"
	"time"

	"golang.org/x/term"
)

type TerminalGame struct {
	screen   *Screen
	input    *Input
	elements []GameElement
	running  bool
	fd       int
}

type GameElement interface {
	Update(screen *Screen, delta float64, input Input)
}

const TARGET_FPS = 60
const TARGET_FRAME_TIME = time.Second / TARGET_FPS

func NewGame() (*TerminalGame, error) {
	fd := int(os.Stdin.Fd())
	width, height, err := term.GetSize(fd)
	if err != nil {
		return nil, err
	}

	screen := newScreen(width, height)
	input := newInput()

	return &TerminalGame{
		running:  false,
		screen:   screen,
		input:    input,
		fd:       fd,
		elements: make([]GameElement, 0),
	}, nil
}

func (game *TerminalGame) AddElement(el GameElement) {
	game.elements = append(game.elements, el)
}

func (game *TerminalGame) Start() {
	game.screen.enable()
	defer game.screen.disable()

	//Raw terminal
	oldState, err := term.MakeRaw(game.fd)
	if err != nil {
		panic(err)
	}
	defer term.Restore(game.fd, oldState)

	//Run loop
	run(game)
}

func run(game *TerminalGame) {
	game.running = true

	last := time.Now()
	for game.running {
		now := time.Now()
		delta := now.Sub(last).Seconds()
		last = now

		game.input.refresh()

		update(game, delta)

		game.screen.render()

		//Print FPS
		fmt.Printf("\x1b[%d;1H\x1b[2K FPS: %d", game.screen.Height, int(1.0/delta))

		//Target FPS
		elapsed := time.Since(now)
		if sleepTime := TARGET_FRAME_TIME - elapsed; sleepTime > 0 {
			time.Sleep(sleepTime)
		}
	}
}

func update(game *TerminalGame, delta float64) {
	if game.input.GetKeyDown('q') {
		game.running = false
	}

	for _, element := range game.elements {
		element.Update(game.screen, delta, *game.input)
	}
}
