package core

import (
	"fmt"
	"io"
	"os"
	"unicode"
	"unicode/utf8"

	"github.com/mattn/go-runewidth"
)

func (e *E) drawRows(w io.Writer) {
	for y := 0; y < e.screenRows; y++ {
		e.drawRow(w, y)

		w.Write([]byte(ClearLineCode))
		w.Write([]byte("\r\n"))
	}
}

func (e *E) drawRow(w io.Writer, y int) {
	filerow := y + e.rowOffset
	if filerow >= len(e.rows) {
		w.Write([]byte("~"))
		return
	}

	var (
		line string
		hl   []SyntaxHL
	)

	// Use the offset to remove the first part of the render string
	row := e.rows[filerow]
	if runewidth.StringWidth(row.render) > e.colOffset {
		line = utf8Slice(row.render, e.colOffset, utf8.RuneCountInString(row.render))
		hl = e.rows[filerow].hl[e.colOffset:]
	}

	// Use the number of columns to truncate the end
	if runewidth.StringWidth(line) > e.screenCols {
		line = runewidth.Truncate(line, e.screenCols, "")
		hl = hl[:utf8.RuneCountInString(line)]
	}

	currentColor := -1
	i := 0
	for _, r := range line {
		if unicode.IsControl(r) {
			// deal with non-printable characters (e.g. Ctrl-A)
			sym := '?'
			if r < 26 {
				sym = '@' + r
			}

			w.Write(getColor(InvertedColor))
			w.Write([]byte(string(sym)))
			w.Write(ClearFormatting)

			// restore the current color
			if currentColor != -1 {
				w.Write(getColor(currentColor))
			}
		} else {
			if color := e.syntaxToColor(hl[i]); color != currentColor {
				currentColor = color
				w.Write(getColor(color))
			}

			w.Write([]byte(string(r)))
		}
		i++
	}

	w.Write(getColor(ClearColor))
}

func (e *E) syntaxToColor(hl SyntaxHL) int {
	color, ok := e.colorscheme[hl]
	if !ok {
		return 37
	}
	return color
}

func (e *E) Render(line int) {
	os.Stdout.Write(HideCursor)

	e.updateRow(line)

	// line is out of bounds
	if line < e.colOffset || e.colOffset+e.screenCols < line {
		return
	}

	e.positionCursor(0, line)
	e.drawRow(os.Stdout, line)

	e.positionCursor(e.cx, e.cy)
	os.Stdout.Write(ShowCursor)
}

func (e *E) positionCursor(x, y int) {
	d := x
	if d > e.rows[y].visibleLength(e.cfg.Tabstop) {
		d = e.rows[y].visibleLength(e.cfg.Tabstop)
	}

	// Ensure the rx is not inside a tabstop
	d = e.rows[y].roundToNearestRealChar(d, e.cfg.Tabstop)

	// position the cursor
	os.Stdout.WriteString(fmt.Sprintf("\x1b[%d;%dH", (y-e.rowOffset)+1, (d-e.colOffset)+1))
}

func (e *E) FullRender() {
	e.scroll()

	os.Stdout.Write(HideCursor)
	os.Stdout.Write(CursorToTopLeft)

	e.drawRows(os.Stdout)
	e.drawStatusBar(os.Stdout)
	e.drawMessageBar(os.Stdout)

	d := e.rx
	if d > e.rows[e.cy].visibleLength(e.cfg.Tabstop) {
		d = e.rows[e.cy].visibleLength(e.cfg.Tabstop)
	}

	// Ensure the rx is not inside a tabstop
	d = e.rows[e.cy].roundToNearestRealChar(d, e.cfg.Tabstop)

	// position the cursor
	os.Stdout.WriteString(fmt.Sprintf("\x1b[%d;%dH", (e.cy-e.rowOffset)+1, (d-e.colOffset)+1))

	// show the cursor
	os.Stdout.Write(ShowCursor)
}

// utf8Slice slice the given string by utf8 character.
func utf8Slice(s string, start, end int) string {
	return string([]rune(s)[start:end])
}

var ClearFromCusorToEndOfLine = []byte("\x1b[K")

func (e *E) drawMessageBar(w io.Writer) {
	msg := e.statusMsg
	if runewidth.StringWidth(msg) > e.screenCols {
		msg = runewidth.Truncate(msg, e.screenCols, "...")
	}

	w.Write(ClearFromCusorToEndOfLine)
	w.Write([]byte(msg))
}

// Cursor position (which is calculated in runes) to the visual position
func CxToRx(row []rune, tabstop int, cx int) int {
	if len(row) == 0 {
		return 0
	}

	if cx > len(row) {
		cx = len(row)
	}

	rx := 0
	for _, r := range row[:cx] {
		if r == '\t' {
			rx += (tabstop) - (rx % tabstop)
		} else {
			rx += runewidth.RuneWidth(r)
		}
	}

	return rx
}

func RxToCx(chars []rune, tabstop, rx int) int {
	if len(chars) == 0 {
		return 0
	}

	curRx := 0
	for i, r := range chars {
		if r == '\t' {
			curRx += tabstop - (curRx % tabstop)
		} else {
			curRx += runewidth.RuneWidth(r)
		}

		if curRx > rx {
			return i
		}
	}

	// If Rx exceeds the length of the row, then put the cursor at the
	// end
	return len(chars)
}

// Scroll so that the cursor is still visible
func (e *E) scroll() {
	d := e.rx
	if d > len(e.rows[e.cy].chars) {
		d = len(e.rows[e.cy].chars)
	}
	// scroll up if the cursor is above the visible window.
	if e.cy < e.rowOffset {
		e.rowOffset = e.cy
	}
	// scroll down if the cursor is below the visible window.
	if e.cy >= e.rowOffset+e.screenRows {
		e.rowOffset = e.cy - e.screenRows + 1
	}
	// scroll left if the cursor is left of the visible window.
	if d < e.colOffset {
		e.colOffset = d
	}
	// scroll right if the cursor is right of the visible window.
	if d >= e.colOffset+e.screenCols {
		e.colOffset = d - e.screenCols + 1
	}
}

// Number of cols a line takes up
func (r Row) visibleLength(tabstop int) int {
	return CxToRx(r.chars, tabstop, len(r.chars))
}

// Round the rx (to the left) to the nearest character so that it is not inside
// a tabstop
func (r Row) roundToNearestRealChar(rx, tabstop int) int {
	acc := 0
	for _, c := range r.chars {
		if c == '\t' {
			acc += tabstop
		} else {
			acc += 1
		}

		if acc == rx {
			return rx
		}

		if acc > rx {
			return rx - tabstop
		}
	}

	return rx
}
