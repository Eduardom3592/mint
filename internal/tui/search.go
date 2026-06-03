package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"

	"github.com/programmersd21/mint/internal/api"
	"github.com/programmersd21/mint/internal/cache"
	"github.com/programmersd21/mint/internal/models"
)

type searchModel struct {
	client            *api.Client
	cache             *cache.Cache
	textInput         textinput.Model
	results           []models.SearchHit
	cursor            int
	loading           bool
	err               error
	page              searchPage
	selectedProjectID string

	filterType    string
	filterLoader  string
	filterVersion string
	sort          api.SearchSort

	totalHits int
	offset    int
	limit     int

	frame         int
	searchVersion int
	height        int
	width         int
	styles        Styles
}

type debounceSearchMsg struct {
	version int
}

func newSearchModel(client *api.Client, c *cache.Cache) searchModel {
	ti := textinput.New()
	ti.Placeholder = "search modrinth..."
	ti.Focus()
	ti.Width = 60
	ti.Prompt = "/ "

	return searchModel{
		client:    client,
		cache:     c,
		textInput: ti,
		results:   make([]models.SearchHit, 0),
		sort:      api.SortRelevance,
		limit:     50,
		page:      searchList,
		styles:    glob,
	}
}

func (m *searchModel) setStyles(s Styles) {
	m.styles = s
}

func (m searchModel) Init() tea.Cmd { return nil }

func (m searchModel) Update(msg tea.Msg) (searchModel, tea.Cmd) {
	switch msg := msg.(type) {
	case searchResultsMsg:
		m.results = msg.results
		m.totalHits = msg.total
		m.loading = false
		return m, nil

	case searchErrorMsg:
		m.err = msg.err
		m.loading = false
		return m, nil

	case debounceSearchMsg:
		if msg.version == m.searchVersion {
			return m, m.search()
		}
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Escape):
			if m.textInput.Focused() {
				m.textInput.Blur()
			}
			return m, nil

		case key.Matches(msg, keys.Enter):
			if m.textInput.Focused() {
				m.textInput.Blur()
				return m, m.search()
			}
			if len(m.results) > 0 && m.cursor >= 0 && m.cursor < len(m.results) {
				m.selectedProjectID = m.results[m.cursor].ProjectID
				m.page = searchProject
			}
			return m, nil

		case key.Matches(msg, keys.Up) && !m.textInput.Focused():
			if m.cursor > 0 {
				m.cursor--
			}
			return m, nil

		case key.Matches(msg, keys.Down) && !m.textInput.Focused():
			if m.cursor < len(m.results)-1 {
				m.cursor++
			}
			return m, nil

		case key.Matches(msg, keys.Home) && !m.textInput.Focused():
			m.cursor = 0
			return m, nil

		case key.Matches(msg, keys.End) && !m.textInput.Focused():
			m.cursor = len(m.results) - 1
			if m.cursor < 0 {
				m.cursor = 0
			}
			return m, nil

		case key.Matches(msg, keys.Slash):
			if m.textInput.Focused() {
				m.textInput.Blur()
			} else {
				m.textInput.Focus()
			}
			return m, nil
		}
	}

	if m.textInput.Focused() {
		var cmd tea.Cmd
		m.textInput, cmd = m.textInput.Update(msg)

		if _, ok := msg.(tea.KeyMsg); ok {
			val := m.textInput.Value()
			if val == "" {
				m.results = nil
				m.totalHits = 0
				m.err = nil
				m.loading = false
				m.cursor = 0
			} else if len(val) >= 3 {
				m.searchVersion++
				sv := m.searchVersion
				return m, tea.Batch(cmd, tea.Tick(300*time.Millisecond, func(t time.Time) tea.Msg {
					return debounceSearchMsg{version: sv}
				}))
			}
		}

		return m, cmd
	}

	return m, nil
}

