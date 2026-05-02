package app

import (
	_ "embed"
	"fmt"
	"os"

	"fyne.io/systray"
	"github.com/niuguy/langmate/llm"
)

//go:embed icon.png
var iconData []byte

var menuState *menuBarState

type menuBarState struct {
	status        *systray.MenuItem
	preview       *systray.MenuItem
	modelItems    map[string]*systray.MenuItem
	hotkeyItems   map[string]*systray.MenuItem
	currentModel  *systray.MenuItem
	currentHotkey *systray.MenuItem
}

// SetMenuBarStatus updates the menu bar title to show current status
func SetMenuBarStatus(status string) {
	if status == "" {
		systray.SetTitle("LM")
	} else {
		systray.SetTitle(status)
	}
}

// RunMenuBar starts the menu bar icon and blocks until quit
func RunMenuBar(state *daemonState, onReady func(), onQuit func()) {
	systray.Run(func() {
		fmt.Printf("Menu bar initializing, icon size: %d bytes\n", len(iconData))
		systray.SetTemplateIcon(iconData, iconData)
		systray.SetTitle("LM")
		updateMenuBarTooltip(state)

		mStatus := systray.AddMenuItem("LangMate Running", "")
		mStatus.Disable()

		systray.AddSeparator()

		mPreview := systray.AddMenuItemCheckbox("Show Preview", "Preview and edit rephrased text before replacing the selection", false)

		mModel := systray.AddMenuItem("Model", "Choose model provider")
		modelItems := map[string]*systray.MenuItem{}
		for _, preset := range llm.ModelPresets {
			modelItems[preset.ID] = mModel.AddSubMenuItemCheckbox(preset.Title, modelTooltip(preset), false)
		}

		mHotkey := systray.AddMenuItem("Hotkey", "Choose the rephrase hotkey")
		hotkeyItems := map[string]*systray.MenuItem{}
		for _, hotkey := range supportedHotkeys {
			hotkeyItems[hotkey.ID] = mHotkey.AddSubMenuItemCheckbox(hotkey.Title, "Set rephrase hotkey to "+hotkey.Title, false)
		}

		mCurrentModel := systray.AddMenuItem("", "")
		mCurrentModel.Disable()
		mCurrentHotkey := systray.AddMenuItem("", "")
		mCurrentHotkey.Disable()

		systray.AddSeparator()

		mQuit := systray.AddMenuItem("Quit", "Quit LangMate")

		menuState = &menuBarState{
			status:        mStatus,
			preview:       mPreview,
			modelItems:    modelItems,
			hotkeyItems:   hotkeyItems,
			currentModel:  mCurrentModel,
			currentHotkey: mCurrentHotkey,
		}
		updateMenuBarState(state)
		startMenuHandlers(state, mQuit)

		// Call the ready callback
		if onReady != nil {
			go onReady()
		}

	}, func() {
		if onQuit != nil {
			onQuit()
		}
		os.Exit(0)
	})
}

func startMenuHandlers(state *daemonState, quitItem *systray.MenuItem) {
	go func() {
		for range menuState.preview.ClickedCh {
			config, _ := state.snapshot()
			state.setPreview(!config.Preview)
			updateMenuBarState(state)
		}
	}()

	for model, item := range menuState.modelItems {
		model := model
		item := item
		go func() {
			for range item.ClickedCh {
				if err := state.setModel(model); err != nil {
					SetMenuBarStatus("Model error!")
					fmt.Printf("Model switch error: %v\n", err)
					clearStatusAfterDelay()
					continue
				}
				updateMenuBarState(state)
				SetMenuBarStatus("Model: " + modelTitle(model))
				clearStatusAfterDelay()
			}
		}()
	}

	for hotkeyID, item := range menuState.hotkeyItems {
		hotkeyID := hotkeyID
		item := item
		go func() {
			for range item.ClickedCh {
				state.setHotkey(hotkeyID)
				SetMenuBarStatus("Hotkey: " + hotkeyByID(hotkeyID).Title)
				clearStatusAfterDelay()
			}
		}()
	}

	go func() {
		<-quitItem.ClickedCh
		systray.Quit()
	}()
}

func updateMenuBarState(state *daemonState) {
	if menuState == nil {
		return
	}

	config, _ := state.snapshot()
	if config.Preview {
		menuState.preview.Check()
	} else {
		menuState.preview.Uncheck()
	}

	for model, item := range menuState.modelItems {
		if config.Model == model {
			item.Check()
		} else {
			item.Uncheck()
		}
	}

	for hotkeyID, item := range menuState.hotkeyItems {
		if config.Hotkey == hotkeyID {
			item.Check()
		} else {
			item.Uncheck()
		}
	}

	menuState.currentModel.SetTitle("Current Model: " + modelTitle(config.Model))
	menuState.currentHotkey.SetTitle("Current Hotkey: " + hotkeyByID(config.Hotkey).Title)
}

func updateMenuBarTooltip(state *daemonState) {
	config, _ := state.snapshot()
	systray.SetTooltip("LangMate - " + hotkeyByID(config.Hotkey).Title + " to rephrase")
}

func modelTitle(model string) string {
	preset, ok := llm.FindModelPreset(llm.NormalizeModelPresetID(model))
	if !ok {
		return model
	}
	return preset.Title
}

func modelTooltip(preset llm.ModelPreset) string {
	switch preset.Provider {
	case llm.ProviderOpenAI:
		return "Use " + preset.Model + " with OPENAI_API_KEY"
	case llm.ProviderOllama:
		return "Use local Ollama model " + preset.Model
	default:
		return preset.Model
	}
}
