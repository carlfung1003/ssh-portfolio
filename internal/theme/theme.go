package theme

import (
	"github.com/charmbracelet/lipgloss"
)

// Colors matching ai-journey web portfolio
var (
	Indigo    = lipgloss.Color("#6366F1")
	Purple    = lipgloss.Color("#8B5CF6")
	Slate     = lipgloss.Color("#64748B")
	Navy      = lipgloss.Color("#0F172A")
	DarkSlate = lipgloss.Color("#1E293B")
	Border    = lipgloss.Color("#334155")
	Text      = lipgloss.Color("#E2E8F0")
	Muted     = lipgloss.Color("#94A3B8")
	Green     = lipgloss.Color("#10B981")
	Amber     = lipgloss.Color("#F59E0B")
	Cyan      = lipgloss.Color("#06B6D4")
)

// Reusable styles
var (
	ActiveTab = lipgloss.NewStyle().
			Foreground(Indigo).
			Bold(true).
			Underline(true)

	InactiveTab = lipgloss.NewStyle().
			Foreground(Slate)

	Title = lipgloss.NewStyle().
		Foreground(Indigo).
		Bold(true)

	Subtitle = lipgloss.NewStyle().
			Foreground(Purple)

	MutedText = lipgloss.NewStyle().
			Foreground(Muted)

	BodyText = lipgloss.NewStyle().
			Foreground(Text)

	Card = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Border).
		Padding(1, 2)

	Tag = lipgloss.NewStyle().
		Foreground(Purple).
		Bold(true)

	StatusActive = lipgloss.NewStyle().
			Foreground(Green).
			Bold(true)

	FooterStyle = lipgloss.NewStyle().
			Foreground(Slate)

	FooterKeyStyle = lipgloss.NewStyle().
			Foreground(Indigo).
			Bold(true)
)

// ProficiencyBar renders a skill level as filled/empty blocks
func ProficiencyBar(level, max int) string {
	filled := ""
	empty := ""
	for i := 0; i < max; i++ {
		if i < level {
			filled += "█"
		} else {
			empty += "░"
		}
	}
	return lipgloss.NewStyle().Foreground(Indigo).Render(filled) +
		lipgloss.NewStyle().Foreground(Border).Render(empty)
}
