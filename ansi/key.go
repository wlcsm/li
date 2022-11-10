package ansi

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
	LeftArrowKey Key = 1000 + iota
	RightArrowKey
	UpArrowKey
	DownArrowKey
	DeleteKey
	PageUpKey
	PageDownKey
	HomeKey
	EndKey
)

// Returns the ASCII for when the character has been pressed with the CTRL key
func Ctrl(char byte) Key {
	return Key(char & 0x1f)
}

type Direction int8

const (
	DirectionUp Direction = iota + 1
	DirectionDown
	DirectionLeft
	DirectionRight
)
