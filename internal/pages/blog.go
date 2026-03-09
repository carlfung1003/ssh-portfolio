package pages

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/carlfung1003/ssh-portfolio/internal/content"
	"github.com/carlfung1003/ssh-portfolio/internal/theme"
)

// BlogMsg is sent when user wants to go back from a post to the list
type BlogBackMsg struct{}

type BlogModel struct {
	posts       []content.Post
	cursor      int
	reading     bool // true = viewing a post
	postContent string
	width       int
	height      int
	scrollY     int
}

func NewBlog(posts []content.Post) BlogModel {
	return BlogModel{posts: posts}
}

func (m BlogModel) Init() tea.Cmd {
	return nil
}

func (m BlogModel) Update(msg tea.Msg) (BlogModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		if m.reading {
			switch msg.String() {
			case "esc", "backspace":
				m.reading = false
				m.postContent = ""
				m.scrollY = 0
			case "j", "down":
				m.scrollY++
			case "k", "up":
				if m.scrollY > 0 {
					m.scrollY--
				}
			}
		} else {
			switch msg.String() {
			case "j", "down":
				if m.cursor < len(m.posts)-1 {
					m.cursor++
				}
			case "k", "up":
				if m.cursor > 0 {
					m.cursor--
				}
			case "enter":
				if m.cursor < len(m.posts) {
					m.reading = true
					m.scrollY = 0
					renderer, _ := glamour.NewTermRenderer(
						glamour.WithStylePath("dark"),
						glamour.WithWordWrap(min(m.width-6, 80)),
					)
					rendered, _ := renderer.Render(m.posts[m.cursor].Content)
					m.postContent = rendered
				}
			}
		}
	}
	return m, nil
}

func (m BlogModel) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	if m.reading {
		return m.viewPost()
	}
	return m.viewList()
}

// IsReading returns true when viewing a blog post (not the list)
func (m BlogModel) IsReading() bool {
	return m.reading
}

func (m BlogModel) viewList() string {
	var b strings.Builder
	b.WriteString(theme.Title.Render("Blog") + "\n\n")

	for i, p := range m.posts {
		cursor := "  "
		if i == m.cursor {
			cursor = theme.ActiveTab.Render("▸ ")
		}

		title := theme.BodyText.Bold(true).Render(p.Title)
		meta := theme.MutedText.Render(fmt.Sprintf("  %s · %s · %d min read", p.Date, p.Tags, p.ReadingTime))

		b.WriteString(cursor + title + "\n" + meta + "\n\n")
	}

	b.WriteString(theme.MutedText.Render("  ↑/↓ navigate · enter read"))

	return b.String()
}

func (m BlogModel) viewPost() string {
	var b strings.Builder

	header := theme.Title.Render(m.posts[m.cursor].Title)
	meta := theme.MutedText.Render(fmt.Sprintf("%s · %d min read", m.posts[m.cursor].Date, m.posts[m.cursor].ReadingTime))
	back := lipgloss.NewStyle().Foreground(theme.Slate).Render("esc ← back")

	b.WriteString(header + "  " + meta + "\n")
	b.WriteString(back + "\n\n")

	// Simple viewport: show lines based on scroll position
	lines := strings.Split(m.postContent, "\n")
	viewHeight := m.height - 8
	if viewHeight < 5 {
		viewHeight = 5
	}

	if m.scrollY >= len(lines) {
		m.scrollY = max(0, len(lines)-viewHeight)
	}

	end := m.scrollY + viewHeight
	if end > len(lines) {
		end = len(lines)
	}

	for _, line := range lines[m.scrollY:end] {
		b.WriteString(line + "\n")
	}

	// Scroll indicator
	if len(lines) > viewHeight {
		pct := 0
		if len(lines)-viewHeight > 0 {
			pct = m.scrollY * 100 / (len(lines) - viewHeight)
		}
		b.WriteString(theme.MutedText.Render(fmt.Sprintf("\n  ↑/↓ scroll · %d%%", pct)))
	}

	return b.String()
}
