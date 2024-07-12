package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	viewportWidth  = 60
	viewportHeight = 10
)

type model struct {
	topViewport    viewport.Model
	bottomViewport viewport.Model
	client         *OpenAIClient
}

func initialModel() model {
	topVp := viewport.New(viewportWidth, viewportHeight)
	topVp.SetContent("Waiting for clipboard content...")

	bottomVp := viewport.New(viewportWidth, viewportHeight)
	bottomVp.SetContent("Transferred text will appear here...")

	return model{
		topViewport:    topVp,
		bottomViewport: bottomVp,
		client:         NewOpenAIClient(),
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
		m.topViewport.SetContent(msg.content)
		transferredText, err := m.client.transferText(msg.content, "en")
		if err != nil {
			m.bottomViewport.SetContent(fmt.Sprintf("Error: %v", err))
		} else {
			m.bottomViewport.SetContent(transferredText)
		}
	}
	m.topViewport, cmd = m.topViewport.Update(msg)
	m.bottomViewport, _ = m.bottomViewport.Update(msg)
	return m, cmd
}

func (m model) View() string {
	viewportStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62"))
		// Padding(1).
		// Width(viewportWidth + 2) // +2 for borders

	return lipgloss.JoinVertical(
		lipgloss.Center,
		viewportStyle.Render(m.topViewport.View()),
		viewportStyle.Render(m.bottomViewport.View()),
	) + "\n\nPress 'q' to quit"
}

type clipboardMsg struct {
	content string
}
