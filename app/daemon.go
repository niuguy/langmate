package app

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/atotto/clipboard"
	"github.com/niuguy/langmate/llm"
	hook "github.com/robotn/gohook"
)

const (
	debounceInterval = 1000 * time.Millisecond
	llmTimeout       = 30 * time.Second
)

var (
	lastTriggerTime time.Time
	debounceMu      sync.Mutex
)

type daemonState struct {
	mu            sync.RWMutex
	config        DaemonConfig
	textProcessor llm.TextProcessor
}

// StartDaemon runs the background daemon that listens for Cmd+Ctrl+R
func StartDaemon(config DaemonConfig) {
	config.normalize()

	textProcessor, err := llm.CreateTextProcessor(config.Model)
	if err != nil {
		fmt.Printf("Model setup error: %v\n", err)
	}

	state := &daemonState{
		config:        config,
		textProcessor: textProcessor,
	}

	fmt.Println("LangMate daemon started")
	fmt.Printf("Press %s to rephrase selected text\n", hotkeyByID(config.Hotkey).Title)
	fmt.Println("Press Ctrl+C to quit")
	if config.Preview {
		fmt.Println("Preview mode enabled")
	}
	fmt.Println()

	// Run menu bar with hotkey listener
	RunMenuBar(state, func() {
		// onReady: start the hotkey listener
		config, _ := state.snapshot()
		NotifyStartup(hotkeyByID(config.Hotkey).Title)

		for _, hotkey := range supportedHotkeys {
			hotkey := hotkey
			hook.Register(hook.KeyDown, hotkey.Commands, func(e hook.Event) {
				if !state.isActiveHotkey(hotkey.ID) {
					return
				}
				fmt.Println("Hotkey triggered")
				go handleRephraseHotkey(state)
			})
		}

		s := hook.Start()
		defer hook.End()
		<-hook.Process(s)
	}, nil)
}

func handleRephraseHotkey(state *daemonState) {
	// Debounce to prevent double-triggers
	debounceMu.Lock()
	if time.Since(lastTriggerTime) < debounceInterval {
		debounceMu.Unlock()
		return
	}
	lastTriggerTime = time.Now()
	debounceMu.Unlock()

	// Save original clipboard content for restoration on error
	originalClipboard, _ := clipboard.ReadAll()
	fmt.Printf("Original clipboard: %s\n", truncateText(originalClipboard, 30))

	selectionMarker := fmt.Sprintf("__LANGMATE_SELECTION_MARKER_%d__", time.Now().UnixNano())
	if err := clipboard.WriteAll(selectionMarker); err != nil {
		SetMenuBarStatus("Clipboard!")
		clearStatusAfterDelay()
		return
	}

	// Simulate Cmd+C to copy selected text
	SimulateCopy()

	// Read the copied text
	selectedText, err := clipboard.ReadAll()
	if err != nil {
		SetMenuBarStatus("Error!")
		restoreClipboard(originalClipboard)
		clearStatusAfterDelay()
		return
	}

	fmt.Printf("After copy clipboard: %s\n", truncateText(selectedText, 30))

	// Validate selection - check if clipboard changed and has content
	selectedText = strings.TrimSpace(selectedText)
	if selectedText == "" || selectedText == selectionMarker {
		SetMenuBarStatus("No text!")
		restoreClipboard(originalClipboard)
		clearStatusAfterDelay()
		return
	}

	config, textProcessor := state.snapshot()
	if textProcessor == nil {
		SetMenuBarStatus("Model error!")
		restoreClipboard(originalClipboard)
		clearStatusAfterDelay()
		return
	}

	// Show processing status in menu bar
	SetMenuBarStatus("Rephrasing...")
	fmt.Printf("Processing: %s\n", truncateText(selectedText, 50))

	// Call LLM with timeout
	ctx, cancel := context.WithTimeout(context.Background(), llmTimeout)
	defer cancel()
	rephrasedText, err := textProcessor.RephraseText(ctx, selectedText, config.Lang)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded), errors.Is(err, context.Canceled), ctx.Err() != nil:
			SetMenuBarStatus("Timeout!")
		default:
			SetMenuBarStatus("Error!")
			fmt.Printf("LLM error: %v\n", err)
		}
		restoreClipboard(originalClipboard)
		clearStatusAfterDelay()
		return
	}

	if config.Preview {
		SetMenuBarStatus("Previewing...")
		result, err := ShowPreview(selectedText, rephrasedText, modelTitle(config.Model))
		if err != nil {
			SetMenuBarStatus("Preview error!")
			fmt.Printf("Preview error: %v\n", err)
			restoreClipboard(originalClipboard)
			clearStatusAfterDelay()
			return
		}

		switch result.Action {
		case PreviewActionReplace:
			rephrasedText = result.Text
		case PreviewActionCopy:
			if err := clipboard.WriteAll(result.Text); err != nil {
				SetMenuBarStatus("Error!")
			} else {
				SetMenuBarStatus("Copied!")
			}
			clearStatusAfterDelay()
			return
		default:
			restoreClipboard(originalClipboard)
			SetMenuBarStatus("Canceled")
			clearStatusAfterDelay()
			return
		}
	}

	// Write rephrased text to clipboard
	if err := clipboard.WriteAll(rephrasedText); err != nil {
		SetMenuBarStatus("Error!")
		restoreClipboard(originalClipboard)
		clearStatusAfterDelay()
		return
	}

	// Wait a bit for clipboard to settle
	time.Sleep(50 * time.Millisecond)

	// Simulate Cmd+V to paste
	SimulatePaste()
	time.Sleep(150 * time.Millisecond)
	restoreClipboard(originalClipboard)

	// Show success status briefly then clear
	SetMenuBarStatus("Done!")
	fmt.Printf("Rephrased: %s\n", truncateText(rephrasedText, 50))
	clearStatusAfterDelay()
}

func (s *daemonState) snapshot() (DaemonConfig, llm.TextProcessor) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.config, s.textProcessor
}

func (s *daemonState) isActiveHotkey(id string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.config.Hotkey == id
}

func (s *daemonState) setPreview(preview bool) {
	s.mu.Lock()
	s.config.Preview = preview
	config := s.config
	s.mu.Unlock()

	saveConfigBestEffort(config)
}

func (s *daemonState) setHotkey(id string) {
	if !isSupportedHotkey(id) {
		return
	}

	s.mu.Lock()
	s.config.Hotkey = id
	config := s.config
	s.mu.Unlock()

	saveConfigBestEffort(config)
	updateMenuBarState(s)
	updateMenuBarTooltip(s)
}

func (s *daemonState) setModel(model string) error {
	textProcessor, err := llm.CreateTextProcessor(model)
	if err != nil {
		return err
	}

	s.mu.Lock()
	s.config.Model = model
	s.textProcessor = textProcessor
	config := s.config
	s.mu.Unlock()

	saveConfigBestEffort(config)
	return nil
}

func saveConfigBestEffort(config DaemonConfig) {
	if err := SaveDaemonConfig(config); err != nil {
		fmt.Printf("Config save error: %v\n", err)
	}
}

func clearStatusAfterDelay() {
	go func() {
		time.Sleep(2 * time.Second)
		SetMenuBarStatus("")
	}()
}

func restoreClipboard(content string) {
	clipboard.WriteAll(content)
}

func truncateText(text string, maxLen int) string {
	text = strings.ReplaceAll(text, "\n", " ")
	if len(text) > maxLen {
		return text[:maxLen] + "..."
	}
	return text
}
