package app

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/atotto/clipboard"
	"github.com/niuguy/langmate/llm"
	hook "github.com/robotn/gohook"
)

const (
	doubleCopyThreshold = 500 * time.Millisecond
)

var (
	lastCopyTime = time.Now()
	lastContent  = ""
)

func StartHook(textProcessor llm.TextProcessor, lang string) {
	// p := tea.NewProgram(initialModel())

	fmt.Println("Type or double cmd+c to translate or rephrase...")
	fmt.Println()

	go func() {
		hook.Register(hook.KeyDown, []string{"cmd", "c"}, func(e hook.Event) {
			currentTime := time.Now()
			content, err := clipboard.ReadAll()
			if err != nil {
				fmt.Println("Error reading clipboard:", err)
				return
			}
			if content == lastContent && currentTime.Sub(lastCopyTime) < doubleCopyThreshold {

				done := make(chan bool)
				go showWaitingAnimation(done)
				processedText, _ := textProcessor.TransferText(content, lang)
				done <- true
				fmt.Print("\033[H\033[2J") // Clear screen
				fmt.Println(content)
				fmt.Println("-----------------")
				fmt.Println()
				fmt.Println(processedText)
			}
			lastCopyTime = currentTime
			lastContent = content
		})
		s := hook.Start()
		defer hook.End()
		<-hook.Process(s)
	}()

	// Handle manual text input
	reader := bufio.NewReader(os.Stdin)
	for {

		input, _ := reader.ReadString('\n')
		if input == "exit\n" {
			break
		}

		done := make(chan bool)
		go showWaitingAnimation(done)
		processedText, err := textProcessor.TransferText(input, lang)
		done <- true
		fmt.Print("\033[H\033[2J") // Clear screen
		fmt.Println(input)

		fmt.Println("-----------------")
		fmt.Println()
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println(processedText)
		}
	}
}

func showWaitingAnimation(done chan bool) {
	for {
		select {
		case <-done:
			return
		default:
			fmt.Print(".")
			time.Sleep(500 * time.Millisecond) // Adjust timing as needed
		}
	}
}
