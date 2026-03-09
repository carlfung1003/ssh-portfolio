package pages

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/carlfung1003/ssh-portfolio/internal/content"
	"github.com/carlfung1003/ssh-portfolio/internal/theme"
)

type LinksModel struct {
	links  []content.Link
	width  int
	height int
}

func NewLinks(links []content.Link) LinksModel {
	return LinksModel{links: links}
}

func (m LinksModel) Init() tea.Cmd {
	return nil
}

func (m LinksModel) Update(msg tea.Msg) (LinksModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m LinksModel) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	var b strings.Builder
	b.WriteString(theme.Title.Render("Links") + "\n\n")

	urlStyle := lipgloss.NewStyle().Foreground(theme.Cyan)

	for _, l := range m.links {
		label := theme.BodyText.Bold(true).Render("  " + l.Icon + " " + l.Label)
		url := urlStyle.Render("    " + l.URL)
		b.WriteString(label + "\n" + url + "\n\n")
	}

	b.WriteString(theme.MutedText.Render("  Copy URLs to open in your browser"))

	return b.String()
}
