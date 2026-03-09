package pages

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/carlfung1003/ssh-portfolio/internal/content"
	"github.com/carlfung1003/ssh-portfolio/internal/theme"
)

type SkillsModel struct {
	tools  []content.Tool
	width  int
	height int
}

func NewSkills(tools []content.Tool) SkillsModel {
	return SkillsModel{tools: tools}
}

func (m SkillsModel) Init() tea.Cmd {
	return nil
}

func (m SkillsModel) Update(msg tea.Msg) (SkillsModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m SkillsModel) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	var b strings.Builder
	b.WriteString(theme.Title.Render("Skills & Tools") + "\n\n")

	// Group by category
	categories := []string{"AI", "Language", "Framework", "Database", "Infra"}
	grouped := make(map[string][]content.Tool)
	for _, t := range m.tools {
		grouped[t.Category] = append(grouped[t.Category], t)
	}

	for _, cat := range categories {
		tools, ok := grouped[cat]
		if !ok {
			continue
		}

		b.WriteString(theme.Subtitle.Render("  "+cat) + "\n")
		for _, t := range tools {
			bar := theme.ProficiencyBar(t.Proficiency, 5)
			name := theme.BodyText.Render(fmt.Sprintf("    %-16s", t.Name))
			b.WriteString(name + bar + "\n")
		}
		b.WriteString("\n")
	}

	return b.String()
}
