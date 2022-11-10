package config

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"codeberg.org/wlcsm/li/ansi"
	"codeberg.org/wlcsm/li/core"
)

type KeyMapName string

type KeyMap struct {
	Name    KeyMapName
	Handler func(e *core.E, k ansi.Key) (bool, error)
}

const (
	BasicMapName    KeyMapName = "Basic"
	InsertModeName  KeyMapName = "Insert"
	CommandModeName KeyMapName = "Command"
	PromptModeName  KeyMapName = "Prompt"
)

var (
	BasicMap       KeyMap
	InsertModeMap  KeyMap
	CommandModeMap KeyMap
)

// Must be init'ed here to prevent an import cycle, since these maps can have a
// function that will set the keymapping and hence refer to themselves
func init() {
	BasicMap = KeyMap{
		Name:    BasicMapName,
		Handler: basicHandler,
	}

	InsertModeMap = KeyMap{
		Name:    InsertModeName,
		Handler: insertModeHandler,
	}

	CommandModeMap = KeyMap{
		Name:    CommandModeName,
		Handler: commandModeHandler,
	}
}

func basicHandler(e *core.E, k ansi.Key) (bool, error) {
	switch k {
	case ansi.UpArrowKey:
		e.SetY(e.Y() - 1)
	case ansi.DownArrowKey:
		e.SetY(e.Y() + 1)
	case ansi.LeftArrowKey:
		e.SetX(e.X() - 1)
	case ansi.RightArrowKey:
		e.SetX(e.X() + 1)
	case ansi.Ctrl('q'):
		return true, core.ErrQuitEditor
	case ansi.Ctrl('s'):
		log.Printf("attempting to save: %s\n", e.Filename())
		if err := e.Save(); err != nil {
			return true, err
		}

		e.SetStatusLine("saved file: %s", e.Filename())

		//	case ansi.Ctrl('w'):
		//		e.SetRow(e.Y(), append(e.rows[e.Y()].chars[:e.BackWord()], e.rows[e.Y()].chars[e.X()-1:]...))
	case ansi.Ctrl('u'):
		e.SetY(e.Y() - (e.ScreenRows() / 2))
	case ansi.Ctrl('d'):
		e.SetY(e.Y() + (e.ScreenRows() / 2))
	case ansi.DeleteKey, ansi.BackspaceKey:
		if row := e.Row(e.Y()); len(row) != 0 {
			e.SetRow(e.Y(), row[:len(row)-1])
			e.SetX(e.X() - 1)
		}
	default:
		log.Println("insert char: ", rune(k))
		e.AppendChar(e.Y(), rune(k))
		e.SetX(e.X() + 1)
		return false, nil
	}

	return true, nil
}

func insertModeHandler(e *core.E, k ansi.Key) (bool, error) {
	switch k {
	case ansi.EnterKey:
		row := e.Row(e.Y())
		//row, row2 := row[:e.X()], row[e.X():]

		e.SetRow(e.Y(), row)
		panic("unimplemented")
		//e.InsertRow(e.Y()+1, row2)

		e.SetY(e.Y() + 1)
		e.SetX(0)

	case ansi.CarriageReturnKey:
		row := e.Row(e.Y())
		//row, row2 := row[:e.X()], row[e.X():]

		e.SetRow(e.Y(), row)
		panic("unimplemented")
		//e.InsertRow(e.Y()+1, row2)

		e.SetY(e.Y() + 1)
		e.SetX(0)

	case ansi.DeleteKey, ansi.BackspaceKey:
		panic("unimplemented")
		//		x, y := e.X(), e.Y()
		//		if x != 0 {
		//			e.SetRow(y, append(e.rows[y].chars[:x-1], e.rows[y].chars[x:]...))
		//			e.SetX(x - 1)
		//		} else {
		//			e.SetY(y - 1)
		//			e.SetX(len(e.Row(y - 1)))
		//
		//			e.SetRow(y-1, append(e.Row(y-1), e.Row(y)...))
		//			e.DeleteRows(y, y)
		//		}

	default:
		panic("unimplemented")
		if k == ansi.Key('\t') || core.IsPrintable(k) {
			//e.InsertChars(e.Y(), e.X(), rune(k))
			e.SetX(e.X() + 1)
		}
	}

	return true, nil
}

const (
	StartSelection = "start"
)

