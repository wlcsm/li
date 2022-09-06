package core

import (
	"os"
	"strconv"
)

type Key int32

// Assign an arbitrary large number to the following special keys
// to avoid conflicts with the normal keys.
const (
	EnterKey          Key = 10
	CarriageReturnKey Key = 13
	BackspaceKey      Key = 127
	EscapeKey         Key = '\x1b'

	// These keys are accomplished with escape sequences in the terminal,
	// and so do not normally fit into a single rune character. As such we
	// have mapped the to runes so that everything just processes runes
	LeftArrowKey Key = iota + 1000
	RightArrowKey
	UpArrowKey
	DownArrowKey
	DeleteKey
	PageUpKey
	PageDownKey
	HomeKey
	EndKey
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

// Returns the ASCII for when the character has been pressed with the CTRL key
func Ctrl(char byte) Key {
	return Key(char & 0x1f)
}

var escapeCodeToKey = map[string]Key{
	"\x1b[A":  UpArrowKey,
	"\x1b[B":  DownArrowKey,
	"\x1b[C":  RightArrowKey,
	"\x1b[D":  LeftArrowKey,
	"\x1b[1~": HomeKey,
	"\x1b[7~": HomeKey,
	"\x1b[H":  HomeKey,
	"\x1bOH":  HomeKey,
	"\x1b[4~": EndKey,
	"\x1b[8~": EndKey,
	"\x1b[F":  EndKey,
	"\x1bOF":  EndKey,
	"\x1b[3~": DeleteKey,
	"\x1b[5~": PageUpKey,
	"\x1b[6~": PageDownKey,
}

type Direction int8

const (
	DirectionUp Direction = iota + 1
	DirectionDown
	DirectionLeft
	DirectionRight
)


func IsPrintable(k Key) bool {
	return !unicode.IsControl(rune(k)) && unicode.IsPrint(rune(k)) && !IsArrowKey(k)
}

func IsArrowKey(k Key) bool {
	return k == UpArrowKey || k == RightArrowKey || k == DownArrowKey || k == LeftArrowKey
}

