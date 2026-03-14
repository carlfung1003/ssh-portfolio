package pages

import (
	"strings"
	"time"
	"unicode/utf8"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/carlfung1003/ssh-portfolio/internal/theme"
)

var asciiArt = `
 @@@@@@@@@@@@@@@@@@@@@@@@@@@@@%%%%%%%%%%%%%%%%%%%%@
 @@@@@@@@@@@@@@@@@@@@@@@@@@@@%%%#################%%%
 @@@@@@@@@@@@@@@@@@@@@%%%%%%%%%%###****###*****###%%
 @@@@@@@@@@@@@@@@%%###%%%%#******####**************#
 @@@@@@@@@@@@@@@%%%%@@@%@#+##%%@%%%%#***************
 @@@@@@@%%%%%%%@@%#**++*#%%####%%@@@@@@@@%###*******
 @@@@@@%%######%#+--==+=-=+###%%%%%@@@@@@@@@@%******
 @%%%%%###*****+----+:.:+**++*%+*@@@@@@@@@@@@@%*****
 %%######******=:::-.:-:::--+++**#*#%%%@@%%%%@@@@#*+*
 %%###********=.. :-:. . :::-****%@@@@@@@@@@@@@%%#++*
 %%###********+.  .:.  . -:.=-+**##%%@@@@@@@@@@@%%*++
 %####*******+     .    :.:+*******##%%@@@@@@@@@@@@#*
 ######******=.   :..  :.-++=====++++**#@@@@@@@@@@@@@
 ######******:   .:.. ::=*=:.....:-=+***%@@@@@@@@@@@@
 ###**#******:  .:::. ++*=..:*#%%#%@%%###%@@@@@@@@@@@
 ****#*******+  :..:.=**+:-*#**-::=*%@@%##%%@@@@@@@@@
 ************. .:=**#*#=.:#*+*#%%%%@@%%####%@@@@@@@@@
 ************=.:-+=+*+*+ :*-:+=--+##*++++*#%@@@@@@@@
 *************+::=@@@=%- :*#......:*..:=+*#%@@@@@@@@@
 ***********##*+=+-:=.  :--:.  . ...:=**###%@@@@@@@%%
 ***************+...   .::-+:  .. .:=+*####%%@%%%@@#%
 *************+=:  ..-+++=+*-     .-=***####%%#%@%#*@
 **************. =#*@@%##+    .:=++**#####%*%%#*##**++
 **************-  -:-:==-.   .:-=++**####**#@#*#@**++
 **************+     -:.....::--=+***####***##%@%++++
 ***************:  -=*%%%@@%*-.:=+***##%%#*+***@@*+++
 ***************= -=:..:---...-=+**##%%%%***++*%%*+++
 ***************+.  :-=--:---=+**#%%@@@@%****+=#*=+++
 ****************-    ..  ..:-=+*#%@@@@@@@%%%#**=+*+*
`

// Indigo gradient palette for the wave animation
var artGradient = []lipgloss.Color{
	"#3730A3",
	"#4338CA",
	"#4F46E5",
	"#5558EB",
	"#6366F1", // primary indigo
	"#7577F4",
	"#818CF8",
	"#9BA3FA",
	"#A5B4FC", // brightest
	"#9BA3FA",
	"#818CF8",
	"#7577F4",
	"#6366F1",
	"#5558EB",
	"#4F46E5",
	"#4338CA",
}

// TypingTickMsg drives the bio typing animation
type TypingTickMsg struct{}

// ArtTickMsg drives the ASCII art animation
type ArtTickMsg struct{}

func typingTick() tea.Cmd {
	return tea.Tick(25*time.Millisecond, func(t time.Time) tea.Msg {
		return TypingTickMsg{}
	})
}

func artTick() tea.Cmd {
	return tea.Tick(120*time.Millisecond, func(t time.Time) tea.Msg {
		return ArtTickMsg{}
	})
}

type AboutModel struct {
	content   string
	width     int
	height    int
	charIndex int  // how many runes of bio revealed so far
	done      bool // bio animation complete
	active    bool // is this page currently visible
	runes     []rune

	// ASCII art animation state
	artLines []string // pre-split art lines
	artLine  int      // lines revealed so far (0 = none)
	artFrame int      // frame counter for gradient wave
	artDone  bool     // reveal complete
}

