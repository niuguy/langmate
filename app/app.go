package app

import (
	"fmt"
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

func StartHook(textProcessor llm.TextProcessor) {
	// p := tea.NewProgram(initialModel())

	fmt.Println("Pending for clipboard content...")
	// go func() {
	hook.Register(hook.KeyDown, []string{"cmd", "c"}, func(e hook.Event) {
		currentTime := time.Now()
		content, err := clipboard.ReadAll()
		if err != nil {
			fmt.Println("Error reading clipboard:", err)
			return
		}
		if content == lastContent && currentTime.Sub(lastCopyTime) < doubleCopyThreshold {
			fmt.Println(content)
			fmt.Println("-----------------")
			fmt.Println(textProcessor.TransferText(content, "en"))
			fmt.Println("-----------------")
			fmt.Println("Pending for clipboard content...")
		}
		lastCopyTime = currentTime
		lastContent = content
	})
	s := hook.Start()
	<-hook.Process(s)
	// }()
	// if _, err := p.Run(); err != nil {
	// 	fmt.Println("Error running program:", err)
	// }
}
