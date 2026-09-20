package ssg

import (
	"bufio"
	"bytes"
	"fmt"
	"html"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type PostFrontmatter struct {
	Title     string `yaml:"title"`
	Slug      string `yaml:"slug"`
	Author    string `yaml:"author"`
	Date      string `yaml:"date"`
	Published *bool  `yaml:"published"`
	Summary   string `yaml:"summary"`
}

type Post struct {
	Title         string
	Slug          string
	Author        string
	Date          time.Time
	FormattedDate string
	Published     bool
	Summary       string
	Content       string
	HTMLContent   template.HTML
	BasePath      string
}

func ParsePost(path string) (*Post, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	content := string(data)
	var fm PostFrontmatter
	body := content

	if strings.HasPrefix(content, "---\n") || strings.HasPrefix(content, "---\r\n") {
		parts := strings.SplitN(content, "---", 3)
		if len(parts) >= 3 {
			if err := yaml.Unmarshal([]byte(parts[1]), &fm); err != nil {
				return nil, fmt.Errorf("parsing frontmatter in %s: %w", path, err)
			}
			body = strings.TrimSpace(parts[2])
		}
	}

	if fm.Title == "" {
		// Fallback title to filename
		base := filepath.Base(path)
		fm.Title = strings.TrimSuffix(base, filepath.Ext(base))
	}
	if fm.Slug == "" {
		fm.Slug = slugify(fm.Title)
	}

	published := true
	if fm.Published != nil {
		published = *fm.Published
	}

	var parsedDate time.Time
	if fm.Date != "" {
		for _, layout := range []string{"2006-01-02", "2006-01-02 15:04:05", time.RFC3339} {
			if t, err := time.Parse(layout, fm.Date); err == nil {
				parsedDate = t
				break
			}
		}
	}
	if parsedDate.IsZero() {
		fi, _ := os.Stat(path)
		if fi != nil {
			parsedDate = fi.ModTime()
		} else {
			parsedDate = time.Now()
		}
	}

	htmlContent := renderMarkdownToHTML(body)

	summary := fm.Summary
	if summary == "" {
		summary = excerpt(body, 160)
	}

	return &Post{
		Title:         fm.Title,
		Slug:          fm.Slug,
		Author:        fm.Author,
		Date:          parsedDate,
		FormattedDate: parsedDate.Format("January 2, 2006"),
		Published:     published,
		Summary:       summary,
		Content:       body,
		HTMLContent:   template.HTML(htmlContent),
	}, nil
}

func LoadPosts(contentDir string) ([]*Post, error) {
	var posts []*Post

	err := filepath.Walk(contentDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(strings.ToLower(info.Name()), ".md") {
			p, parseErr := ParsePost(path)
			if parseErr != nil {
				return parseErr
			}
			if p.Published {
				posts = append(posts, p)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	// Sort posts by date descending
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Date.After(posts[j].Date)
	})

	return posts, nil
}

var nonWordRegex = regexp.MustCompile(`[^\w\s-]`)
var spaceRegex = regexp.MustCompile(`[-\s]+`)

func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = nonWordRegex.ReplaceAllString(s, "")
	s = spaceRegex.ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}

func excerpt(s string, maxLen int) string {
	lines := strings.Split(s, "\n")
	var textLines []string
	for _, l := range lines {
		trimmed := strings.TrimSpace(l)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, "```") {
			continue
		}
		textLines = append(textLines, trimmed)
		if len(strings.Join(textLines, " ")) >= maxLen {
			break
		}
	}
	combined := strings.Join(textLines, " ")
	if len(combined) > maxLen {
		return combined[:maxLen] + "..."
	}
	return combined
}

