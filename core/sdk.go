package core

import (
	"errors"
	"os"
)

func (e *E) Row(y int) []rune {
	return e.rows[y].chars
}

func (e *E) SetRow(y int, r []rune) {
	e.rows[y].chars = r
	e.updateRow(y)
}

func (e *E) NumRows() int {
	return len(e.rows)
}

func (e *E) ScreenBottom() int {
	return e.rowOffset + e.screenRows - 1
}

func (e *E) ScreenTop() int {
	return e.rowOffset + 1
}

func (e *E) Filename() string {
	return e.filename
}

func (e *E) Save() error {
	if len(e.filename) == 0 {
		return errors.New("file has no name")
	}
	return e.SaveTo(e.filename)
}

func (e *E) SaveTo(filename string) error {
	f, err := os.OpenFile(e.filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, row := range e.rows {
		if _, err := f.Write([]byte(string(row.chars))); err != nil {
			return err
		}
		if _, err := f.Write([]byte{'\n'}); err != nil {
			return err
		}
	}

	e.modified = false
	return nil
}

func (e *E) SetY(y int) {
	switch {
	case y < 0:
		e.cy = 0
	case y >= len(e.rows):
		e.cy = len(e.rows) - 1
	default:
		e.cy = y
	}

	e.cx = RxToCx(e.rows[e.cy].chars, e.cfg.Tabstop, e.rx)
}

func (e *E) SetX(x int) {
	switch {
	case x < 0:
		e.cx = 0
	case x > len(e.rows[e.cy].chars):
		e.cx = len(e.rows[e.cy].chars)
	default:
		e.cx = x
	}

	e.rx = CxToRx(e.rows[e.cy].chars, e.cfg.Tabstop, x)
}

func (e *E) SetRowOffset(y int) {
	if y < 0 {
		y = 0
	}
	e.rowOffset = y
}

func (e *E) SetColOffset(x int) {
	if x < 0 {
		x = 0
	}
	e.colOffset = x
}

func (e *E) X() int {
	return e.cx
}

func (e *E) Y() int {
	return e.cy
}


//func (s *SDK) InsertChars(y, x int, chars []rune) {
//	row := s.rows[s.cy].Copy()
//	row.chars = append(row.chars[x:], append(chars, row.chars[:x]...)...)
//
//	s.changeChan <- Edit(y, &Row{chars: row.chars})
//}
//

//func (s *SDK) DeleteRows(from, to int) {
//	log.Printf("sent to the change channel, from=%d, to=%d", from, to)
//	s.changeChan <- Delete(from, to)
//}

//func (s *SDK) FindGeneral(x1, y1 int, f func([]rune) int) (x, y int) {
//	if x = f(s.rows[y1].chars[x1:]); x != -1 {
//		return x1 + x, y1
//	}
//
//	for y = y1 + 1; y < len(s.rows); y++ {
//		if x = f(s.rows[y].chars); x != -1 {
//			return x, y
//		}
//	}
//
//	return -1, -1
//}
//
//func (s *SDK) FindBackGeneral(x1, y1 int, f func([]rune) int) (x, y int) {
//	if x = f(s.rows[y1].chars[:x1]); x != -1 {
//		return x, y1
//	}
//
//	for y = y1 - 1; y >= 0; y-- {
//		if x = f(s.rows[y].chars); x != -1 {
//			return x, y
//		}
//	}
//
//	return -1, -1
//}
//
//func (s *SDK) goForwardOneStep() (x, y int) {
//	if s.cx == len(s.rows[s.cy].chars) {
//		if s.cy == len(s.rows) {
//			return s.cx, s.cy
//		}
//
//		return 0, s.cy + 1
//	}
//
//	return s.cx + 1, s.cy
//}
//
//func (s *SDK) goBackOneStep() (x, y int) {
//	if s.e.X() == 0 {
//		if s.e.Y() == 0 {
//			return 0, s.e.Y()
//		}
//
//		return len(s.e.Row(s.e.Y()-1)), s.e.Y() - 1
//	}
//
//	return s.e.X() - 1, s.e.Y()
//}
//
//// maybe just get rid of this offset
//func findSubstringBack(text, query []rune, offset int) int {
//	if len(text) < len(query) {
//		return -1
//	}
//
//	// Make sure text[i+j] doesn't overflow
//	if offset > len(text)-len(query) {
//		offset = len(text) - len(query)
//	}
//
//outer:
//	for i := offset; i >= 0; i-- {
//		log.Printf("text: %s, i: %d", string(text[i:i+len(query)]), i)
//		for j := range query {
//			if text[i+j] != query[j] {
//				continue outer
//			}
//		}
//
//		return i
//	}
//
//	return -1
//}
//
//
//func (s *SDK) Save() error {
//	s.StaticPrompt("Save as: ", s.e.SaveTo)
//	return nil
//}

//func (s *SDK) FindInteractive() {
//	savedCx := s.cx
//	savedCy := s.cy
//	savedColOffset := s.colOffset
//	savedRowOffset := s.rowOffset
//
//	var query []rune
//
//	onKeyPress := func(k Key) (string, bool) {
//		switch k {
//		case keyDelete, keyBackspace:
//			if len(query) != 0 {
//				query = query[:len(query)-1]
//
//				// This forces the editor to search again to
//				// see if the current word is indeed the
//				// closest match..cyes making a stack containing
//				// the previous matches would be better, but it
//				// is somewhat unecessary at the moment
//				s.cx = savedCx
//				s.cy = savedCy
//
//			}
//		case keyEscape, Key(ctrl('q')):
//			// restore cursor position when the user cancels search
//			s.cx = savedCx
//			s.cy = savedCy
//			s.colOffset = savedColOffset
//			s.rowOffset = savedRowOffset
//
//			s.SetStatusLine("")
//
//			return "", true
//		case keyEnter, keyCarriageReturn:
//			s.SetStatusLine("")
//			s.lastSearch = query
//
//			return "", true
//		default:
//			if isPrintable(k) {
//				query = append(query, rune(k))
//			}
//		}
//
//		x, y := s.Find(s.cx, s.cy, query)
//		if x == -1 {
//			s.cx = savedCx
//			s.cy = savedCy
//			s.colOffset = savedColOffset
//			s.rowOffset = savedRowOffset
//
//			return string(query), false
//		}
//
//		// Set cursor to beginning of match
//		s.SetX(x)
//		s.SetY(y)
//
//		// Try to make the text in the middle of the screen
//		s.SetRowOffset(s.cy - s.screenRows/2)
//
//		return string(query), false
//	}
//
//	s.Prompt("Search: ", onKeyPress)
//}

//func (s *SDK) Find(x1, y1 int, query []rune) (x, y int) {
//	return s.FindGeneral(x1, y1, func(s []rune) int {
//		return FindSubstring(s, query)
//	})
//}
//
//func (s *SDK) FindBack(x1, y1 int, query []rune) (x, y int) {
//	return s.FindBackGeneral(x1, y1, func(s []rune) int {
//		return findSubstringBack(s, query, len(s))
//	})
//}
//
//func (s *SDK) FindRegex(x1, y1 int, query string) (x, y int) {
//	r, err := regexp.Compile(query)
//	if err != nil {
//		return -1, -1
//	}
//
//	return s.FindGeneral(x1, y1, func(s []rune) int {
//		loc := r.FindIndex([]byte(string(s)))
//		if len(loc) == 0 {
//			return -1
//		}
//
//		return loc[0]
//	})
//}
//
//func (s *SDK) FindBackRegex(x1, y1 int, query string) (x, y int) {
//	r, err := regexp.Compile(query)
//	if err != nil {
//		return -1, -1
//	}
//
//	return s.FindBackGeneral(x1, y1, func(s []rune) int {
//		loc := r.FindAllIndex([]byte(string(s)), 100)
//		log.Printf("loc: %v", loc)
//		if len(loc) == 0 {
//			return -1
//		}
//
//		return loc[len(loc)-1][0]
//	})
//}
//
