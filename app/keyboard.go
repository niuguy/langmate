package app

import (
	"fmt"
	"os/exec"
	"time"
)

// SimulateCopy sends Cmd+C to copy selected text to clipboard using AppleScript
func SimulateCopy() {
	script := `tell application "System Events" to keystroke "c" using command down`
	cmd := exec.Command("osascript", "-e", script)
	if err := cmd.Run(); err != nil {
		fmt.Printf("SimulateCopy error: %v\n", err)
	}
	time.Sleep(200 * time.Millisecond)
}

// SimulatePaste sends Cmd+V to paste clipboard content using AppleScript
func SimulatePaste() {
	script := `tell application "System Events" to keystroke "v" using command down`
	cmd := exec.Command("osascript", "-e", script)
	if err := cmd.Run(); err != nil {
		fmt.Printf("SimulatePaste error: %v\n", err)
	}
}
