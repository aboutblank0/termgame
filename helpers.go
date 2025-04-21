package termgame

import (
	"fmt"
	"math/rand"
)

type Color int

const (
	BlackColor   Color = 0
	RedColor           = 1
	GreenColor         = 2
	YellowColor        = 3
	BlueColor          = 4
	MagentaColor       = 5
	CyanColor          = 6
	WhiteColor         = 7

	MaxColor = 8
)

func GetRandomColor() Color {
	randColor := rand.Intn(MaxColor) + 1
	return Color(randColor)
}

const ESC = 0x1b

func printAnsi(ansi string) { fmt.Printf("%c%s", ESC, ansi) }

// Screen
func EraseScreen()   { printAnsi("[2J") }
func SaveScreen()    { printAnsi("[?47h") }
func RestoreScreen() { printAnsi("[?47l") }

// Cursor
func SaveCursor()                       { printAnsi("[ 7") }
func RestoreCursor()                    { printAnsi("[ 8") }
func SetCursorInvisible()               { printAnsi("[?25l") }
func SetCursorVisible()                 { printAnsi("[?25h") }
func MoveCursor(x, y int)               { printAnsi(fmt.Sprintf("[%d;%dH", y, x)) }
func GetMoveCursorCode(x, y int) string { return fmt.Sprintf("%c[%d;%dH", ESC, y, x) }

// Back/Foreground
func SetBackgroundColor(color Color)               { printAnsi(fmt.Sprintf("[48;5;%dm", color)) }
func GetSetBackgroundColorCode(color Color) string { return fmt.Sprintf("%c[48;5;%dm", ESC, color) }

// Mouse
func EnableMouseTracking()  {
	printAnsi("[?1000h") 
	printAnsi("[?1006h") 
}

func DisableMouseTracking() {
	printAnsi("[?1000l") 
	printAnsi("[?1006l") 
}


func ResetAttributes() { printAnsi("[0m") }
