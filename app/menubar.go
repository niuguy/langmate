package app

import (
	_ "embed"
	"fmt"
	"os"

	"fyne.io/systray"
)

//go:embed icon.png
var iconData []byte

// SetMenuBarStatus updates the menu bar title to show current status
func SetMenuBarStatus(status string) {
	if status == "" {
		systray.SetTitle("LM")
	} else {
		systray.SetTitle(status)
	}
}

// RunMenuBar starts the menu bar icon and blocks until quit
func RunMenuBar(onReady func(), onQuit func()) {
	systray.Run(func() {
		fmt.Printf("Menu bar initializing, icon size: %d bytes\n", len(iconData))
		systray.SetTemplateIcon(iconData, iconData)
		systray.SetTitle("LM")
		systray.SetTooltip("LangMate - Cmd+Ctrl+R to rephrase")

		mStatus := systray.AddMenuItem("LangMate Running", "")
		mStatus.Disable()

		systray.AddSeparator()

		mQuit := systray.AddMenuItem("Quit", "Quit LangMate")

		// Call the ready callback
		if onReady != nil {
			go onReady()
		}

		// Wait for quit
		<-mQuit.ClickedCh
		systray.Quit()
	}, func() {
		if onQuit != nil {
			onQuit()
		}
		os.Exit(0)
	})
}
