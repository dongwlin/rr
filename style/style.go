package style

import "github.com/charmbracelet/lipgloss"

var (
	Success = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FF00")).
		Bold(true)

	Error = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF0000")).
		Bold(true)

	Warning = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFF00")).
		Bold(true)

	Info = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FFFF")).
		Bold(true)
)
