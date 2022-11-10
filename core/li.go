package core

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"codeberg.org/wlcsm/li/ansi"
	"github.com/mattn/go-runewidth"
	"github.com/pkg/errors"
	"golang.org/x/term"
)

var (
	ErrPromptCanceled = errors.New("user canceled the input prompt")
	ErrQuitEditor     = errors.New("quit editor")
)

var (
	LogFile = filepath.Join(os.Getenv("HOME"), ".li.log")
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

	keymap func(*E, ansi.Key) error

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
	Keymap    func(*E, ansi.Key) error
	Callbacks Callbacks
}

func NewEditor(conf EditorConf, args []string) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.Wrap(e.(error), "panic")
		}
	}()

	logFile, err := os.OpenFile(LogFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o666)
	if err != nil {
		return errors.Wrap(err, "opening log file: "+LogFile)
	}
	defer logFile.Close()

	log.SetOutput(logFile)
	log.Println("Begin logging")

	SwitchToAlternateScreen(os.Stdout)
	defer SwitchBackFromAlternateScreen(os.Stdout)

	// Set the terminal to raw mode. This allows us to directly receive the
	// user's raw input without further processing by the terminal
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	e := E{}
	e.setWindowSize()
	e.cfg = conf.Config
	e.keymap = conf.Keymap

	if len(args) > 1 {
		err := e.OpenFile(args[1])
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}
	} else {
		e.rows = []*Row{{}}
	}

	e.FullRender()

	e.signals = make(chan os.Signal)
	signal.Notify(e.signals, syscall.SIGWINCH)

	go func() {
		d := ansi.NewDecoder(os.Stdin)
		for {
			key, err := d.Decode()
			if err != nil {
				e.Errs <- err
			}

			if err := e.keymap(&e, key); err != nil {
				e.Errs <- err
			}
		}
	}()

	e.Errs = make(chan error)
	for {
		select {
		case err := <-e.Errs:
			if err == ErrQuitEditor {
				return nil
			}
			e.SetStatusLine("err: " + err.Error())
		}
	}
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

func (e *E) ScreenRows() int {
	return e.screenRows
}

func (e *E) ScreenCols() int {
	return e.screenCols
}

func (e *E) ScreenCenter() int {
	return e.rowOffset + (e.screenRows / 2)
}
