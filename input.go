package termgame

import (
	"os"
)

type Input struct {
	Key   byte
	Mouse *MouseInput
}

type MouseInput struct {
	X      int
	Y      int
	Button int
}

func getInputChannel() <-chan []byte {
	inputCh := make(chan []byte)
	go readInputLoop(inputCh)
	return inputCh
}

func getInput(ch <-chan []byte) Input {
	select {
	case b := <-ch:
		switch len(b) {
		case 1:
			return Input{Key: b[0]}
		case 6:
			//Kind of bad to assume any 6 byte slice is guaranteed to be a mouse input, but... meh
			mouseInput := getMouseInput([6]byte(b))
			return Input{Mouse: &mouseInput}

		}
	default: //Do nothing
	}

	return Input{}
}

func getMouseInput(bytes [6]byte) MouseInput {
	button := bytes[3] - 32
	x := bytes[4] - 32
	y := bytes[5] - 32
	return MouseInput{X: int(x), Y: int(y), Button: int(button)}
}

func readInputLoop(inputCh chan<- []byte) {
	inputBuffer := make([]byte, 6)
	for {
		n, err := os.Stdin.Read(inputBuffer)
		if err == nil && n > 0 {
			data := make([]byte, n)
			copy(data, inputBuffer[:n])
			inputCh <- data
		}
	}
}
