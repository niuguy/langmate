package app

type HotkeyOption struct {
	ID       string
	Title    string
	Commands []string
}

var supportedHotkeys = []HotkeyOption{
	{
		ID:       "cmd_ctrl_r",
		Title:    "Cmd+Ctrl+R",
		Commands: []string{"cmd", "ctrl", "r"},
	},
	{
		ID:       "cmd_shift_r",
		Title:    "Cmd+Shift+R",
		Commands: []string{"cmd", "shift", "r"},
	},
	{
		ID:       "cmd_alt_r",
		Title:    "Cmd+Option+R",
		Commands: []string{"cmd", "alt", "r"},
	},
}

func hotkeyByID(id string) HotkeyOption {
	for _, hotkey := range supportedHotkeys {
		if hotkey.ID == id {
			return hotkey
		}
	}
	return supportedHotkeys[0]
}

func isSupportedHotkey(id string) bool {
	for _, hotkey := range supportedHotkeys {
		if hotkey.ID == id {
			return true
		}
	}
	return false
}
