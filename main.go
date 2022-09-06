package main

import (
	"errors"
	"fmt"
	"os"

	"codeberg.org/wlcsm/li/core"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("err: %+v\n", err)
		os.Exit(1)
	}
}

func run() error {
	conf := core.EditorConf{
		Config:    core.DisplayConfig{},
		Callbacks: core.Callbacks{},
	}

	e, err := core.NewEditor(conf)
	if err != nil {
		return err
	}

	if len(os.Args) > 1 {
		err := e.OpenFile(os.Args[1])
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			panic(err)
		}
	}

	return nil
}
