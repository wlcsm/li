// Defines the recognised keys for the li editor as well as associated ANSI
// parser.
package ansi

import (
	"bufio"
	"errors"
	"io"
)

// Simple ANSI decoder. Only working on a subset of ANSI codes relevant for
// usage inside the editor program. For example, arrow keys, home key, and more
// are supported, but colour codes are not.
// 
// Will only decode the keys defined in the Key type in this package
type Decoder struct {
	rd *bufio.Reader
	// The last rune processed in the escape sequence
	state rune
}

var ErrInvalidEscape = errors.New("invalid escape code")

func NewDecoder(r io.Reader) Decoder {
	return Decoder{
		rd: bufio.NewReader(r),
	}
}

func (d *Decoder) Decode() (Key, error) {
	for {
		ru, n, err := d.rd.ReadRune()
		if err != nil {
			return 0, err
		}
		if n == 0 {
			continue
		}

		switch d.state {
		case 0:
			// hot path: normal character input
			if ru != rune(EscapeKey) {
				return Key(ru), nil
			} else {
				d.state = ru
			}
		case rune(EscapeKey):
			if ru == '[' {
				d.state = ru
			} else {
				return 0, ErrInvalidEscape
			}
		case '[':
			d.state = 0

			switch ru {
			case 'A':
				return UpArrowKey, nil
			case 'B':
				return DownArrowKey, nil
			case 'C':
				return RightArrowKey, nil
			case 'D':
				return LeftArrowKey, nil
			case 'F':
				return EndKey, nil
			case 'H':
				return HomeKey, nil
			case '1', '2', '3', '4', '5', '6', '7', '8':
				d.state = ru
			default:
				return 0, ErrInvalidEscape
			}
		case '1', '2', '3', '4', '5', '6', '7', '8':
			state := d.state
			d.state = 0

			if ru != '~' {
				return 0, ErrInvalidEscape
			}

			switch state {
			case '1', '7':
				return HomeKey, nil
			case '3':
				return DeleteKey, nil
			case '4', '8':
				return EndKey, nil
			case '5':
				return PageUpKey, nil
			case '6':
				return PageDownKey, nil
			}
		}
	}
}