func NewAbout(markdownContent string) AboutModel {
	lines := strings.Split(strings.TrimSpace(asciiArt), "\n")
	return AboutModel{
		content:  markdownContent,
		runes:    []rune(markdownContent),
		active:   true,
		artLines: lines,
	}
}

func (m AboutModel) Init() tea.Cmd {
	return nil
}

// SetActive starts/resumes animations when page becomes visible
func (m AboutModel) SetActive(active bool) (AboutModel, tea.Cmd) {
	m.active = active
	if active {
		var cmds []tea.Cmd
		if !m.done {
			cmds = append(cmds, typingTick())
		}
		cmds = append(cmds, artTick())
		return m, tea.Batch(cmds...)
	}
	return m, nil
}

func (m AboutModel) Update(msg tea.Msg) (AboutModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case ArtTickMsg:
		if !m.active {
			return m, nil
		}
		if !m.artDone {
			// Reveal phase: 3 lines per tick
			m.artLine += 3
			if m.artLine >= len(m.artLines) {
				m.artLine = len(m.artLines)
				m.artDone = true
			}
		}
		m.artFrame++
		return m, artTick()
	case TypingTickMsg:
		if !m.active || m.done {
			return m, nil
		}
		// Reveal 2 characters per tick for faster typing
		m.charIndex += 2
		if m.charIndex >= len(m.runes) {
			m.charIndex = len(m.runes)
			m.done = true
			return m, nil
		}
		return m, typingTick()
	case tea.KeyMsg:
		// Any key press skips to full reveal
		if (!m.done || !m.artDone) && m.active {
			m.charIndex = len(m.runes)
			m.done = true
			m.artLine = len(m.artLines)
			m.artDone = true
			return m, nil
		}
	}
	return m, nil
}

func (m AboutModel) renderArt() string {
	if len(m.artLines) == 0 {
		return ""
	}

	visibleCount := m.artLine
	if visibleCount > len(m.artLines) {
		visibleCount = len(m.artLines)
	}

	var result strings.Builder
	for i := 0; i < len(m.artLines); i++ {
		if i > 0 {
			result.WriteString("\n")
		}
		if i >= visibleCount {
			// Not yet revealed — blank line for layout stability
			result.WriteString(strings.Repeat(" ", len(m.artLines[i])))
			continue
		}

		// Color based on line position + frame for moving gradient wave
		colorIdx := (i + m.artFrame) % len(artGradient)
		color := artGradient[colorIdx]
		styled := lipgloss.NewStyle().
			Foreground(color).
			Bold(true).
			Render(m.artLines[i])
		result.WriteString(styled)
	}

	return result.String()
}

func (m AboutModel) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	art := m.renderArt()

	// Show bio text up to current charIndex
	visibleText := string(m.runes[:m.charIndex])

	// Add blinking cursor if still typing
	cursor := ""
	if !m.done {
		cursor = theme.ActiveTab.Render("▌")
	}

	// Word-wrap the visible text
	bio := wordWrap(visibleText, min(m.width-10, 60))
	bio = theme.BodyText.Render(bio) + cursor

	// Skip glamour — plain styled text looks cleaner for typing effect
	_ = utf8.RuneCountInString // keep import

	var result string
	if m.width >= 120 {
		artBox := lipgloss.NewStyle().
			Width(54).
			Render(art)
		bioBox := lipgloss.NewStyle().
			Width(m.width - 58).
			Render(bio)
		result = lipgloss.JoinHorizontal(lipgloss.Top, artBox, bioBox)
	} else {
		result = art + "\n\n" + bio
	}

	return result
}

// wordWrap wraps text at word boundaries to fit within maxWidth
func wordWrap(text string, maxWidth int) string {
	if maxWidth <= 0 {
		return text
	}
	var result strings.Builder
	for _, paragraph := range strings.Split(text, "\n") {
		if paragraph == "" {
			result.WriteString("\n")
			continue
		}
		words := strings.Fields(paragraph)
		lineLen := 0
		for i, word := range words {
			wLen := utf8.RuneCountInString(word)
			if i > 0 && lineLen+1+wLen > maxWidth {
				result.WriteString("\n")
				lineLen = 0
			} else if i > 0 {
				result.WriteString(" ")
				lineLen++
			}
			result.WriteString(word)
			lineLen += wLen
		}
		result.WriteString("\n")
	}
	return strings.TrimRight(result.String(), "\n")
}
