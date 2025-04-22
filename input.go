package termgame

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

// TODO: Keep some sort of state for input. So that we can easily check when mouse is up/down.
// Probably keep a pointer to Input in TerminalGame struct
type Input struct {
	inputChannel        chan []byte
	keys                [256]bool
	mouseButtons        [3]bool
	frameMousePositions []Position

	Mouse Position
}

type Position struct {
	X, Y int
}

const BUFFER_SIZE = 512

func newInput() *Input {
	inputCh := make(chan []byte, BUFFER_SIZE)
	go readInputLoop(inputCh)

	return &Input{
		inputChannel: inputCh,
	}
}

func (in Input) GetKeyDown(key byte) bool {
	return in.keys[key]
}

func (in Input) GetMouseButtonDown(button int) bool {
	return in.mouseButtons[button]
}

// Desperately needs a new name lol.
// This function returns all the cells the mouse moved "past" in between the last frame and this one.
func (in Input) GetAllFrameMousePositions() []Position {
	return in.frameMousePositions
}

func (in *Input) refresh() {
	var keys [256]bool
	var finalMousePosition *Position
	seenMousePositions := make(map[Position]bool)
	var uniqueMousePositions []Position

	for {
		select {
		case b := <-in.inputChannel:
			switch len(b) {
			case 1:
				keys[b[0]] = true
			default:
				mouseInput, valid := parseSGR(b)
				if valid {
					if mouseInput.button <= 2 { //Mouse move "buttons"
						in.mouseButtons[mouseInput.button] = (mouseInput.eventType == 'M')
					}

					if _, exists := seenMousePositions[mouseInput.pos]; !exists {
						seenMousePositions[mouseInput.pos] = true
						uniqueMousePositions = append(uniqueMousePositions, mouseInput.pos)
					}

					finalMousePosition = &mouseInput.pos
				}
			}
		default: //No more data in channel
			in.keys = keys
			in.frameMousePositions = uniqueMousePositions
			if finalMousePosition != nil {
				in.Mouse = *finalMousePosition
			}
			return
		}
	}
}

type SGREvent struct {
	button    int
	eventType byte
	pos       Position
}

// TODO: More efficient parsing (don't use strings)
// Also, do this in the input goroutine, so that parsing doesn't "add" to the frame time
func parseSGR(bytes []byte) (SGREvent, bool) {
	str := string(bytes)

	if !strings.HasPrefix(str, "\x1b[<") {
		panic(errors.New("Parsing input []byte that is not SGR ?"))
	}

	body := str[3 : len(str)-1]
	parts := strings.Split(body, ";")
	if len(parts) != 3 {
		return SGREvent{}, false
	}

	button, err1 := strconv.Atoi(parts[0])
	x, err2 := strconv.Atoi(parts[1])
	y, err3 := strconv.Atoi(parts[2])
	if err1 != nil || err2 != nil || err3 != nil {
		return SGREvent{}, false
	}

	eventType := str[len(str)-1] // 'M' or 'm'
	return SGREvent{
		button:    button,
		pos:       Position{X: x, Y: y},
		eventType: eventType,
	}, true
}

func readInputLoop(inputCh chan<- []byte) {
	inputBuffer := make([]byte, 64)
	for {
		n, err := os.Stdin.Read(inputBuffer)
		if err == nil && n > 0 {
			data := make([]byte, n)
			copy(data, inputBuffer[:n])
			inputCh <- data
		}
	}
}
