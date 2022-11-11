package main

import (
	"fmt"
	"os"

	"codeberg.org/wlcsm/li/config"
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
		Config: core.DisplayConfig{
			Tabstop: 8,
		},
		Keymap:    config.ProcessKey,
	}

	return core.NewEditor(conf, os.Args)
}