// renderMarkdownToHTML renders markdown to safe HTML
func renderMarkdownToHTML(md string) string {
	var out bytes.Buffer
	scanner := bufio.NewScanner(strings.NewReader(md))
	inCodeBlock := false
	codeLang := ""
	var codeBuf bytes.Buffer
	inList := false
	listType := "" // "ul" or "ol"

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		// Code blocks
		if strings.HasPrefix(trimmed, "```") {
			if inCodeBlock {
				inCodeBlock = false
				out.WriteString(fmt.Sprintf("<pre><code class=\"language-%s\">%s</code></pre>\n", html.EscapeString(codeLang), html.EscapeString(codeBuf.String())))
				codeBuf.Reset()
				codeLang = ""
			} else {
				if inList {
					out.WriteString(fmt.Sprintf("</%s>\n", listType))
					inList = false
				}
				inCodeBlock = true
				codeLang = strings.TrimSpace(strings.TrimPrefix(trimmed, "```"))
			}
			continue
		}

		if inCodeBlock {
			codeBuf.WriteString(line)
			codeBuf.WriteString("\n")
			continue
		}

		// Close list if blank line or heading
		if trimmed == "" {
			if inList {
				out.WriteString(fmt.Sprintf("</%s>\n", listType))
				inList = false
			}
			continue
		}

		// Headings
		if strings.HasPrefix(line, "### ") {
			if inList {
				out.WriteString(fmt.Sprintf("</%s>\n", listType))
				inList = false
			}
			out.WriteString(fmt.Sprintf("<h3>%s</h3>\n", renderInline(strings.TrimPrefix(line, "### "))))
			continue
		}
		if strings.HasPrefix(line, "## ") {
			if inList {
				out.WriteString(fmt.Sprintf("</%s>\n", listType))
				inList = false
			}
			out.WriteString(fmt.Sprintf("<h2>%s</h2>\n", renderInline(strings.TrimPrefix(line, "## "))))
			continue
		}
		if strings.HasPrefix(line, "# ") {
			if inList {
				out.WriteString(fmt.Sprintf("</%s>\n", listType))
				inList = false
			}
			out.WriteString(fmt.Sprintf("<h1>%s</h1>\n", renderInline(strings.TrimPrefix(line, "# "))))
			continue
		}

		// Blockquote
		if strings.HasPrefix(trimmed, "> ") {
			if inList {
				out.WriteString(fmt.Sprintf("</%s>\n", listType))
				inList = false
			}
			out.WriteString(fmt.Sprintf("<blockquote><p>%s</p></blockquote>\n", renderInline(strings.TrimPrefix(trimmed, "> "))))
			continue
		}

		// Unordered list
		if strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "* ") {
			if !inList || listType != "ul" {
				if inList {
					out.WriteString(fmt.Sprintf("</%s>\n", listType))
				}
				out.WriteString("<ul>\n")
				inList = true
				listType = "ul"
			}
			itemText := strings.TrimPrefix(strings.TrimPrefix(trimmed, "- "), "* ")
			out.WriteString(fmt.Sprintf("  <li>%s</li>\n", renderInline(itemText)))
			continue
		}

		// Ordered list: number followed by ". "
		if isOrderedList(trimmed) {
			if !inList || listType != "ol" {
				if inList {
					out.WriteString(fmt.Sprintf("</%s>\n", listType))
				}
				out.WriteString("<ol>\n")
				inList = true
				listType = "ol"
			}
			parts := strings.SplitN(trimmed, ". ", 2)
			itemText := parts[1]
			out.WriteString(fmt.Sprintf("  <li>%s</li>\n", renderInline(itemText)))
			continue
		}

		// Regular paragraph
		if inList {
			out.WriteString(fmt.Sprintf("</%s>\n", listType))
			inList = false
		}
		out.WriteString(fmt.Sprintf("<p>%s</p>\n", renderInline(line)))
	}

	if inList {
		out.WriteString(fmt.Sprintf("</%s>\n", listType))
	}
	if inCodeBlock {
		out.WriteString(fmt.Sprintf("<pre><code class=\"language-%s\">%s</code></pre>\n", html.EscapeString(codeLang), html.EscapeString(codeBuf.String())))
	}

	return out.String()
}

func isOrderedList(s string) bool {
	idx := strings.Index(s, ". ")
	if idx <= 0 || idx > 4 {
		return false
	}
	for i := 0; i < idx; i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}

var (
	linkRegex       = regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
	boldRegex       = regexp.MustCompile(`\*\*([^*]+)\*\*`)
	italicRegex     = regexp.MustCompile(`\*([^*]+)\*`)
	inlineCodeRegex = regexp.MustCompile("`([^`]+)`")
)

func renderInline(s string) string {
	// First escape HTML entities
	escaped := html.EscapeString(s)

	// Restore inline code
	escaped = inlineCodeRegex.ReplaceAllString(escaped, `<code>$1</code>`)

	// Bold
	escaped = boldRegex.ReplaceAllString(escaped, `<strong>$1</strong>`)

	// Italic
	escaped = italicRegex.ReplaceAllString(escaped, `<em>$1</em>`)

	// Links
	escaped = linkRegex.ReplaceAllString(escaped, `<a href="$2">$1</a>`)

	return escaped
}