func commandModeHandler(e *core.E, k ansi.Key) (bool, error) {
	switch k {
	case ansi.Key('j'):
		e.SetY(e.Y() + 1)
	case ansi.Key('k'):
		e.SetY(e.Y() - 1)
	case ansi.Key('h'):
		e.SetX(e.X() - 1)
	case ansi.Key('l'):
		e.SetX(e.X() + 1)
	case ansi.Key('J'):
		e.SetY(e.NumRows() - 1)
	case ansi.Key('K'):
		e.SetY(0)
	case ansi.Key('H'), ansi.Key('0'):
		e.SetX(0)
	case ansi.Key('G'):
		e.SetY(e.NumRows())
	case ansi.Key('C'):
		e.SetRow(e.Y(), []rune{})
	case ansi.Key('e'):
		StaticPrompt(e, "File name: ", func(f string) error {
			if len(f) == 0 {
				return fmt.Errorf("No file name")
			}

			return e.OpenFile(f)
		}, FileCompletion)
	case ansi.Key('s'):
		StaticPrompt(e, "$ ", func(res string) error {
			if len(res) == 0 {
				return nil
			}

			c := strings.Split(res, " ")
			out, err := exec.Command(c[0], c[1:]...).Output()
			if err != nil {
				return err
			}

			e.SetStatusLine(string(out))

			return nil
		})
	default:
		return false, nil
	}

	return true, nil
}

type Line struct {
	File string
	Row  int
	Orig string
}

func CreateList(l string) []Line {
	if len(l) == 0 {
		return nil
	}

	lines := strings.Split(l, "\n")

	parsed := make([]Line, 0)

	for _, line := range lines {
		// Just ignore empty lines
		if len(line) != 0 {
			parsed = append(parsed, parseLine(line))
		}
	}

	return parsed
}

func parseLine(l string) Line {
	i := strings.Index(l, ":")
	if i == -1 {
		return Line{Orig: l}
	}

	j := strings.Index(l[i+1:], ":")
	if j == -1 {
		return Line{Orig: l}
	}

	row, err := strconv.Atoi(l[i+1 : i+j+1])
	if err != nil {
		return Line{Orig: l}
	}

	return Line{
		File: l[:i],
		Row:  row,
		Orig: l,
	}
}

type CompletionFunc func(a string) ([]CmplItem, error)

type CmplItem struct {
	Display string
	Real    string
}

func FileCompletion(a string) ([]CmplItem, error) {
	//.cyes this will break on windows, idc
	i := strings.LastIndex(a, "/")
	if i == -1 {
		i = 0
	} else {
		i++
	}

	fileBasename := a[:i]
	fileHead := a[i:]

	log.Printf("fileBase: %s", fileBasename)

	files, err := os.ReadDir("./" + fileBasename)
	if err != nil {
		return nil, err
	}

	var res []CmplItem
	for _, f := range files {
		log.Printf("fil: %s", f.Name())
		if !strings.HasPrefix(f.Name(), fileHead) {
			continue
		}

		if f.IsDir() {
			res = append(res, CmplItem{
				Display: f.Name() + "/",
				Real:    fileBasename + f.Name() + "/",
			})
		} else if f.Type().IsRegular() {
			res = append(res, CmplItem{
				Display: f.Name(),
				Real:    fileBasename + f.Name(),
			})
		}
	}

	return res, nil
}

// StaticPrompt is a "normal" prompt designed to only get input from the user.
// It you want things to happen when you press any key, then use Prompt
func StaticPrompt(e *core.E, prompt string, end func(string) error, comp ...CompletionFunc) {
	var input string
	var cachedComp []CmplItem
	var compIndex int

	e.Prompt(prompt, func(k ansi.Key) string {
		log.Printf("key is: %s", string(k))

		switch k {
		case ansi.EnterKey, ansi.CarriageReturnKey:
			if err := end(input); err != nil {
				e.Errs <- err
			}

			return input
		case ansi.EscapeKey, ansi.Ctrl('q'):
			return ""
		case ansi.BackspaceKey, ansi.DeleteKey:
			if len(input) > 0 {
				input = input[:len(input)-1]
			}
		case ansi.Key('\t'):
			if len(comp) == 0 {
				break
			}

			if len(cachedComp) == 0 {
				for _, c := range comp {
					res, err := c(input)
					if err != nil {
						break
					}

					cachedComp = append(cachedComp, res...)
				}

				log.Printf("completion options: %v", cachedComp)
			}

			if len(cachedComp) == 0 {
				break
			}

			compIndex = (compIndex + 1) % len(cachedComp)
			input = cachedComp[compIndex].Real

		default:
			if core.IsPrintable(k) {
				input += string(k)
				cachedComp = nil
				compIndex = 0
			}
		}

		return input
	})
}
