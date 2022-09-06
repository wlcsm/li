package sdk

import "codeberg.org/wlcsm/li/core"

// Prompt shows the given prompt in the status bar and get user input
func (s *SDK) Prompt(prompt string, keymap func(k core.Key) string) {
	if keymap == nil {
		panic("can't give a nil function to prompt")
	}

	s.e.SetStatusLine(prompt)

	s.SetBackupKeymap()
	s.SetKeymap(func(e *core.E, k core.Key) error {
		if k == core.EnterKey {
			s.RestoreBackup()
			return nil
		}

		s := keymap(k)
		e.SetStatusLine(prompt + s)
		return nil
	})
}