func (m *searchModel) search() tea.Cmd {
	if len(m.textInput.Value()) < 3 {
		return nil
	}

	m.loading = true
	m.err = nil

	query := m.textInput.Value()

	cached, err := m.cache.GetSearch(query)
	if err == nil && cached != nil {
		m.results = cached.Hits
		m.totalHits = cached.TotalHits
		m.loading = false
		return nil
	}

	return func() tea.Msg {
		filter := api.SearchFilter{
			Query:  query,
			Sort:   m.sort,
			Offset: m.offset,
			Limit:  m.limit,
		}

		if m.filterType != "" {
			filter.ProjectType = m.filterType
		}
		if m.filterLoader != "" {
			filter.Loaders = []string{m.filterLoader}
		}
		if m.filterVersion != "" {
			filter.GameVersions = []string{m.filterVersion}
		}

		result, err := m.client.Search(filter)
		if err != nil {
			return searchErrorMsg{err: err}
		}

		_ = m.cache.CacheSearch(query, result, 5*time.Minute)

		return searchResultsMsg{results: result.Hits, total: result.TotalHits}
	}
}

type searchResultsMsg struct {
	results []models.SearchHit
	total   int
}

type searchErrorMsg struct {
	err error
}

func spinnerFrame(frame int) string {
	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	return frames[frame%len(frames)]
}

func (m searchModel) View() string {
	var sb strings.Builder
	s := m.styles

	if m.page != searchList {
		return ""
	}

	sb.WriteString(s.Title.Render("search"))
	sb.WriteString("\n\n")
	sb.WriteString(m.textInput.View())
	sb.WriteString("\n")

	if m.err != nil {
		sb.WriteString("\n")
		sb.WriteString(s.Error.Render(fmt.Sprintf("error: %s", m.err)))
		return sb.String()
	}

	if len(m.results) > 0 {
		if m.loading {
			sp := spinnerFrame(m.frame)
			fmt.Fprintf(&sb, "\n  %s\n\n", s.Info.Render(fmt.Sprintf("%s searching...", sp)))
		}
		sb.WriteString("\n")
		sb.WriteString(s.Info.Render(fmt.Sprintf("* %d results", m.totalHits)))
		sb.WriteString("\n\n")

		maxResults := len(m.results)
		h := m.height
		if h < 10 {
			h = 10
		}
		itemsPerPage := (h - 9) / 3
		if itemsPerPage < 3 {
			itemsPerPage = 3
		}

		start := m.cursor - itemsPerPage/2
		if start < 0 {
			start = 0
		}
		end := start + itemsPerPage
		if end > maxResults {
			end = maxResults
			start = end - itemsPerPage
			if start < 0 {
				start = 0
			}
		}

		for i := start; i < end; i++ {
			hit := m.results[i]

			if i == m.cursor {
				sb.WriteString(s.Selected.Render(fmt.Sprintf("> %s", hit.Title)))
				sb.WriteString("\n")
				fmt.Fprintf(&sb, "  %s", s.Dimmed.Render(truncateText(hit.Description, 90)))
				sb.WriteString("\n")
				fmt.Fprintf(&sb, "  %s\n",
					s.Info.Render(fmt.Sprintf("%s  %s  ↓ %s",
						strings.Join(hit.Loaders, ", "),
						strings.Join(hit.GameVersions, ", "),
						formatNumber(hit.Downloads))))
			} else {
				fmt.Fprintf(&sb, "  %s\n", hit.Title)
				fmt.Fprintf(&sb, "    %s\n", s.Dimmed.Render(truncateText(hit.Description, 80)))
			}

			if i < end-1 {
				sb.WriteString("\n")
			}
		}

		if m.totalHits > itemsPerPage {
			sb.WriteString("\n")
			sb.WriteString(s.Dimmed.Render(scrollLine(m.cursor, m.totalHits, itemsPerPage)))
		}
	} else if m.loading {
		sp := spinnerFrame(m.frame)
		fmt.Fprintf(&sb, "\n  %s\n", s.Info.Render(fmt.Sprintf("%s searching...", sp)))
	} else if len(m.textInput.Value()) > 0 {
		sb.WriteString("\n")
		sb.WriteString(s.Info.Render("no results found"))
	} else {
		sb.WriteString("\n")
		sb.WriteString(s.Info.Render("/ to search - type 3+ chars"))
	}

	sb.WriteString("\n\n")
	sb.WriteString(s.Dimmed.Render("j/k navigate - enter select - d: download - D: quick dl - esc back"))

	return sb.String()
}
