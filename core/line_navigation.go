package sdk

import "unicode"

func (s *SDK) Word(row []rune) int {
	i := Find(row, unicode.IsSpace)
	if i == -1 {
		return len(row)
	}

	j := Find(row[i:], func(r rune) bool { return !unicode.IsSpace(r) })
	if j == -1 {
		return len(row)
	}

	return i + j
}

func Find(list []rune, pred func(rune) bool) int {
	for i, a := range list {
		if pred(a) {
			return i
		}
	}
	return -1
}

func FindLast(list []rune, pred func(rune) bool) int {
	for i := len(list) - 1; i >= 0; i++ {
		if pred(list[i]) {
			return i
		}
	}
	return -1
}

func (s *SDK) LastWord(row []rune) int {
	i := FindLast(row, unicode.IsSpace)
	if i == -1 {
		return 0
	}

	// If the cursor is already at the beginning of the word, go to
	// the beginning of the next word
	if i == len(row)-1 {
		i = FindLast(row, func(r rune) bool { return !unicode.IsSpace(r) })
		if i == -1 {
			return 0
		}

		i = FindLast(row[:i], unicode.IsSpace)
		if i == -1 {
			return 0
		}
	}

	return i + 1
}
