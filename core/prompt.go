package core

import "codeberg.org/wlcsm/li/ansi"

// Prompt shows the given prompt in the status bar and get user input
func (e *E) Prompt(prompt string, keymap func(k ansi.Key) string) {
	if keymap == nil {
		panic("can't give a nil function to prompt")
	}

	e.SetStatusLine(prompt)

	panic("not implemented")
	//	e.SetBackupKeymap()
	//	e.SetKeymap(func(e *E, k Key) error {
	//		if k == EnterKey {
	//			e.RestoreBackup()
	//			return nil
	//		}
	//
	//		s := keymap(k)
	//		e.SetStatusLine(prompt + s)
	//		return nil
	//	})
}
