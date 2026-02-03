package cmd

import "github.com/charmbracelet/lipgloss"

var (
	primaryColor = lipgloss.Color("#7C3AED")
	successColor = lipgloss.Color("#10B981")
	warningColor = lipgloss.Color("#F59E0B")
	errorColor   = lipgloss.Color("#EF4444")
	mutedColor   = lipgloss.Color("#6B7280")
	infoColor    = lipgloss.Color("#3B82F6")

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor)

	idStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#A78BFA"))

	subjectStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F3F4F6"))

	labelStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Width(12)

	valueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E5E7EB"))

	descriptionStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#D1D5DB")).
				PaddingLeft(2).
				BorderLeft(true).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(mutedColor)

	detailBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(1, 2).
			MarginTop(1)

	headerBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#A78BFA")).
			Padding(0, 2).
			MarginBottom(1)

	propLabelStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Width(14)

	propValueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E5E7EB"))
)

func statusStyle(status string) lipgloss.Style {
	base := lipgloss.NewStyle().Bold(true).Padding(0, 1)

	switch status {
	case "New":
		return base.Background(lipgloss.Color("#1E3A5F")).Foreground(lipgloss.Color("#60A5FA"))
	case "In Progress", "Doing":
		return base.Background(lipgloss.Color("#164E63")).Foreground(lipgloss.Color("#22D3EE"))
	case "Code review":
		return base.Background(lipgloss.Color("#4C1D95")).Foreground(lipgloss.Color("#A78BFA"))
	case "Homolog":
		return base.Background(lipgloss.Color("#713F12")).Foreground(lipgloss.Color("#FCD34D"))
	case "Done", "Closed":
		return base.Background(lipgloss.Color("#14532D")).Foreground(lipgloss.Color("#4ADE80"))
	case "Blocked":
		return base.Background(lipgloss.Color("#7F1D1D")).Foreground(lipgloss.Color("#FCA5A5"))
	default:
		return base.Background(lipgloss.Color("#374151")).Foreground(lipgloss.Color("#D1D5DB"))
	}
}
