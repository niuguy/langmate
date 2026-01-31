package app

import (
	gosxnotifier "github.com/deckarep/gosx-notifier"
)

const appName = "LangMate"

// NotifyStartup shows a notification that the daemon has started
func NotifyStartup() {
	note := gosxnotifier.NewNotification("Ready! Press Cmd+Shift+R to rephrase selected text.")
	note.Title = appName
	note.Sound = gosxnotifier.Default
	note.Push()
}

// NotifyProcessing shows a notification that text is being processed
func NotifyProcessing() {
	note := gosxnotifier.NewNotification("Processing your text...")
	note.Title = appName
	note.Sound = gosxnotifier.Default
	note.Push()
}

// NotifySuccess shows a notification that text was successfully rephrased
func NotifySuccess() {
	note := gosxnotifier.NewNotification("Text rephrased and pasted!")
	note.Title = appName
	note.Sound = gosxnotifier.Default
	note.Push()
}

// NotifyError shows an error notification with the given message
func NotifyError(msg string) {
	note := gosxnotifier.NewNotification(msg)
	note.Title = appName
	note.Sound = gosxnotifier.Basso
	note.Push()
}
