package style

import (
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-isatty"
)

var (
	Enabled bool

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

// DisableColor switches rendering to plain output (no color/no bold).
func DisableColor() {
	Enabled = false
}

// EnableColor switches rendering back to colored output.
func EnableColor() {
	Enabled = true
}

// Render helpers respect the Enabled flag. When disabled, they return
// the original text unmodified; when enabled, they apply the lipgloss style.
func RenderSuccess(s string) string {
	if !Enabled {
		return s
	}
	return Success.Render(s)
}

func RenderError(s string) string {
	if !Enabled {
		return s
	}
	return Error.Render(s)
}

func RenderWarning(s string) string {
	if !Enabled {
		return s
	}
	return Warning.Render(s)
}

func RenderInfo(s string) string {
	if !Enabled {
		return s
	}
	return Info.Render(s)
}

// init sets the default Enabled based on whether stdout is a terminal.
func init() {
	Enabled = isatty.IsTerminal(os.Stdout.Fd())
}
