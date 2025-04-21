package termgame

import (
	"fmt"
	"os"
	"time"

	"golang.org/x/term"
)

type TerminalGame struct {
	Screen   *Screen
	Elements []GameElement
	running  bool
	fd       int
}

type GameElement interface {
	Update(delta float64, input Input)
	Render(screen *Screen)
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
	return &TerminalGame{
		running:  false,
		Screen:   screen,
		fd:       fd,
		Elements: make([]GameElement, 0),
	}, nil
}

func (game *TerminalGame) AddElement(el GameElement) {
	game.Elements = append(game.Elements, el)
}

func (game *TerminalGame) Start() {
	game.Screen.enable()
	defer game.Screen.disable()

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

	inputCh := getInputChannel()

	last := time.Now()
	for game.running {
		now := time.Now()
		delta := now.Sub(last).Seconds()
		last = now

		input := getInput(inputCh)

		update(game, delta, input)

		render(game)
	
		//Print FPS
		fmt.Printf("\x1b[%d;1H\x1b[2KFPS: %d", game.Screen.Height+1, int(1.0/delta))

		//Target FPS 
		elapsed := time.Since(now)
		if sleepTime := TARGET_FRAME_TIME - elapsed; sleepTime > 0 {
			time.Sleep(sleepTime)
		}
	}
}

func update(game *TerminalGame, delta float64, input Input) {
	if input.Key == 'q' {
		game.running = false
	}

	for _, element := range game.Elements {
		element.Update(delta, input)
	}
}

func render(app *TerminalGame) {
	for _, element := range app.Elements {
		element.Render(app.Screen)
	}

	app.Screen.render()
}
