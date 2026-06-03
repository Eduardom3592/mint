package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbletea"

	"github.com/programmersd21/mint/internal/cache"
	"github.com/programmersd21/mint/internal/models"
)

type dependencyNode struct {
	Project    *models.Project
	Version    *models.Version
	Dependency models.Dependency
	Children   []dependencyNode
	Expanded   bool
	Depth      int
}

type dependencyModel struct {
	cache  *cache.Cache
	root   *dependencyNode
	cursor int
	page   dependencyPage
	styles Styles
}

func newDependencyModel(c *cache.Cache) dependencyModel {
	return dependencyModel{
		cache:  c,
		page:   dependencyBrowse,
		styles: glob,
	}
}

func (m *dependencyModel) setStyles(s Styles) {
	m.styles = s
}

func (m dependencyModel) Init() tea.Cmd { return nil }

func (m dependencyModel) Update(msg tea.Msg) (dependencyModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Escape):
			return m, nil

		case key.Matches(msg, keys.Up):
			if m.cursor > 0 {
				m.cursor--
			}
			return m, nil

		case key.Matches(msg, keys.Down):
			if m.cursor < m.totalVisible()-1 {
				m.cursor++
			}
			return m, nil

		case key.Matches(msg, keys.Enter):
			if m.cursor >= 0 {
				visible := m.getVisibleNodes()
				if m.cursor < len(visible) {
					visible[m.cursor].Expanded = !visible[m.cursor].Expanded
				}
			}
			return m, nil
		}
	}

	return m, nil
}

func (m dependencyModel) totalVisible() int {
	return len(m.getVisibleNodes())
}

func (m dependencyModel) getVisibleNodes() []*dependencyNode {
	var visible []*dependencyNode
	m.collectVisible(m.root, &visible)
	return visible
}

func (m dependencyModel) collectVisible(node *dependencyNode, result *[]*dependencyNode) {
	*result = append(*result, node)
	if node.Expanded {
		for i := range node.Children {
			m.collectVisible(&node.Children[i], result)
		}
	}
}

func (m dependencyModel) View() string {
	var sb strings.Builder
	s := m.styles

	if m.root == nil {
		return fmt.Sprintf("\n  %s", s.Info.Render("no dependencies to display"))
	}

	sb.WriteString(s.Title.Render("dependencies"))
	sb.WriteString("\n\n")

	visible := m.getVisibleNodes()

	for i := range visible {
		node := visible[i]

		prefix := "  "
		expanded := " "

		if i == m.cursor {
			prefix = "> "
		}

		if len(node.Children) > 0 {
			if node.Expanded {
				expanded = "v"
			} else {
				expanded = ">"
			}
		}

		indent := strings.Repeat("  ", node.Depth)

		var label string
		if node.Project != nil {
			label = fmt.Sprintf("%s %s", expanded, node.Project.Title)
		} else if node.Dependency.ProjectID != nil {
			label = fmt.Sprintf("%s %s", expanded, *node.Dependency.ProjectID)
		} else if node.Dependency.VersionID != nil {
			label = fmt.Sprintf("%s version:%s", expanded, *node.Dependency.VersionID)
		} else {
			label = fmt.Sprintf("%s unknown", expanded)
		}

		depType := node.Dependency.DependencyType
		var typeTag string
		switch depType {
		case "required":
			typeTag = s.Warning.Render("required")
		case "optional":
			typeTag = s.Info.Render("optional")
		case "incompatible":
			typeTag = s.Error.Render("incompatible")
		default:
			typeTag = s.Info.Render(depType)
		}

		line := fmt.Sprintf("%s%s%s  %s", prefix, indent, label, typeTag)

		if i == m.cursor {
			sb.WriteString(s.Selected.Render(line))
		} else {
			sb.WriteString(line)
		}
		sb.WriteString("\n")
	}

	sb.WriteString("\n")
	sb.WriteString(s.Dimmed.Render("enter: expand/collapse - j/k: navigate - esc: back"))

	return sb.String()
}
