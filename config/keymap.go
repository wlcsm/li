package config

import (
	"log"

	"github.com/pkg/errors"
	"codeberg.org/wlcsm/li/core"
)

type EditorMode int8

const (
	InsertMode EditorMode = iota + 1
	CommandMode
	PromptMode
)

// ProcessKey processes a key read from stdin.
// Returns errQuitEditor when user requests to quit.
func ProcessKey(e *core.E, k core.Key) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.Wrap(e.(error), "panicked")
		}
	}()


	for _, keymap := range e.keymapping {
		log.Printf("processing key: %s, with keymap: %s", string(k), keymap.Name)

		handled, err := keymap.Handler(e, k)
		if err != nil {
			return err
		}

		if handled {
			return nil
		}
	}

	return nil
}
