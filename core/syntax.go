package core

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/mattn/go-runewidth"
)

type SyntaxHL uint8

// Syntax highlight enums
const (
	HLNormal SyntaxHL = iota + 1
	HLComment
	HLMlComment
	HLKeyword1
	HLKeyword2
	HLString
	HLNumber
	HLMatch
)

type EditorSyntax struct {
	// Name of the filetype displayed in the status bar.
	Filetype string
	// List of keywords to highlight.
	Keywords map[SyntaxHL][]string
	// Second highlight group
	Keywords2 []string
	// scs is a single-line comment start pattern (e.g. "//" for golang).
	// set to an empty string if comment highlighting is not needed.
	Scs string
	// mcs is a multi-line comment start pattern (e.g. "/*" for golang).
	Mcs string
	// mce is a multi-line comment end pattern (e.g. "*/" for golang).
	Mce string

	HighlightStrings bool
	HighlightNumbers bool
}

func (e *E) updateRow(y int) {
	var b strings.Builder
	row := e.rows[y]
	cols := 0
	for _, r := range row.chars {
		if r != '\t' {
			b.WriteRune(r)
			cols += runewidth.RuneWidth(r)
			continue
		}

		// each tab must advance the cursor forward at least one column
		b.WriteRune(' ')
		cols++

		// append spaces until we get to a tab stop
		for cols%e.cfg.Tabstop != 0 {
			b.WriteRune(' ')
			cols++
		}

	}

	row.render = b.String()
	e.updateHighlight(y)
}

func isSeparator(r rune) bool {
	return unicode.IsSpace(r) || strings.IndexRune(",.()+-/*=~%<>[]{}:;", r) != -1
}

func (e *E) updateHighlight(y int) {
	row := e.rows[y]

	row.hl = make([]SyntaxHL, utf8.RuneCountInString(row.render))
	for i := range row.hl {
		row.hl[i] = HLNormal
	}

	if e.syntax == nil {
		return
	}

	// whether the previous rune was a separator
	prevSep := true

	// zero when outside a string, set to the quote character ( ' or ")  in the string
	var strQuote rune

	// indicates whether we are inside a multi-line comment.
	inComment := y > 0 && e.rows[y-1].hasUnclosedComment

	idx := 0
	runes := []rune(row.render)
	for idx < len(runes) {
		r := runes[idx]

		prevHl := HLNormal
		if idx > 0 {
			prevHl = row.hl[idx-1]
		}

		// Single line comments
		if e.syntax.Scs != "" && strQuote == 0 && !inComment {
			if strings.HasPrefix(string(runes[idx:]), e.syntax.Scs) {
				for idx < len(runes) {
					row.hl[idx] = HLComment
					idx++
				}
				break
			}
		}

		// Multiline comments
		if e.syntax.Mcs != "" && e.syntax.Mce != "" && strQuote == 0 {
			if inComment {
				row.hl[idx] = HLMlComment
				if strings.HasPrefix(string(runes[idx:]), e.syntax.Mce) {
					for j := 0; j < len(e.syntax.Mce); j++ {
						row.hl[idx] = HLMlComment
						idx++
					}
					inComment = false
					prevSep = true
				} else {
					idx++
				}
				continue
			} else if strings.HasPrefix(string(runes[idx:]), e.syntax.Mcs) {
				for j := 0; j < len(e.syntax.Mcs); j++ {
					row.hl[idx] = HLMlComment
					idx++
				}
				inComment = true
				continue
			}
		}

		if e.syntax.HighlightStrings {
			if strQuote != 0 {
				row.hl[idx] = HLString
				// deal with escape quote when inside a string
				if r == '\\' && idx+1 < len(runes) {
					row.hl[idx+1] = HLString
					idx += 2
					continue
				}

				if r == strQuote {
					strQuote = 0
				}

				idx++
				prevSep = true
				continue
			} else {
				if r == '"' || r == '\'' {
					strQuote = r
					row.hl[idx] = HLString
					idx++
					continue
				}
			}
		}

		if e.syntax.HighlightNumbers {
			if unicode.IsDigit(r) && (prevSep || prevHl == HLNumber) ||
				r == '.' && prevHl == HLNumber {
				row.hl[idx] = HLNumber
				idx++
				prevSep = false
				continue
			}
		}

		if prevSep {
			if kw, hl := e.checkIfKeyword(runes[idx:]); kw != "" {
				end := idx + len(kw)
				for idx < end {
					row.hl[idx] = hl
					idx++
				}
			}
		}

		prevSep = isSeparator(r)
		idx++
	}

	changed := row.hasUnclosedComment != inComment
	row.hasUnclosedComment = inComment
	if changed && y+1 < len(e.rows) {
		e.updateHighlight(y + 1)
	}
}

func (e *E) checkIfKeyword(text []rune) (string, SyntaxHL) {
	for group := range e.syntax.Keywords {
		kw := isKeyword(e.syntax.Keywords[group], text)
		if len(kw) != 0 {
			return kw, group
		}
	}

	return "", 0
}

// Check if any of the keywords are a prefix of text, and also that it isn't
// just a substring of the a bigger word in text
func isKeyword(keywords []string, text []rune) string {
	for _, kw := range keywords {
		length := utf8.RuneCountInString(kw)
		if length > len(text) {
			continue
		}

		// check if we have a match
		if kw != string(text[:length]) {
			continue
		}

		// check that this is the entire word, either
		// there are no characters after this, or the
		// next character is a separator
		if length != len(text) && !isSeparator(text[length]) {
			continue
		}

		return kw
	}

	return ""
}
