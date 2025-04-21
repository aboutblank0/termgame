package termgame

import (
	"fmt"
	"math/rand"
)

type BackgroundColor int

const (
	BlackColor   BackgroundColor = 0
	RedColor                     = 1
	GreenColor                   = 2
	YellowColor                  = 3
	BlueColor                    = 4
	MagentaColor                 = 5
	CyanColor                    = 6
	WhiteColor                   = 7
	MaxColor                     = 8
)

func GetRandomColor() BackgroundColor {
	randColor := rand.Intn(MaxColor) + 1
	return BackgroundColor(randColor)
}

const ESC = "\x1b"

func printAnsi(ansi string) { fmt.Printf("%s%s", ESC, ansi) }

// Screen
func EraseScreen()   { printAnsi("[2J") }
func SaveScreen()    { printAnsi("[?47h") }
func RestoreScreen() { printAnsi("[?47l") }

// Cursor
func SaveCursor()                       { printAnsi("[s") }
func RestoreCursor()                    { printAnsi("[u") }
func SetCursorInvisible()               { printAnsi("[?25l") }
func SetCursorVisible()                 { printAnsi("[?25h") }
func MoveCursor(x, y int)               { printAnsi(fmt.Sprintf("[%d;%dH", y, x)) }
func GetMoveCursorCode(x, y int) string { return fmt.Sprintf("%s[%d;%dH", ESC, y, x) }

// Back/Foreground
func SetBackgroundColor(color BackgroundColor) { printAnsi(fmt.Sprintf("[48;5;%dm", color)) }
func GetSetBackgroundColorCode(color BackgroundColor) string {
	return fmt.Sprintf("%s[48;5;%dm", ESC, color)
}

// Mouse
func EnableMouseTracking()  { printAnsi("[?1000h") }
func DisableMouseTracking() { printAnsi("[?1000l") }

func ResetAttributes() { printAnsi("[0m") }
