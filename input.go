package termgame

import (
	"os"
	"strconv"
	"strings"
)

type Input struct {
	Key   byte
	Mouse MouseInput
}

type MouseInput struct {
	X       int
	Y       int
	Button  ButtonType
	Pressed bool
}

type ButtonType int

const (
	LeftClick   = 0
	MiddleClick = 1
	RightClick  = 2
)

func getInputChannel() <-chan []byte {
	inputCh := make(chan []byte)
	go readInputLoop(inputCh)
	return inputCh
}

func getInput(ch <-chan []byte) Input {
	select {
	case b := <-ch:
		parseSGR(b)
		switch len(b) {
		case 1:
			return Input{Key: b[0]}
		default:
			//Kind of bad to assume any 6 byte slice is guaranteed to be a mouse input, but... meh
			mouseInput := parseSGR(b)
			return Input{Mouse: mouseInput}
		}
	default: //Do nothing
	}
	return Input{}
}

// TODO: More efficient parsing
func parseSGR(bytes []byte) MouseInput {
	str := string(bytes)

	if !strings.HasPrefix(str, "\x1b[<") {
		return MouseInput{}
	}

	body := str[3 : len(str)-1]
	parts := strings.Split(body, ";")
	if len(parts) != 3 {
		return MouseInput{}
	}

	button, err1 := strconv.Atoi(parts[0])
	x, err2 := strconv.Atoi(parts[1])
	y, err3 := strconv.Atoi(parts[2])
	if err1 != nil || err2 != nil || err3 != nil {
		return MouseInput{}
	}
	eventType := str[len(str)-1] // 'M' or 'm'

	return MouseInput{
		X:       x,
		Y:       y,
		Button:  ButtonType(button),
		Pressed: eventType == 'M',
	}
}

func readInputLoop(inputCh chan<- []byte) {
	inputBuffer := make([]byte, 32)
	for {
		n, err := os.Stdin.Read(inputBuffer)
		if err == nil && n > 0 {
			data := make([]byte, n)
			copy(data, inputBuffer[:n])
			inputCh <- data
		}
	}
}
