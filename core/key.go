package core

import (
	"os"
	"strconv"
	"unicode"

	"codeberg.org/wlcsm/li/ansi"
)

var (
	ShowCursor      = []byte("\x1b[?25h")
	HideCursor      = []byte("\x1b[?25l")
	CursorToTopLeft = []byte("\x1b[H")
)

func RepositionCursor() {
	os.Stdout.WriteString(RepositionCursorCode)
}

func ClearScreen() {
	os.Stdout.WriteString(ClearScreenCode)
}

type EscapeCodes string

const (
	RepositionCursorCode = "\x1b[H"
	ResetColorCode       = "\x1b[39m"
	ClearLineCode        = "\x1b[K"
	ClearScreenCode      = "\x1b[2J"
)

const (
	ClearColor    = 39
	InvertedColor = 7
)

var ClearFormatting = []byte("\x1b[m")

func getColor(c int) []byte {
	return []byte("\x1b[" + strconv.Itoa(c) + "m")
}

func IsPrintable(k ansi.Key) bool {
	return !unicode.IsControl(rune(k)) && unicode.IsPrint(rune(k)) && !IsArrowKey(k)
}

func IsArrowKey(k ansi.Key) bool {
	return k == ansi.UpArrowKey || k == ansi.RightArrowKey || k == ansi.DownArrowKey || k == ansi.LeftArrowKey
}
