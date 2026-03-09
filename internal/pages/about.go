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

// TypingTickMsg drives the typing animation
type TypingTickMsg struct{}

func typingTick() tea.Cmd {
	return tea.Tick(25*time.Millisecond, func(t time.Time) tea.Msg {
		return TypingTickMsg{}
	})
}

type AboutModel struct {
	content   string
	width     int
	height    int
	charIndex int  // how many runes of bio revealed so far
	done      bool // animation complete
	active    bool // is this page currently visible
	runes     []rune
}

func NewAbout(markdownContent string) AboutModel {
	return AboutModel{
		content: markdownContent,
		runes:   []rune(markdownContent),
		active:  true,
	}
}

func (m AboutModel) Init() tea.Cmd {
	return nil
}

// SetActive starts/resumes the typing animation when page becomes visible
func (m AboutModel) SetActive(active bool) (AboutModel, tea.Cmd) {
	m.active = active
	if active && !m.done {
		return m, typingTick()
	}
	return m, nil
}

func (m AboutModel) Update(msg tea.Msg) (AboutModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
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
		if !m.done && m.active {
			m.charIndex = len(m.runes)
			m.done = true
			return m, nil
		}
	}
	return m, nil
}

func (m AboutModel) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	art := lipgloss.NewStyle().
		Foreground(theme.Indigo).
		Bold(true).
		Render(asciiArt)

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
