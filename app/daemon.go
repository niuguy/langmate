package app

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/atotto/clipboard"
	"github.com/niuguy/langmate/llm"
	hook "github.com/robotn/gohook"
)

const (
	debounceInterval = 500 * time.Millisecond
	llmTimeout       = 30 * time.Second
)

var (
	lastTriggerTime time.Time
	debounceMu      sync.Mutex
)

// StartDaemon runs the background daemon that listens for Cmd+Shift+R
func StartDaemon(textProcessor llm.TextProcessor, lang string) {
	fmt.Println("LangMate daemon started")
	fmt.Println("Press Cmd+Ctrl+R to rephrase selected text")
	fmt.Println("Press Ctrl+C to quit")
	fmt.Println()

	// Run menu bar with hotkey listener
	RunMenuBar(func() {
		// onReady: start the hotkey listener
		NotifyStartup()

		hook.Register(hook.KeyDown, []string{"cmd", "ctrl", "r"}, func(e hook.Event) {
			go handleRephraseHotkey(textProcessor, lang)
		})

		s := hook.Start()
		defer hook.End()
		<-hook.Process(s)
	}, nil)
}

func handleRephraseHotkey(textProcessor llm.TextProcessor, lang string) {
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

	// Simulate Cmd+C to copy selected text
	SimulateCopy()

	// Read the copied text
	selectedText, err := clipboard.ReadAll()
	if err != nil {
		SetMenuBarStatus("Error!")
		go func() {
			time.Sleep(2 * time.Second)
			SetMenuBarStatus("")
		}()
		return
	}

	// Validate selection
	selectedText = strings.TrimSpace(selectedText)
	if selectedText == "" {
		SetMenuBarStatus("No text!")
		restoreClipboard(originalClipboard)
		go func() {
			time.Sleep(2 * time.Second)
			SetMenuBarStatus("")
		}()
		return
	}

	// Show processing status in menu bar
	SetMenuBarStatus("Rephrasing...")
	fmt.Printf("Processing: %s\n", truncateText(selectedText, 50))

	// Call LLM with timeout
	ctx, cancel := context.WithTimeout(context.Background(), llmTimeout)
	defer cancel()

	resultChan := make(chan string, 1)
	errChan := make(chan error, 1)

	go func() {
		result, err := textProcessor.RephraseText(selectedText, lang)
		if err != nil {
			errChan <- err
		} else {
			resultChan <- result
		}
	}()

	var rephrasedText string
	select {
	case <-ctx.Done():
		SetMenuBarStatus("Timeout!")
		restoreClipboard(originalClipboard)
		go func() {
			time.Sleep(2 * time.Second)
			SetMenuBarStatus("")
		}()
		return
	case err := <-errChan:
		SetMenuBarStatus("Error!")
		fmt.Printf("LLM error: %v\n", err)
		restoreClipboard(originalClipboard)
		go func() {
			time.Sleep(2 * time.Second)
			SetMenuBarStatus("")
		}()
		return
	case rephrasedText = <-resultChan:
	}

	// Write rephrased text to clipboard
	if err := clipboard.WriteAll(rephrasedText); err != nil {
		SetMenuBarStatus("Error!")
		restoreClipboard(originalClipboard)
		go func() {
			time.Sleep(2 * time.Second)
			SetMenuBarStatus("")
		}()
		return
	}

	// Wait a bit for clipboard to settle
	time.Sleep(50 * time.Millisecond)

	// Simulate Cmd+V to paste
	SimulatePaste()

	// Show success status briefly then clear
	SetMenuBarStatus("Done!")
	fmt.Printf("Rephrased: %s\n", truncateText(rephrasedText, 50))
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
