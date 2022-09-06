package core

func Filter[T any](list *[]T, pred func(int) bool) {
	if list == nil {
		return
	}

	l := *list

	filtered := l[:0]
	for i := range l {
		if pred(i) {
			filtered = append(filtered, l[i])
		}
	}

	list = &filtered
}

func Find[T any](s []T, f func(T) bool) int {
	for i := range s {
		if f(s[i]) {
			return i
		}
	}

	return -1
}

func FindReverse[T any](s []T, f func(T) bool) int {
	for i := len(s) - 1; i >= 0; i-- {
		if f(s[i]) {
			return i
		}
	}

	return -1
}

// FindSubstring return the beginnin of the substring
func FindSubstring[T comparable](text, query []T) int {
	if len(text) < len(query) {
		return -1
	}

outer:
	for i := range text[:len(text)-len(query)+1] {
		for j := range query {
			if text[i+j] != query[j] {
				continue outer
			}
		}

		return i
	}

	return -1
}
