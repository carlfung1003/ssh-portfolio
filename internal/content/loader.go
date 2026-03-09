package content

import (
	"embed"
	"encoding/json"
	"path"
	"sort"
	"strings"
)

//go:embed all:data
var dataFS embed.FS

type Project struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	Status      string   `json:"status"`
	LiveURL     string   `json:"liveUrl"`
	RepoURL     string   `json:"repoUrl"`
}

type Tool struct {
	Name        string `json:"name"`
	Category    string `json:"category"`
	Proficiency int    `json:"proficiency"`
}

type Experience struct {
	Company    string   `json:"company"`
	Role       string   `json:"role"`
	StartDate  string   `json:"startDate"`
	EndDate    string   `json:"endDate"`
	Type       string   `json:"type"`
	Highlights []string `json:"highlights"`
}

type Link struct {
	Label string `json:"label"`
	URL   string `json:"url"`
	Icon  string `json:"icon"`
}

type Post struct {
	Title       string `json:"title"`
	Slug        string `json:"slug"`
	Date        string `json:"date"`
	Tags        string `json:"tags"`
	ReadingTime int    `json:"readingTime"`
	Content     string `json:"-"`
}

type Data struct {
	About      string
	Projects   []Project
	Tools      []Tool
	Experience []Experience
	Links      []Link
	Posts      []Post
}

func Load() Data {
	d := Data{}

	// About
	if b, err := dataFS.ReadFile("data/about.md"); err == nil {
		d.About = string(b)
	}

	// Projects
	if b, err := dataFS.ReadFile("data/projects.json"); err == nil {
		json.Unmarshal(b, &d.Projects)
	}

	// Tools
	if b, err := dataFS.ReadFile("data/tools.json"); err == nil {
		json.Unmarshal(b, &d.Tools)
	}

	// Experience
	if b, err := dataFS.ReadFile("data/experience.json"); err == nil {
		json.Unmarshal(b, &d.Experience)
	}

	// Links
	if b, err := dataFS.ReadFile("data/links.json"); err == nil {
		json.Unmarshal(b, &d.Links)
	}

	// Blog posts
	entries, err := dataFS.ReadDir("data/posts")
	if err == nil {
		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
				continue
			}
			b, err := dataFS.ReadFile(path.Join("data/posts", entry.Name()))
			if err != nil {
				continue
			}
			post := parsePost(entry.Name(), string(b))
			d.Posts = append(d.Posts, post)
		}
		sort.Slice(d.Posts, func(i, j int) bool {
			return d.Posts[i].Date > d.Posts[j].Date
		})
	}

	return d
}

func parsePost(filename, raw string) Post {
	p := Post{
		Slug: strings.TrimSuffix(filename, ".md"),
	}

	// Simple frontmatter parser (---\n...\n---)
	if strings.HasPrefix(raw, "---\n") {
		parts := strings.SplitN(raw[4:], "\n---\n", 2)
		if len(parts) == 2 {
			parseFrontmatter(parts[0], &p)
			p.Content = strings.TrimSpace(parts[1])
		} else {
			p.Content = raw
		}
	} else {
		p.Content = raw
	}

	if p.Title == "" {
		p.Title = p.Slug
	}

	return p
}

func parseFrontmatter(fm string, p *Post) {
	for _, line := range strings.Split(fm, "\n") {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		val = strings.Trim(val, "\"")
		switch key {
		case "title":
			p.Title = val
		case "date":
			p.Date = val
		case "tags":
			p.Tags = val
		case "readingTime":
			// ignore parse error
			p.ReadingTime = 0
			for _, c := range val {
				if c >= '0' && c <= '9' {
					p.ReadingTime = p.ReadingTime*10 + int(c-'0')
				}
			}
		}
	}
}
