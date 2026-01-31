package app

import (
	"os/exec"
	"time"
)

// SimulateCopy sends Cmd+C to copy selected text to clipboard using AppleScript
func SimulateCopy() {
	script := `tell application "System Events" to keystroke "c" using command down`
	exec.Command("osascript", "-e", script).Run()
	time.Sleep(100 * time.Millisecond)
}

// SimulatePaste sends Cmd+V to paste clipboard content using AppleScript
func SimulatePaste() {
	script := `tell application "System Events" to keystroke "v" using command down`
	exec.Command("osascript", "-e", script).Run()
}
