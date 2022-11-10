package core

import (
	"bufio"
	"bytes"
	"os"

	"github.com/pkg/errors"
)

// OpenFile opens a file with the given filename.
// If a file does not exist, it returns os.ErrNotExist.
func (e *E) OpenFile(filename string) error {
	e.filename = filename

	f, err := os.Open(filename)
	if errors.Is(err, os.ErrNotExist) {
		f, err = os.Create(filename)
		e.modified = true
	} else {
		e.modified = false
		err = nil
	}

	if err != nil {
		return errors.Wrapf(err, "opening file: %s", filename)
	}
	defer f.Close()

	e.rows = e.rows[:0]

	s := bufio.NewScanner(f)
	for i := 0; s.Scan(); i++ {
		line := s.Bytes()

		// strip off newline or cariage return
		bytes.TrimRightFunc(line, func(r rune) bool { return r == '\n' || r == '\r' })
		e.rows = append(e.rows, &Row{
			chars: []rune(string(line)),
		})

		e.updateRow(i)
	}

	if err := s.Err(); err != nil {
		return errors.Wrapf(err, "reading %s", e.filename)
	}

	if e.modified {
		e.cx = 0
		e.cy = 0
		e.rx = 0
	}

	return nil
}
