package core

import "io"

func SwitchToAlternateScreen(w io.Writer) {
	w.Write([]byte("\033[?1049h"))
}

func SwitchBackFromAlternateScreen(w io.Writer) {
	w.Write([]byte("\033[?1049l"))
}
