package app

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/niuguy/langmate/llm"
)

const (
	viewportWidth      = 80 // Adjusted for two wider viewports side by side
	viewportUpHeight   = 10
	viewportDownHeight = 15
)

type model struct {
	upViewport   viewport.Model
	downViewport viewport.Model
	client       *llm.OpenAIClient
}

func initialModel() model {
	leftVp := viewport.New(viewportWidth, viewportUpHeight)
	leftVp.SetContent("Waiting for clipboard content...")

	rightVp := viewport.New(viewportWidth, viewportDownHeight)
	rightVp.SetContent("Transferred text will appear here...")

	return model{
		upViewport:   leftVp,
		downViewport: rightVp,
		client:       llm.NewOpenAIClient(),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case clipboardMsg:
		m.upViewport.SetContent(msg.content)
		transferredText, err := m.client.TransferText(msg.content, "en")
		if err != nil {
			m.downViewport.SetContent(fmt.Sprintf("Error: %v", err))
		} else {
			m.downViewport.SetContent(transferredText)
		}
	}
	m.upViewport, cmd = m.upViewport.Update(msg)
	m.downViewport, _ = m.downViewport.Update(msg)
	return m, cmd
}

func (m model) View() string {
	viewportStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("240")).Padding(1)

	// divider := lipgloss.NewStyle().
	// 	Background(lipgloss.Color("242")).
	// 	Width(1).
	// 	Render(" ")

	return lipgloss.JoinVertical(
		lipgloss.Top,
		viewportStyle.Render(m.upViewport.View()),
		// divider,
		viewportStyle.Render(m.downViewport.View()),
	) + "\n\nPress 'q' to quit"
}

type clipboardMsg struct {
	content string
}
