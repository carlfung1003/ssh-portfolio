package pages

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/carlfung1003/ssh-portfolio/internal/content"
	"github.com/carlfung1003/ssh-portfolio/internal/theme"
)

type ExperienceModel struct {
	items  []content.Experience
	width  int
	height int
}

func NewExperience(items []content.Experience) ExperienceModel {
	return ExperienceModel{items: items}
}

func (m ExperienceModel) Init() tea.Cmd {
	return nil
}

func (m ExperienceModel) Update(msg tea.Msg) (ExperienceModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m ExperienceModel) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	var b strings.Builder
	b.WriteString(theme.Title.Render("Experience") + "\n\n")

	workStyle := lipgloss.NewStyle().Foreground(theme.Indigo)
	eduStyle := lipgloss.NewStyle().Foreground(theme.Green)

	for i, e := range m.items {
		// Timeline connector
		connector := "│"
		if i == len(m.items)-1 {
			connector = "└"
		}
		dot := "●"

		var typeStyle lipgloss.Style
		if e.Type == "education" {
			typeStyle = eduStyle
			dot = "◆"
		} else {
			typeStyle = workStyle
		}

		dates := theme.MutedText.Render(fmt.Sprintf("%s – %s", e.StartDate, e.EndDate))
		company := typeStyle.Bold(true).Render(e.Company)
		role := theme.BodyText.Render(e.Role)

		b.WriteString(fmt.Sprintf("  %s %s  %s\n", typeStyle.Render(dot), company, dates))
		b.WriteString(fmt.Sprintf("  %s    %s\n", theme.MutedText.Render(connector), role))

		for _, h := range e.Highlights {
			b.WriteString(fmt.Sprintf("  %s      %s\n", theme.MutedText.Render("│"), theme.MutedText.Render("· "+h)))
		}
		b.WriteString("\n")
	}

	b.WriteString(theme.MutedText.Render("  ● work  ◆ education"))

	return b.String()
}
