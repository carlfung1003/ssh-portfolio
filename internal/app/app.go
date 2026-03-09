package app

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/carlfung1003/ssh-portfolio/internal/content"
	"github.com/carlfung1003/ssh-portfolio/internal/pages"
	"github.com/carlfung1003/ssh-portfolio/internal/theme"
)

type page int

const (
	pageAbout page = iota
	pageProjects
	pageBlog
	pageSkills
	pageExperience
	pageLinks
)

var pageNames = []string{"About", "Projects", "Blog", "Skills", "Experience", "Links"}

type Model struct {
	activePage page
	about      pages.AboutModel
	projects   pages.ProjectsModel
	blog       pages.BlogModel
	skills     pages.SkillsModel
	experience pages.ExperienceModel
	links      pages.LinksModel
	width      int
	height     int
}

func NewModel(data content.Data, width, height int) Model {
	return Model{
		activePage: pageAbout,
		about:      pages.NewAbout(data.About),
		projects:   pages.NewProjects(data.Projects),
		blog:       pages.NewBlog(data.Posts),
		skills:     pages.NewSkills(data.Tools),
		experience: pages.NewExperience(data.Experience),
		links:      pages.NewLinks(data.Links),
		width:      width,
		height:     height,
	}
}

func (m Model) Init() tea.Cmd {
	// Start typing animation on About page
	var cmd tea.Cmd
	m.about, cmd = m.about.SetActive(true)
	return cmd
}

// pageNeedsUpDown returns true for pages that use up/down keys internally
func (m Model) pageNeedsUpDown() bool {
	switch m.activePage {
	case pageProjects, pageBlog:
		return true
	}
	return false
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Forward to all pages
		m.about, _ = m.about.Update(msg)
		m.projects, _ = m.projects.Update(msg)
		m.blog, _ = m.blog.Update(msg)
		m.skills, _ = m.skills.Update(msg)
		m.experience, _ = m.experience.Update(msg)
		m.links, _ = m.links.Update(msg)
		return m, nil

	case tea.KeyMsg:
		// Global quit — ctrl+c always works
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
		// q only quits when not reading a blog post
		if msg.Type == tea.KeyRunes && msg.String() == "q" && !m.blog.IsReading() {
			return m, tea.Quit
		}

		// When reading a blog post, only forward to blog (no page navigation)
		if m.blog.IsReading() && m.activePage == pageBlog {
			var cmd tea.Cmd
			m.blog, cmd = m.blog.Update(msg)
			return m, cmd
		}

		// Page navigation — use Type for arrow keys (reliable across terminals)
		goLeft := msg.Type == tea.KeyLeft || msg.Type == tea.KeyShiftTab
		goRight := msg.Type == tea.KeyRight || msg.Type == tea.KeyTab

		// Also support h/l for vim users
		if msg.Type == tea.KeyRunes {
			switch msg.String() {
			case "h":
				goLeft = true
			case "l":
				goRight = true
			}
		}

		newPage := m.activePage
		if goLeft {
			newPage = (m.activePage - 1 + page(len(pageNames))) % page(len(pageNames))
		}
		if goRight {
			newPage = (m.activePage + 1) % page(len(pageNames))
		}

		// Number key jump (only for rune keys)
		if msg.Type == tea.KeyRunes {
			switch msg.String() {
			case "1":
				newPage = pageAbout
			case "2":
				newPage = pageProjects
			case "3":
				newPage = pageBlog
			case "4":
				newPage = pageSkills
			case "5":
				newPage = pageExperience
			case "6":
				newPage = pageLinks
			}
		}

		if newPage != m.activePage {
			return m.switchPage(newPage)
		}
	}

	// Forward remaining keys to active page (up/down, enter, etc.)
	var cmd tea.Cmd
	switch m.activePage {
	case pageAbout:
		m.about, cmd = m.about.Update(msg)
	case pageProjects:
		m.projects, cmd = m.projects.Update(msg)
	case pageBlog:
		m.blog, cmd = m.blog.Update(msg)
	case pageSkills:
		m.skills, cmd = m.skills.Update(msg)
	case pageExperience:
		m.experience, cmd = m.experience.Update(msg)
	case pageLinks:
		m.links, cmd = m.links.Update(msg)
	}

	return m, cmd
}

func (m Model) switchPage(newPage page) (tea.Model, tea.Cmd) {
	// Deactivate About typing if leaving
	if m.activePage == pageAbout && newPage != pageAbout {
		m.about, _ = m.about.SetActive(false)
	}

	m.activePage = newPage

	// Activate About typing if arriving
	if newPage == pageAbout {
		var cmd tea.Cmd
		m.about, cmd = m.about.SetActive(true)
		return m, cmd
	}

	return m, nil
}

func (m Model) View() string {
	if m.width == 0 {
		return "Initializing..."
	}

	header := m.renderHeader()
	footer := m.renderFooter()

	// Page content
	var pageView string
	switch m.activePage {
	case pageAbout:
		pageView = m.about.View()
	case pageProjects:
		pageView = m.projects.View()
	case pageBlog:
		pageView = m.blog.View()
	case pageSkills:
		pageView = m.skills.View()
	case pageExperience:
		pageView = m.experience.View()
	case pageLinks:
		pageView = m.links.View()
	}

	content := lipgloss.NewStyle().
		Padding(1, 2).
		Width(m.width).
		Height(m.height - 4). // Reserve space for header + footer
		Render(pageView)

	return header + "\n" + content + "\n" + footer
}

func (m Model) renderHeader() string {
	var tabs []string
	for i, name := range pageNames {
		label := fmt.Sprintf(" %d %s ", i+1, name)
		if page(i) == m.activePage {
			tabs = append(tabs, theme.ActiveTab.Render(label))
		} else {
			tabs = append(tabs, theme.InactiveTab.Render(label))
		}
	}

	tabBar := strings.Join(tabs, theme.MutedText.Render("│"))

	title := lipgloss.NewStyle().
		Foreground(theme.Indigo).
		Bold(true).
		Render(" Carl Fung ")

	line := lipgloss.NewStyle().
		Foreground(theme.Border).
		Render(strings.Repeat("─", max(0, m.width-lipgloss.Width(title)-2)))

	return title + line + "\n" + tabBar
}

func (m Model) renderFooter() string {
	keys := theme.FooterKeyStyle.Render("←/→") + theme.FooterStyle.Render(" navigate · ") +
		theme.FooterKeyStyle.Render("1-6") + theme.FooterStyle.Render(" jump · ") +
		theme.FooterKeyStyle.Render("↑/↓") + theme.FooterStyle.Render(" scroll · ") +
		theme.FooterKeyStyle.Render("q") + theme.FooterStyle.Render(" quit")

	version := theme.MutedText.Render("v0.1.0")

	gap := max(0, m.width-lipgloss.Width(keys)-lipgloss.Width(version)-4)

	return lipgloss.NewStyle().
		Foreground(theme.Border).
		Render(strings.Repeat("─", m.width)) + "\n" +
		"  " + keys + strings.Repeat(" ", gap) + version
}
