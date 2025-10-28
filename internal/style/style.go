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

// renderWith is a small helper that applies the provided lipgloss style to
// the input string only when color rendering is enabled. When disabled it
// returns the original, unmodified string. Centralizing this check avoids
// repeating the Enabled guard in each exported RenderX function.
func renderWith(style lipgloss.Style, s string) string {
	if !Enabled {
		return s
	}
	return style.Render(s)
}

// RenderSuccess returns the string rendered with the Success style when
// coloring is enabled; otherwise it returns the plain string. Use this in
// callers that want a standardized success label (e.g. "success").
func RenderSuccess(s string) string {
	return renderWith(Success, s)
}

// RenderError returns the string rendered with the Error style when
// coloring is enabled; otherwise it returns the plain string.
func RenderError(s string) string {
	return renderWith(Error, s)
}

// RenderWarning returns the string rendered with the Warning style when
// coloring is enabled; otherwise it returns the plain string.
func RenderWarning(s string) string {
	return renderWith(Warning, s)
}

// RenderInfo returns the string rendered with the Info style when
// coloring is enabled; otherwise it returns the plain string.
func RenderInfo(s string) string {
	return renderWith(Info, s)
}

// init sets the default Enabled based on whether stdout is a terminal.
func init() {
	Enabled = isatty.IsTerminal(os.Stdout.Fd())
}
