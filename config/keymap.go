package config

import (
	"log"

	"codeberg.org/wlcsm/li/ansi"
	"codeberg.org/wlcsm/li/core"
	"github.com/pkg/errors"
)

type EditorMode int8

const (
	InsertMode EditorMode = iota + 1
	CommandMode
	PromptMode
)

// ProcessKey processes a key read from stdin.
// Returns errQuitEditor when user requests to quit.
func ProcessKey(e *core.E, k ansi.Key) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.Wrap(e.(error), "panicked")
		}
	}()

	log.Printf("processing key: %s", string(k))

	_, err = basicHandler(e, k)
	if err != nil {
		return err
	}

	return nil
}
