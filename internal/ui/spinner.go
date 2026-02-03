package ui

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/charmbracelet/lipgloss"
)

var s *spinner.Spinner

func init() {
	s = spinner.New(spinner.CharSets[14], 80*time.Millisecond)
	s.Color("cyan")
}

func StartSpinner(msg string) {
	s.Suffix = " " + msg
	s.Start()
}

func StopSpinner() {
	s.Stop()
}

func StartThinkingSpinner(msg string) {
	StartSpinner(msg)
}

var (
	gray  = lipgloss.Color("#6B7280")
	white  = lipgloss.Color("#E5E7EB")
	green  = lipgloss.Color("#10B981")
	red    = lipgloss.Color("#EF4444")

	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(gray).
			Padding(1, 2)

	LabelStyle = lipgloss.NewStyle().
			Foreground(gray)

	ValueStyle = lipgloss.NewStyle().
			Foreground(white)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(green)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(red)

	MutedStyle = lipgloss.NewStyle().
			Foreground(gray)
)

func RenderAnalysisResult(title, description string) string {
	titleLine := LabelStyle.Render("Título: ") + ValueStyle.Bold(true).Render(title)
	descLine := LabelStyle.Render("Descrição:")
	descContent := ValueStyle.Width(78).Render(description)

	content := lipgloss.JoinVertical(lipgloss.Left,
		titleLine,
		"",
		descLine,
		descContent,
	)

	return BoxStyle.Render(content)
}

func PrintSuccess(msg string) {
	fmt.Println(SuccessStyle.Render("[OK] " + msg))
}

func PrintError(msg string) {
	fmt.Println(ErrorStyle.Render("[ERRO] " + msg))
}

func PrintInfo(msg string) {
	fmt.Println(MutedStyle.Render(msg))
}
