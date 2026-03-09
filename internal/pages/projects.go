package pages

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/carlfung1003/ssh-portfolio/internal/content"
	"github.com/carlfung1003/ssh-portfolio/internal/theme"
)

type ProjectsModel struct {
	projects []content.Project
	cursor   int
	width    int
	height   int
}

func NewProjects(projects []content.Project) ProjectsModel {
	return ProjectsModel{projects: projects}
}

func (m ProjectsModel) Init() tea.Cmd {
	return nil
}

func (m ProjectsModel) Update(msg tea.Msg) (ProjectsModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			if m.cursor < len(m.projects)-1 {
				m.cursor++
			}
		case "k", "up":
			if m.cursor > 0 {
				m.cursor--
			}
		}
	}
	return m, nil
}

func (m ProjectsModel) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	var b strings.Builder
	b.WriteString(theme.Title.Render("Projects") + "\n\n")

	cardWidth := min(m.width-4, 70)

	for i, p := range m.projects {
		cursor := "  "
		if i == m.cursor {
			cursor = theme.ActiveTab.Render("▸ ")
		}

		// Title line
		title := theme.BodyText.Bold(true).Render(p.Title)
		status := ""
		if p.Status == "active" {
			status = theme.StatusActive.Render(" ● active")
		} else if p.Status == "in-progress" {
			status = lipgloss.NewStyle().Foreground(theme.Amber).Render(" ◐ in-progress")
		}

		// Tags
		tags := ""
		for _, t := range p.Tags {
			tags += theme.Tag.Render("["+t+"]") + " "
		}

		// Description
		desc := theme.MutedText.Render(truncate(p.Description, cardWidth-4))

		// URLs
		urls := ""
		if p.LiveURL != "" {
			urls += theme.MutedText.Render("  live: ") + lipgloss.NewStyle().Foreground(theme.Cyan).Render(p.LiveURL)
		}
		if p.RepoURL != "" {
			if urls != "" {
				urls += "\n"
			}
			urls += theme.MutedText.Render("  repo: ") + lipgloss.NewStyle().Foreground(theme.Cyan).Render(p.RepoURL)
		}

		card := fmt.Sprintf("%s%s%s\n  %s\n  %s", cursor, title, status, desc, tags)
		if i == m.cursor && urls != "" {
			card += "\n" + urls
		}

		if i == m.cursor {
			card = theme.Card.Width(cardWidth).Render(card)
		}

		b.WriteString(card + "\n")
	}

	b.WriteString("\n" + theme.MutedText.Render("  ↑/↓ navigate"))

	return b.String()
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}
