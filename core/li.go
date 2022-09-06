package core

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime/debug"
	"syscall"

	"github.com/mattn/go-runewidth"
	"github.com/pkg/errors"
	"golang.org/x/term"
)

var (
	ErrPromptCanceled = errors.New("user canceled the input prompt")
	ErrQuitEditor     = errors.New("quit editor")
)

var (
	LogFile = "/var/log/li.log"
)

// E is the editor kernel
//
// Terminology:
// +-----------------------------------------+
// |             ^                           |
// |  rowOffset  |   *Entire File Contents*  |
// |             _                           |
// |             ^ +-----------------+       |
// |             | |                 |       |
// |             | |                 |       |
// |             | | *Editor Window* |       |
// |  screenRows | |                 |       |
// |             | |                 |       |
// |             | |                 |       |
// |             _ +-----------------+       |
// |<-------------><----------------->       |
// |  colOffset       screenCols             |
// +-----------------------------------------+
//
// E is a very lightweight kernel that handles:
// * Terminal configuration
// * Rendering
// * Line editing
// * Input processing
// * File IO
type E struct {
	signals chan os.Signal
	keys    chan Key
	Errs    chan error

	// cursor coordinates
	cx, cy int // cx is an index into Row.chars
	rx     int // rx is an index into []rune(Row.render)

	filetypeLookup func(string) *EditorSyntax

	// Row offset is the number of rows above the row on the top of the screen
	// Offset is calculated in the number of runes
	rowOffset int
	colOffset int

	// screen size
	screenRows int
	screenCols int

	// file content
	rows []*Row

	// whether or not the file has been modified
	modified bool

	filename string

	// status message and time the message was set
	statusMsg string

	// General settings like tabstop
	cfg DisplayConfig

	// specify which syntax highlight to use.
	syntax *EditorSyntax

	colorscheme map[SyntaxHL]int

	// Callbacks.
	// Currently only has the open file callback
	callbacks Callbacks
}

type Callbacks struct {
	FileOpen func(e *E, filename string) error
}

type DisplayConfig struct {
	Tabstop int
}

type Row struct {
	// Raw character data for the row as an array of runes.
	chars []rune
	// Actual chracters to draw on the screen. It is primarily about
	// expanding the tab character to a variable number of spaces
	render string
	// Syntax highlight value for each rune in the render string.
	hl []SyntaxHL
	// Indicates whether this row has unclosed multiline comment.
	hasUnclosedComment bool
}

type EditorConf struct {
	Config    DisplayConfig
	Callbacks Callbacks
}

func NewEditor(conf EditorConf) (E, error) {
	teardown, err := enableLogs()
	if err != nil {
		panic(err)
	}
	defer teardown()

	defer func() {
		SwitchBackFromAlternateScreen(os.Stdout)

		os.Stdout.WriteString(ClearScreenCode)
		os.Stdout.WriteString(RepositionCursorCode)
		if err := recover(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %+v\n", err)
			fmt.Fprintf(os.Stderr, "stack: %s\n", debug.Stack())
			os.Exit(1)
		}
	}()

	// Set the terminal to raw mode
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}

	// Restore the old terminal settings when we finish
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	e := E{}
	e.setWindowSize()
	e.cfg = conf.Config

	e.keys = make(chan Key)
	go Parsekeys(os.Stdin, e.keys, e.Errs)

	e.signals = make(chan os.Signal)
	signal.Notify(e.signals, syscall.SIGWINCH)

	go func() {
		for err := range e.Errs {
			e.SetStatusLine("err: " + err.Error())
		}
	}()

	e.FullRender()

	return e, nil
}

func (e *E) SetStatusLine(format string, a ...interface{}) {
	e.statusMsg = fmt.Sprintf(format, a...)
}

func (e *E) detectSyntax() {
	e.syntax = nil
	if len(e.filename) == 0 {
		return
	}

	ext := filepath.Ext(e.filename)
	if len(ext) != 0 {
		// ext[1:] don't include the leading period
		e.syntax = e.filetypeLookup(ext[1:])
	}
}

func Parsekeys(r io.Reader, keys chan<- Key, errs chan<- error) {
	reader := bufio.NewReaderSize(r, 16)
	buf := make([]rune, 0, 4)
	for {
		// TODO abstract this
		r, _, err := reader.ReadRune()
		if err != nil && err != io.EOF {
			errs <- err
			continue
		}

		if r == '\x1b' {
			buf = append(buf, r)
			continue
		}

		if len(buf) == 0 {
			keys <- Key(r)
			continue
		}

		// I'm pretty sure there aren't any escape sequences
		// with more than 4 runes
		if len(buf) == 4 {
			for _, d := range buf {
				keys <- Key(d)
			}

			keys <- Key(r)
			buf = buf[:0]
			continue
		}

		buf = append(buf, r)

		key, ok := escapeCodeToKey[string(buf)]
		if ok {
			keys <- key
			buf = buf[:0]
		}
	}
}

func (e *E) Keys() <-chan Key {
	return e.keys
}

func (e *E) Signals() <-chan os.Signal {
	return e.signals
}

func (e *E) setWindowSize() error {
	cols, rows, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}

	// make room for status-bar and message-bar
	e.screenRows = rows - 2
	e.screenCols = cols

	return nil
}

// Enables logging and returns teardown function
func enableLogs() (func() error, error) {
	f, err := os.OpenFile(LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o666)
	if err != nil {
		return nil, errors.Wrap(err, "open: "+LogFile)
	}

	log.SetOutput(f)
	log.Println("Begin logging")

	return f.Close, nil
}

func (e *E) drawStatusBar(w io.Writer) {
	w.Write(getColor(InvertedColor))

	filename := e.filename
	if len(filename) == 0 {
		filename = "[No Name]"
	}

	modifiedStatus := ""
	if e.modified {
		modifiedStatus = "(modified)"
	}

	lmsg := fmt.Sprintf("%.20s - %d lines %s", filename, len(e.rows), modifiedStatus)
	if runewidth.StringWidth(lmsg) > e.screenCols {
		lmsg = runewidth.Truncate(lmsg, e.screenCols, "...")
	}
	w.Write([]byte(lmsg))

	filetype := "no filetype"
	if e.syntax != nil {
		filetype = e.syntax.Filetype
	}
	rmsg := fmt.Sprintf("%s | %d/%d", filetype, e.cy+1, len(e.rows))

	// Add padding between the left and right message
	l := runewidth.StringWidth(lmsg)
	r := runewidth.StringWidth(rmsg)
	for i := 0; i < e.screenCols-l-r; i++ {
		w.Write([]byte{' '})
	}

	w.Write([]byte(rmsg))
	w.Write([]byte("\r\n"))
	w.Write(ClearFormatting)
}
